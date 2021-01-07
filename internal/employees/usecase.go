package employees

import (
	"context"

	"github.com/moguchev/service/internal/models"
)

// Usecase - business logic
type Usecase interface {
	CountEmployees(ctx context.Context, f models.EmployeeFilter) (int, error)
	GetEmployees(ctx context.Context, f models.EmployeeFilter) (models.Employees, error)
}
