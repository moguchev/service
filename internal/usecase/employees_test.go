package usecase

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/moguchev/service/internal/employees"
	"github.com/moguchev/service/internal/models"
)

type repoSuccess struct {
}

func (r *repoSuccess) CountEmployees(ctx context.Context, f models.EmployeeFilter) (uint, error) {
	return 1, nil
}

func (r *repoSuccess) GetEmployees(ctx context.Context, f models.EmployeeFilter) (models.Employees, error) {
	return models.Employees{models.Employee{}}, nil
}

func newRepoSuccess() employees.Repository {
	return &repoSuccess{}
}

func TestGetEmployees_Success(t *testing.T) {
	uc := NewEmployeesUsecase(newRepoSuccess())
	total, employees, err := uc.GetEmployees(context.Background(), models.EmployeeFilter{})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if total != 1 {
		t.Errorf("wrong total expected: %v, got: %v", 1, total)
	}

	expected := models.Employees{models.Employee{}}
	if !reflect.DeepEqual(employees, expected) {
		t.Errorf("wrong eployees expected: %v, got: %v", expected, employees)
	}
}

type repoCountFail struct {
}

func (r *repoCountFail) CountEmployees(ctx context.Context, f models.EmployeeFilter) (uint, error) {
	return 1, fmt.Errorf("error")
}

func (r *repoCountFail) GetEmployees(ctx context.Context, f models.EmployeeFilter) (models.Employees, error) {
	return models.Employees{models.Employee{}}, nil
}

func newRepoCountFail() employees.Repository {
	return &repoCountFail{}
}

func TestGetEmployees_CountFail(t *testing.T) {
	uc := NewEmployeesUsecase(newRepoCountFail())
	_, _, err := uc.GetEmployees(context.Background(), models.EmployeeFilter{})
	if err == nil {
		t.Errorf("expected error: %v", fmt.Errorf("error"))
	}
}

type repoGetEmployeesFail struct{}

func (r *repoGetEmployeesFail) CountEmployees(ctx context.Context, f models.EmployeeFilter) (uint, error) {
	return 1, nil
}

func (r *repoGetEmployeesFail) GetEmployees(ctx context.Context, f models.EmployeeFilter) (models.Employees, error) {
	return models.Employees{models.Employee{}}, fmt.Errorf("error")
}

func newRepoGetEmployeesFail() employees.Repository {
	return &repoGetEmployeesFail{}
}

func TestGetEmployees_GetEmployeesFail(t *testing.T) {
	uc := NewEmployeesUsecase(newRepoGetEmployeesFail())
	_, _, err := uc.GetEmployees(context.Background(), models.EmployeeFilter{})
	if err == nil {
		t.Errorf("expected error: %v", fmt.Errorf("error"))
	}
}

type repoNoEmployee struct{}

func (r *repoNoEmployee) CountEmployees(ctx context.Context, f models.EmployeeFilter) (uint, error) {
	return 0, nil
}

func (r *repoNoEmployee) GetEmployees(ctx context.Context, f models.EmployeeFilter) (models.Employees, error) {
	return models.Employees{models.Employee{}}, nil
}

func newRepoNoEmployee() employees.Repository {
	return &repoNoEmployee{}
}

func TestGetEmployees_NoEmployees(t *testing.T) {
	uc := NewEmployeesUsecase(newRepoNoEmployee())
	_, _, err := uc.GetEmployees(context.Background(), models.EmployeeFilter{})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
