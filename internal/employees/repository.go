package employees

import (
	"context"

	"github.com/moguchev/service/internal/models"
)

// Repository - database level
type Repository interface {
	CountEmployees(ctx context.Context, f models.EmployeeFilter) (uint, error)
	GetEmployees(ctx context.Context, f models.EmployeeFilter) (models.Employees, error)
}
