package repository

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/moguchev/service/internal/employees"
	"github.com/moguchev/service/internal/models"
	"github.com/moguchev/service/pkg/logger"
	"github.com/sirupsen/logrus"
)

type employeesRepository struct {
	db *sqlx.DB
}

// NewEmployeesRepository will create an object that represent the employees.Repository interface
func NewEmployeesRepository(db *sqlx.DB) employees.Repository {
	return &employeesRepository{db: db}
}

func applyEmployeeFilter(sb sq.SelectBuilder, f models.EmployeeFilter) sq.SelectBuilder {
	expr := sq.And{}

	if f.AssignmentID != nil {
		expr = append(expr, sq.Eq{"employees.assignment_id": *f.AssignmentID})
	}

	if f.EmployeeID != nil {
		expr = append(expr, sq.Eq{"employees.employee_id": *f.EmployeeID})
	}

	if f.FIO != nil {
		expr = append(expr, sq.Expr("employees.fio ILIKE ?", fmt.Sprint("%", *f.FIO, "%")))
	}

	if f.JobName != nil {
		expr = append(expr, sq.Expr("employees.job_name ILIKE ?", fmt.Sprint("%", *f.JobName, "%")))
	}

	if len(expr) > 0 {
		sb = sb.Where(expr)
	}

	orders := []string{}

	if f.DateFromSort != nil {
		orders = append(orders, fmt.Sprintf("salaries.date_from %s", *f.DateFromSort))
	}

	if f.SalarySort != nil {
		orders = append(orders, fmt.Sprintf("salaries.salary %s", *f.DateFromSort))
	}

	if len(orders) > 0 {
		sb = sb.OrderBy(orders...)
	}

	if f.Limit != nil {
		sb = sb.Limit(*f.Limit)
	}

	if f.Offset != nil {
		sb = sb.Offset(*f.Offset)
	}

	return sb
}

func (r *employeesRepository) CountEmployees(ctx context.Context, f models.EmployeeFilter) (uint, error) {
	log := logger.GetLogger(ctx).WithFields(logrus.Fields{
		"actor":  "repository",
		"func":   "CountEmployees",
		"filter": f,
	})

	query := sq.Select("COUNT(employee_id)").From("employees").
		Join("salaries ON employees.assignment_id = salaries.assignment_id").
		PlaceholderFormat(sq.Dollar)

	f.Limit = nil
	f.Offset = nil
	f.SalarySort = nil
	f.DateFromSort = nil

	query = applyEmployeeFilter(query, f)

	sql, args, err := query.ToSql()
	if err != nil {
		return 0, fmt.Errorf("to sql: %w", err)
	}

	log = log.WithFields(logrus.Fields{"query": sql, "args": args})

	log.Debug("count employees")

	var count uint
	if err = r.db.QueryRowxContext(ctx, sql, args...).Scan(&count); err != nil {
		log.WithError(err).Error("count employees")
		return 0, fmt.Errorf("count employees: %w", err)
	}

	return count, nil
}

func (r *employeesRepository) GetEmployees(ctx context.Context, f models.EmployeeFilter) (models.Employees, error) {
	log := logger.GetLogger(ctx).WithFields(logrus.Fields{
		"actor": "repository",
		"func":  "GetEmployees",
	})

	cols := []string{"employees.employee_id", "employees.assignment_id", "employees.fio",
		"employees.job_name", "salaries.salary", "salaries.date_from"}

	query := sq.Select(cols...).From("employees").
		Join("salaries ON employees.assignment_id = salaries.assignment_id").
		PlaceholderFormat(sq.Dollar)

	query = applyEmployeeFilter(query, f)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("to sql: %w", err)
	}

	log = log.WithFields(logrus.Fields{"query": sql, "args": args})

	log.Debug("get employees")

	rows, err := r.db.QueryxContext(ctx, sql, args...)
	if err != nil {
		log.WithError(err).Error("get employees")
		return nil, fmt.Errorf("get employees: %w", err)
	}
	defer rows.Close()

	emps := models.Employees{}

	for rows.Next() {
		employee := models.Employee{}

		if err = rows.StructScan(&employee); err != nil {
			log.WithError(err).Error("scan employee")
			return nil, fmt.Errorf("scan employee: %w", err)
		}

		emps = append(emps, employee)
	}

	return emps, nil
}
