package delivery

import (
	"context"
	"net/http"

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

func (h *EmployeesHandler) GetEmployees(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	log := logger.GetLogger(ctx).WithField("handler", "GetEmployees")

	log.Info("OK")

	utils.RespondWithJSON(w, r, http.StatusOK, struct{}{})
}

func (h *EmployeesHandler) GetEmployeeByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	log := logger.GetLogger(ctx).WithField("handler", "GetEmployees")

	log.Info("OK")

	utils.RespondWithJSON(w, r, http.StatusOK, struct{}{})
}
