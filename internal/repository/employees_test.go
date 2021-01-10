package repository

import (
	"context"
	"database/sql/driver"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/moguchev/service/internal/models"

	sq "github.com/Masterminds/squirrel"
)

func TestApplyEmployeeFilter(t *testing.T) {
	var (
		id   int64  = 1
		u    uint64 = 1
		str         = "string"
		sort        = models.ASC
	)

	f := models.EmployeeFilter{
		AssignmentID: &id,
		EmployeeID:   &id,
		FIO:          &str,
		JobName:      &str,
		Limit:        &u,
		Offset:       &u,
		DateFromSort: &sort,
		SalarySort:   &sort,
	}

	query := sq.Select("*").PlaceholderFormat(sq.Dollar)
	query = applyEmployeeFilter(query, f)

	sql, _, err := query.ToSql()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	expected := `SELECT * WHERE (employees.assignment_id = $1 AND employees.employee_id = $2 AND employees.fio ILIKE $3 AND employees.job_name ILIKE $4) ORDER BY salaries.date_from ASC, salaries.salary ASC LIMIT 1 OFFSET 1`
	if sql != expected {
		t.Errorf("func returned unexpected query: got %v want %v", sql, expected)
	}
}

func TestCountEmployees_Success(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer mockDB.Close()
	db := sqlx.NewDb(mockDB, "sqlmock")

	var expected uint = 1
	rows := sqlmock.NewRows([]string{"employee_id"}).AddRow(expected)
	mock.ExpectQuery("SELECT COUNT(.+) FROM (.+)").WillReturnRows(rows)

	repo := NewEmployeesRepository(db)
	count, err := repo.CountEmployees(context.Background(), models.EmployeeFilter{})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if count != expected {
		t.Errorf("func returned unexpected value: got %v want %v", count, expected)
	}
}

func TestCountEmployees_Fail(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer mockDB.Close()
	db := sqlx.NewDb(mockDB, "sqlmock")

	dberr := fmt.Errorf("error")
	mock.ExpectQuery("SELECT COUNT(.+) FROM (.+)").WillReturnError(dberr)

	repo := NewEmployeesRepository(db)
	_, err = repo.CountEmployees(context.Background(), models.EmployeeFilter{})
	if err == nil {
		t.Error("expected error")
	}
}

func TestGetEmployees_Success(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer mockDB.Close()
	db := sqlx.NewDb(mockDB, "sqlmock")

	cols := []string{"employee_id", "assignment_id", "fio",
		"job_name", "salary", "date_from"}
	sal := 400000.99
	date := time.Now()
	employee := models.Employee{
		EmployeeID:   1,
		AssignmentID: 1,
		FIO:          "string",
		JobName:      "string",
		Salary:       &sal,
		DateFrom:     &date,
	}

	row := []driver.Value{employee.EmployeeID, employee.AssignmentID, employee.FIO, employee.JobName,
		employee.Salary, employee.DateFrom}
	rows := sqlmock.NewRows(cols).AddRow(row...)
	mock.ExpectQuery("SELECT (.+) FROM (.+)").WillReturnRows(rows)

	repo := NewEmployeesRepository(db)
	employees, err := repo.GetEmployees(context.Background(), models.EmployeeFilter{})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	expexted := models.Employees{employee}
	if !reflect.DeepEqual(employees, expexted) {
		t.Errorf("expected: %v, got: %v", expexted, employees)
	}
}

func TestGetEmployees_Fail(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer mockDB.Close()
	db := sqlx.NewDb(mockDB, "sqlmock")

	mock.ExpectQuery("SELECT (.+) FROM (.+)").WillReturnError(fmt.Errorf("error"))

	repo := NewEmployeesRepository(db)
	_, err = repo.GetEmployees(context.Background(), models.EmployeeFilter{})
	if err == nil {
		t.Error("expected error")
	}
}

func TestGetEmployees_FailScan(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer mockDB.Close()
	db := sqlx.NewDb(mockDB, "sqlmock")

	cols := []string{"employee_id", "assignment_id", "fio",
		"job_name", "salary", "date_from"}
	sal := 400000.99
	date := time.Now()
	employee := models.Employee{
		EmployeeID:   1,
		AssignmentID: 1,
		FIO:          "string",
		JobName:      "string",
		Salary:       &sal,
		DateFrom:     &date,
	}

	row := []driver.Value{nil, employee.AssignmentID, employee.FIO, employee.JobName,
		employee.Salary, employee.DateFrom}
	rows := sqlmock.NewRows(cols).AddRow(row...)
	mock.ExpectQuery("SELECT (.+) FROM (.+)").WillReturnRows(rows)

	repo := NewEmployeesRepository(db)
	_, err = repo.GetEmployees(context.Background(), models.EmployeeFilter{})
	if err == nil {
		t.Error("expected error")
	}
}
