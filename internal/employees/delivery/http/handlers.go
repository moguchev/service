package delivery

import (
	"context"
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

// EmployeesHandler represent the http handler for employees
type EmployeesHandler struct {
	Usecase employees.Usecase
}

// SetEmployeesHandler will initialize the employee(s)/ resources endpoint
func SetEmployeesHandler(router *mux.Router, us employees.Usecase) {
	handler := &EmployeesHandler{
		Usecase: us,
	}

	router.HandleFunc("/employees", handler.GetEmployees).Methods(http.MethodGet)
	router.HandleFunc("/employees/{employee_id}", handler.GetEmployeeByID).Methods(http.MethodGet)
}

// GetEmployees -
func (h *EmployeesHandler) GetEmployees(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	log := logger.GetLogger(ctx).WithField("handler", "GetEmployees")

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
		case "employee_id":
			employeeID, e := strconv.ParseInt(v, 10, 64)
			if e != nil {
				err = fmt.Errorf("employee_id: %w", e)
				break
			}
			f.EmployeeID = &employeeID
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

// GetEmployeeByID -
func (h *EmployeesHandler) GetEmployeeByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	log := logger.GetLogger(ctx).WithField("handler", "GetEmployees")

	log.Info("OK")

	utils.RespondWithJSON(w, r, http.StatusOK, struct{}{})
}
