package models

import "time"

// SortOrder - type of order
type SortOrder int

const (
	// ASC - ascending order
	ASC SortOrder = 0
	// DESC - descending order
	DESC SortOrder = 1
)

type (
	// EmployeeFilter - struct with filter
	EmployeeFilter struct {
		Limit        *int
		Offset       *int
		FIO          *string
		EmployeeID   *int64
		AssignmentID *int64
		JobName      *string
		DateFromSort *SortOrder
		SalarySort   *SortOrder
	}

	// Employee - employee info
	Employee struct {
		EmployeeID   int64     `json:"employee_id" db:"employee_id"`
		AssignmentID int64     `json:"assignment_id" db:"assignment_id"`
		FIO          string    `json:"fio" db:"fio"`
		JobName      string    `json:"job_name" db:"job_name"`
		Salary       float64   `json:"salary" db:"salary"`
		DateFrom     time.Time `json:"date_from" db:"date_from"`
	}

	// Employees - array of employees info
	Employees []Employee
)
