package delivery

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gorilla/mux"

	"bytes"
	"time"

	"github.com/moguchev/service/internal/employees"
	"github.com/moguchev/service/internal/models"
)

func TestGetEmployeeFilter(t *testing.T) {
	str := "string"
	values := url.Values{
		"limit":          {"1"},
		"offset":         {"1"},
		"fio":            {str},
		"assignment_id":  {"1"},
		"job_name":       {str},
		"date_from_sort": {string(models.ASC)},
		"salary_sort":    {string(models.DESC)},
	}

	filter, err := getEmployeeFilter(values)
	if err != nil {
		t.Errorf("unexpected error")
	}

	if filter.Limit == nil || *filter.Limit != 1 {
		t.Errorf("Limit")
	}

	if filter.Offset == nil || *filter.Offset != 1 {
		t.Errorf("Offset")
	}

	if filter.FIO == nil || *filter.FIO != str {
		t.Errorf("FIO")
	}

	if filter.AssignmentID == nil || *filter.AssignmentID != 1 {
		t.Errorf("AssignmentID")
	}

	if filter.JobName == nil || *filter.JobName != str {
		t.Errorf("JobName")
	}

	if filter.DateFromSort == nil || *filter.DateFromSort != models.ASC {
		t.Errorf("DateFromSort")
	}

	if filter.SalarySort == nil || *filter.SalarySort != models.DESC {
		t.Errorf("SalarySort")
	}
}

func TestGetEmployeeFilter_Error(t *testing.T) {
	type testCase struct {
		values url.Values
		err    string
	}

	testCases := []testCase{
		{url.Values{"limit": {"not number"}}, "limit"},
		{url.Values{"offset": {"not number"}}, "offset"},
		{url.Values{"assignment_id": {"not number"}}, "assignment_id"},
		{url.Values{"date_from_sort": {"not order"}}, "date_from_sort"},
		{url.Values{"salary_sort": {"not order"}}, "salary_sort"},
	}

	for i, test := range testCases {
		_, err := getEmployeeFilter(test.values)
		if err == nil {
			t.Errorf("test = %v, expected error", i)
		}
		if !strings.Contains(err.Error(), test.err) {
			t.Errorf("test = %v, expected error with: %v, got: %v", i, test.err, err)
		}
	}
}

type employeesUsecaseSuccessMock struct{}

func (mock *employeesUsecaseSuccessMock) GetEmployees(ctx context.Context, f models.EmployeeFilter) (uint, models.Employees, error) {
	date, _ := time.Parse(time.RFC3339, "2020-07-23T00:00:00Z")
	var salary float64 = 400000
	return 1, models.Employees{models.Employee{
		AssignmentID: 648078,
		EmployeeID:   775900,
		FIO:          "Могучев Леонид Алексеевич",
		JobName:      "старший разработчик",
		Salary:       &salary,
		DateFrom:     &date,
	}}, nil
}

func NewEmployeesUsecaseSuccessMock() employees.Usecase {
	return &employeesUsecaseSuccessMock{}
}

func TestSetEmployeesHandler(t *testing.T) {
	router := mux.NewRouter()
	uc := NewEmployeesUsecaseSuccessMock()

	SetEmployeesHandler(router, uc)

	req, err := http.NewRequest(http.MethodGet, "/employees/775900", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status == http.StatusNotFound {
		t.Errorf("handler not found")
	}

	req, err = http.NewRequest(http.MethodGet, "/employees", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status == http.StatusNotFound {
		t.Errorf("handler not found")
	}
}

func TestGetEmployeesHandler_Success(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/employees", nil)
	if err != nil {
		t.Fatal(err)
	}

	uc := NewEmployeesUsecaseSuccessMock()
	h := EmployeesHandler{uc}
	handler := http.HandlerFunc(h.GetEmployeesHandler)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := []byte(`{"total":1,"employees":[{"employee_id":775900,"assignment_id":648078,"fio":"Могучев Леонид Алексеевич","job_name":"старший разработчик","salary":400000,"date_from":"2020-07-23T00:00:00Z"}]}` + "\n")
	if !bytes.Equal(rr.Body.Bytes(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			string(rr.Body.Bytes()), string(expected))
	}
}

func TestGetEmployeesHandler_BadRequest(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/employees?limit=str", nil)
	if err != nil {
		t.Fatal(err)
	}

	uc := NewEmployeesUsecaseSuccessMock()
	h := EmployeesHandler{uc}
	handler := http.HandlerFunc(h.GetEmployeesHandler)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	msg := models.ErrorMessage{}

	if err := json.Unmarshal(rr.Body.Bytes(), &msg); err != nil {
		t.Errorf("handler returned unexpected body: got %v want %v",
			string(rr.Body.Bytes()), msg)
	}

	if len(msg.Message) == 0 {
		t.Errorf("expected error message")
	}
}

type employeesUsecaseBadMock struct{}

func (mock *employeesUsecaseBadMock) GetEmployees(ctx context.Context, f models.EmployeeFilter) (uint, models.Employees, error) {
	return 0, models.Employees{}, models.ErrInternal
}

func NewEmployeesUsecaseBadMock() employees.Usecase {
	return &employeesUsecaseBadMock{}
}

func TestGetEmployeesHandler_InternalError(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/employees/", nil)
	if err != nil {
		t.Fatal(err)
	}

	uc := NewEmployeesUsecaseBadMock()
	h := EmployeesHandler{uc}
	handler := http.HandlerFunc(h.GetEmployeesHandler)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}

	msg := models.ErrorMessage{}

	if err := json.Unmarshal(rr.Body.Bytes(), &msg); err != nil {
		t.Errorf("handler returned unexpected body: got %v want %v",
			string(rr.Body.Bytes()), msg)
	}

	if len(msg.Message) == 0 {
		t.Errorf("expected error message")
	}
}

func TestGetEmployeeByIDHandler_BadRequest(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/employees/abcd", nil)
	if err != nil {
		t.Fatal(err)
	}

	uc := NewEmployeesUsecaseSuccessMock()
	h := EmployeesHandler{uc}
	handler := http.HandlerFunc(h.GetEmployeeByIDHandler)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	msg := models.ErrorMessage{}

	if err := json.Unmarshal(rr.Body.Bytes(), &msg); err != nil {
		t.Errorf("handler returned unexpected body: got %v want %v",
			string(rr.Body.Bytes()), msg)
	}

	if len(msg.Message) == 0 {
		t.Errorf("expected error message")
	}
}

func TestGetEmployeeByIDHandler_InternalError(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/employees/775900", nil)
	if err != nil {
		t.Fatal(err)
	}

	uc := NewEmployeesUsecaseBadMock()
	h := EmployeesHandler{uc}
	router := mux.NewRouter()
	router.HandleFunc(fmt.Sprintf("/employees/{%s}", employeeIDParam), h.GetEmployeeByIDHandler)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}

	msg := models.ErrorMessage{}

	if err := json.Unmarshal(rr.Body.Bytes(), &msg); err != nil {
		t.Errorf("handler returned unexpected body: got %v want %v",
			string(rr.Body.Bytes()), msg)
	}

	if len(msg.Message) == 0 {
		t.Errorf("expected error message")
	}
}

type employeesUsecaseEmptyMock struct{}

func (mock *employeesUsecaseEmptyMock) GetEmployees(ctx context.Context, f models.EmployeeFilter) (uint, models.Employees, error) {
	return 0, models.Employees{}, nil
}

func NewEmployeesUsecaseEmptyMock() employees.Usecase {
	return &employeesUsecaseEmptyMock{}
}

func TestGetEmployeeByIDHandler_NotFound(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/employees/775900", nil)
	if err != nil {
		t.Fatal(err)
	}

	uc := NewEmployeesUsecaseEmptyMock()
	h := EmployeesHandler{uc}
	router := mux.NewRouter()
	router.HandleFunc(fmt.Sprintf("/employees/{%s}", employeeIDParam), h.GetEmployeeByIDHandler)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}

	expected := []byte("{}\n")
	if !bytes.Equal(rr.Body.Bytes(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			string(rr.Body.Bytes()), string(expected))
	}

}

func TestGetEmployeeByIDHandler(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/employees/775900", nil)
	if err != nil {
		t.Fatal(err)
	}

	uc := NewEmployeesUsecaseSuccessMock()
	h := EmployeesHandler{uc}
	router := mux.NewRouter()
	router.HandleFunc(fmt.Sprintf("/employees/{%s}", employeeIDParam), h.GetEmployeeByIDHandler)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}

	expected := []byte(`{"employee_id":775900,"assignment_id":648078,"fio":"Могучев Леонид Алексеевич","job_name":"старший разработчик","salary":400000,"date_from":"2020-07-23T00:00:00Z"}` + "\n")
	if !bytes.Equal(rr.Body.Bytes(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			string(rr.Body.Bytes()), string(expected))
	}
}
