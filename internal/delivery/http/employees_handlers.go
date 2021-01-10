package delivery

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/moguchev/service/internal/models"

	"github.com/gorilla/mux"
	"github.com/moguchev/service/internal/employees"
	"github.com/moguchev/service/pkg/logger"
	"github.com/moguchev/service/pkg/utils"
)

const (
	employeeIDParam = "employee_id"
)

// EmployeesHandler represent the http handler for employees
type EmployeesHandler struct {
	Usecase employees.Usecase
}

// SetEmployeesHandler will initialize the employee(s)/ resources endpoint
func SetEmployeesHandler(router *mux.Router, uc employees.Usecase) {
	handler := &EmployeesHandler{
		Usecase: uc,
	}

	router.HandleFunc("/employees", handler.GetEmployeesHandler).Methods(http.MethodGet)
	router.HandleFunc(fmt.Sprintf("/employees/{%s}", employeeIDParam),
		handler.GetEmployeeByIDHandler).Methods(http.MethodGet)
}

// GetEmployeesHandler -
func (h *EmployeesHandler) GetEmployeesHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.GetLogger(ctx).WithField("handler", "GetEmployeesHandler")

	filter, err := getEmployeeFilter(r.URL.Query())
	if err != nil {
		log.WithError(err).Error("parse query parameters")
		utils.RespondWithError(w, r, http.StatusBadRequest, err)
		return
	}

	log.WithField("filter", filter).Debug("get employees")

	total, emps, err := h.Usecase.GetEmployees(ctx, filter)
	if err != nil {
		log.WithError(err).Error("get employees")
		utils.RespondWithError(w, r, http.StatusInternalServerError, models.ErrInternal)
		return
	}

	type Response struct {
		Total     uint             `json:"total"`
		Employees models.Employees `json:"employees"`
	}

	utils.RespondWithJSON(w, r, http.StatusOK, Response{Total: total, Employees: emps})
}

// nolint:gocyclo // mapping
func getEmployeeFilter(values url.Values) (models.EmployeeFilter, error) {
	f := models.EmployeeFilter{}
	for k, vs := range values {
		v := vs[0]
		var err error
		switch k {
		case "limit":
			limit, e := strconv.ParseUint(v, 10, 64)
			if e != nil {
				err = fmt.Errorf("limit: %w", e)
				break
			}
			f.Limit = &limit
		case "offset":
			offset, e := strconv.ParseUint(v, 10, 64)
			if e != nil {
				err = fmt.Errorf("offset: %w", e)
				break
			}
			f.Offset = &offset
		case "fio":
			fio := v
			f.FIO = &fio
		case "assignment_id":
			assignment, e := strconv.ParseInt(v, 10, 64)
			if e != nil {
				err = fmt.Errorf("assignment_id: %w", e)
				break
			}
			f.AssignmentID = &assignment
		case "job_name":
			job := v
			f.JobName = &job
		case "date_from_sort":
			order := models.SortOrder(v)
			if order != models.ASC && order != models.DESC {
				err = fmt.Errorf("date_from_sort: wrong order %s", order)
				break
			}
			f.DateFromSort = &order
		case "salary_sort":
			order := models.SortOrder(v)
			if order != models.ASC && order != models.DESC {
				err = fmt.Errorf("salary_sort: wrong order %s", order)
				break
			}
			f.SalarySort = &order
		}

		if err != nil {
			return models.EmployeeFilter{}, err
		}
	}

	return f, nil
}

// GetEmployeeByIDHandler -
func (h *EmployeesHandler) GetEmployeeByIDHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.GetLogger(ctx).WithField("handler", "GetEmployeeByIDHandler")

	id := mux.Vars(r)[employeeIDParam]

	log.WithField("url", r.URL).Info("url")
	log.WithField(employeeIDParam, id).Info("id")

	empID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.WithError(err).WithField(employeeIDParam, id).Error("parse")
		utils.RespondWithError(w, r, http.StatusBadRequest, err)
		return
	}

	_, emps, err := h.Usecase.GetEmployees(ctx, models.EmployeeFilter{EmployeeID: &empID})
	if err != nil {
		log.WithError(err).Error("get employee by id")
		utils.RespondWithError(w, r, http.StatusInternalServerError, models.ErrInternal)
		return
	}

	if len(emps) == 0 {
		utils.RespondWithJSON(w, r, http.StatusNotFound, struct{}{})
		return
	}

	utils.RespondWithJSON(w, r, http.StatusOK, emps[0])
}
