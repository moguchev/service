package models

import "time"

// SortOrder - type of order
type SortOrder string

const (
	// ASC - ascending order
	ASC SortOrder = "ASC"
	// DESC - descending order
	DESC SortOrder = "DESC"
)

type (
	// EmployeeFilter - struct with filter
	EmployeeFilter struct {
		Limit        *uint64
		Offset       *uint64
		FIO          *string
		EmployeeID   *int64
		AssignmentID *int64
		JobName      *string
		DateFromSort *SortOrder
		SalarySort   *SortOrder
	}

	// Employee - employee info
	Employee struct {
		EmployeeID   int64      `json:"employee_id" db:"employee_id"`
		AssignmentID int64      `json:"assignment_id" db:"assignment_id"`
		FIO          string     `json:"fio" db:"fio"`
		JobName      string     `json:"job_name" db:"job_name"`
		Salary       float64    `json:"salary" db:"salary"`
		DateFrom     *time.Time `json:"date_from,omitempty" db:"date_from"`
	}

	// Employees - array of employees info
	Employees []Employee
)
