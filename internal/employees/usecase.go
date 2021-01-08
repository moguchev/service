package employees

import (
	"context"

	"github.com/moguchev/service/internal/models"
)

// Usecase - business logic
type Usecase interface {
	GetEmployees(ctx context.Context, f models.EmployeeFilter) (uint, models.Employees, error)
}
