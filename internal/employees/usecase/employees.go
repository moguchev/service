package usecase

import (
	"context"
	"fmt"

	"github.com/moguchev/service/internal/employees"
	"github.com/moguchev/service/internal/models"
	"github.com/moguchev/service/pkg/logger"
	"github.com/sirupsen/logrus"
)

type employeesUsecase struct {
	empRepo employees.Repository
}

// NewEmployeesUsecase will create new an employeesUsecase object representation of employees.Usecase interface
func NewEmployeesUsecase(eRepo employees.Repository) employees.Usecase {
	return &employeesUsecase{empRepo: eRepo}
}

func (e *employeesUsecase) GetEmployees(ctx context.Context, f models.EmployeeFilter) (uint, models.Employees, error) {
	log := logger.GetLogger(ctx).WithFields(logrus.Fields{
		"actor":  "usecase",
		"func":   "GetEmployees",
		"filter": f,
	})

	emps := models.Employees{}

	total, err := e.empRepo.CountEmployees(ctx, f)
	if err != nil {
		log.WithError(err).Error("count employees")
		return 0, nil, fmt.Errorf("count employees: %w", err)
	}

	if total == 0 {
		return total, emps, nil
	}

	emps, err = e.empRepo.GetEmployees(ctx, f)
	if err != nil {
		log.WithError(err).Error("get employees")
		return 0, nil, fmt.Errorf("get employees: %w", err)
	}

	return total, emps, nil
}
