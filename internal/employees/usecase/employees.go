package usecase

import (
	"context"

	"github.com/moguchev/service/internal/employees"
	"github.com/moguchev/service/internal/models"
)

type employeesUsecase struct{}

// NewEmployeesUsecase will create new an employeesUsecase object representation of employees.Usecase interface
func NewEmployeesUsecase() employees.Usecase {
	return &employeesUsecase{}
}

func (e *employeesUsecase) CountEmployees(ctx context.Context, f models.EmployeeFilter) (int, error) {
	return 0, nil
}

func (e *employeesUsecase) GetEmployees(ctx context.Context, f models.EmployeeFilter) (models.Employees, error) {
	return nil, nil
}
