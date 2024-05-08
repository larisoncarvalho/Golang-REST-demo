package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

// CREATE EMPLOYEE
func TestCreateEmployeeHandler_PASS(t *testing.T) {
	db := setupDatabase()
	handler := Handler{db: db}
	defer handler.db.Close()

	// Create a new request with a JSON body
	employee := Employee{ID: 1, Name: "John Doe", Position: "Engineer", Salary: 50000}
	reqBody, _ := json.Marshal(employee)
	req := httptest.NewRequest("POST", "/createEmployee", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler function
	handler.createEmployeeHandler(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusCreated, rr.Code)
}

func TestCreateEmployeeHandler_FAIL_Missing_ID(t *testing.T) {
	db := setupDatabase()
	handler := Handler{db: db}
	defer handler.db.Close()

	// Create a new request with a JSON body without ID
	employee := Employee{Name: "John Doe", Position: "Engineer", Salary: 50000}
	reqBody, _ := json.Marshal(employee)
	req := httptest.NewRequest("POST", "/createEmployee", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler function
	handler.createEmployeeHandler(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "Employee ID cannot be 0")
}

func TestCreateEmployeeHandler_FAIL_Duplicate(t *testing.T) {
	db := setupDatabase()
	handler := Handler{db: db}
	defer handler.db.Close()

	// Create a new request with a JSON body without ID
	employee := Employee{ID: 44, Name: "John Doe", Position: "Engineer", Salary: 50000}
	reqBody, _ := json.Marshal(employee)
	req := httptest.NewRequest("POST", "/createEmployee", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler function
	handler.createEmployeeHandler(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusConflict, rr.Code)
	assert.Contains(t, rr.Body.String(), "UNIQUE constraint failed: employees.ID")
}

// UPDATE EMPLOYEE
func TestUpdateEmployeeHandler_PASS(t *testing.T) {
	db := setupDatabase()
	handler := Handler{db: db}
	defer handler.db.Close()

	// Create a new employee object for updating
	newEmployee := Employee{ID: 2, Name: "Alice Smith", Position: "Senior Manager", Salary: 70000}

	// Encode the new employee object to JSON
	reqBody, _ := json.Marshal(newEmployee)

	// Create a request to update the employee
	req := httptest.NewRequest("PUT", "/updateEmployee", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler function
	handler.updateEmployeeHandler(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestUpdateEmployeeHandler_FAIL_Employee_Doesnt_Exist(t *testing.T) {
	db := setupDatabase()
	handler := Handler{db: db}
	defer handler.db.Close()

	// Create a new employee object for updating
	newEmployee := Employee{ID: 99, Name: "Alice Smith", Position: "Senior Manager", Salary: 70000}

	// Encode the new employee object to JSON
	reqBody, _ := json.Marshal(newEmployee)

	// Create a request to update the employee
	req := httptest.NewRequest("PUT", "/updateEmployee", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler function
	handler.updateEmployeeHandler(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestDeleteEmployeeHandler_PASS(t *testing.T) {
	db := setupDatabase()
	handler := Handler{db: db}
	defer handler.db.Close()

	// Create a request to update the employee
	req := httptest.NewRequest("DELETE", "/deleteEmployee/{id}", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "2")

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// Create a response recorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler function
	handler.deleteEmployeeHandler(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)
}

// DELETE EMPLOYEE
func TestDeleteEmployeeHandler_FAIL_Does_Not_Exist(t *testing.T) {
	db := setupDatabase()
	handler := Handler{db: db}
	defer handler.db.Close()

	// Create a request to update the employee
	req := httptest.NewRequest("DELETE", "/deleteEmployee/{id}", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "22")

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// Create a response recorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler function
	handler.deleteEmployeeHandler(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusNotFound, rr.Code)
	assert.Contains(t, rr.Body.String(), "Employee does not exist.")
}

func TestDeleteEmployeeHandler_FAIL_Invalid_Id(t *testing.T) {
	db := setupDatabase()
	handler := Handler{db: db}
	defer handler.db.Close()

	// Create a request to update the employee
	req := httptest.NewRequest("DELETE", "/deleteEmployee/{id}", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "aa")

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// Create a response recorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler function
	handler.deleteEmployeeHandler(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "Error parsing the ID, make sure it is an integer")
}

// GET EMPLOYEE BY ID
func TestGetEmployeeHandler_PASS(t *testing.T) {
	db := setupDatabase()
	handler := Handler{db: db}
	defer handler.db.Close()

	// Create a request to update the employee
	req := httptest.NewRequest("GET", "/employees/{id}", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "2")

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// Create a response recorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler function
	handler.getEmployeeByIdHandler(rr, req)
	// Retrieve the response body
	responseBody := rr.Body.Bytes()

	var resultEmployee Employee
	if err := json.Unmarshal(responseBody, &resultEmployee); err != nil {
		t.Errorf("Error unmarshalling JSON: %v", err)
	}

	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, resultEmployee.ID, 2)
	assert.Equal(t, resultEmployee.Name, "Alice")
}

func TestGetEmployeeHandler_FAIL_Does_not_exist(t *testing.T) {
	db := setupDatabase()
	handler := Handler{db: db}
	defer handler.db.Close()

	// Create a request to update the employee
	req := httptest.NewRequest("GET", "/employees/{id}", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "22")

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// Create a response recorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler function
	handler.getEmployeeByIdHandler(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusNotFound, rr.Code)
}

// LIST EMPLOYEE
func TestListEmployeeHandler_PASS_page1_size2(t *testing.T) {
	db := setupDatabase()
	handler := Handler{db: db}
	defer handler.db.Close()

	// Create a request to update the employee
	req := httptest.NewRequest("GET", "/getEmployees?page=1&size=2", nil)

	// Create a response recorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler function
	handler.getEmployeesListHandler(rr, req)

	// Retrieve the response body
	responseBody := rr.Body.Bytes()

	var resultEmployees []Employee
	if err := json.Unmarshal(responseBody, &resultEmployees); err != nil {
		t.Errorf("Error unmarshalling JSON: %v", err)
	}

	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, len(resultEmployees), 2)
}

func TestListEmployeeHandler_PASS_page2_size3(t *testing.T) {
	db := setupDatabase()
	handler := Handler{db: db}
	defer handler.db.Close()

	// Create a request to update the employee
	req := httptest.NewRequest("GET", "/getEmployees?page=2&size=3", nil)

	// Create a response recorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler function
	handler.getEmployeesListHandler(rr, req)

	// Retrieve the response body
	responseBody := rr.Body.Bytes()

	var resultEmployees []Employee
	if err := json.Unmarshal(responseBody, &resultEmployees); err != nil {
		t.Errorf("Error unmarshalling JSON: %v", err)
	}

	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, len(resultEmployees), 1)
}

func TestListEmployeeHandler_PASS_default(t *testing.T) {
	db := setupDatabase()
	handler := Handler{db: db}
	defer handler.db.Close()

	// Create a request to update the employee
	req := httptest.NewRequest("GET", "/getEmployees", nil)

	// Create a response recorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler function
	handler.getEmployeesListHandler(rr, req)

	// Retrieve the response body
	responseBody := rr.Body.Bytes()

	var resultEmployees []Employee
	if err := json.Unmarshal(responseBody, &resultEmployees); err != nil {
		t.Errorf("Error unmarshalling JSON: %v", err)
	}

	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, len(resultEmployees), 4)
}

// SET UP
func setupDatabase() *sql.DB {
	db, _ := sql.Open("sqlite3", ":memory:")
	_, _ = db.Exec(`CREATE TABLE IF NOT EXISTS employees (
                        ID INTEGER PRIMARY KEY,
                        NAME TEXT,
                        POSITION TEXT,
                        SALARY REAL
                     )`)
	db.Exec("INSERT INTO employees (id, name, position, salary) VALUES (?, ?, ?, ?)",
		"44", "Duplicate", "Redundant", "99999")
	db.Exec("INSERT INTO employees (id, name, position, salary) VALUES (?, ?, ?, ?)",
		"2", "Alice", "Manager", "60000")
	db.Exec("INSERT INTO employees (id, name, position, salary) VALUES (?, ?, ?, ?)",
		"3", "Jack", "Writer", "2000")
	db.Exec("INSERT INTO employees (id, name, position, salary) VALUES (?, ?, ?, ?)",
		"4", "Mary", "Assistant", "1000")
	return db
}
