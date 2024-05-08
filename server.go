package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/mattn/go-sqlite3"
)

type Handler struct {
	db *sql.DB
}

func main() {
	port := "3000"

	// Open DB connection
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		log.Fatal(err)
	}

	// Store db in a handler struct so we can use it in our handler functions in a safe way
	handler := Handler{db: db}
	defer handler.db.Close()

	// Create a Chi Router, This handles concurrency of the mulitple requests
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Post("/createEmployee", handler.createEmployeeHandler)

	r.Get("/employees/{id}", handler.getEmployeeByIdHandler)

	r.Post("/updateEmployee", handler.updateEmployeeHandler)

	r.Delete("/deleteEmployee/{id}", handler.deleteEmployeeHandler)

	r.Get("/getEmployees", handler.getEmployeesListHandler)

	log.Println("Starting server on " + port)
	http.ListenAndServe(":"+port, r)
}

func (h *Handler) createEmployeeHandler(w http.ResponseWriter, r *http.Request) {
	// Parse request
	var employee Employee
	if err := json.NewDecoder(r.Body).Decode(&employee); err != nil {
		http.Error(w, "Request body is invalid", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	err := validateEmployee(employee)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// call DB layer
	err = createEmployee(h.db, employee)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			http.Error(w, "Employee with ID already exists. Error: "+
				err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, "Error while inserting employee. Error: "+
			err.Error(), http.StatusInternalServerError)
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) getEmployeeByIdHandler(w http.ResponseWriter, r *http.Request) {
	// Parse Request
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Error parsing the ID, make sure it is an integer. Error: "+err.Error(),
			http.StatusBadRequest)
		return
	}

	// Call DB layer
	employee, err := getEmployeeById(h.db, id)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			http.Error(w, "Employee does not exist.",
				http.StatusNotFound)
			return
		}
		http.Error(w, "Error while getting employee "+
			err.Error(), http.StatusInternalServerError)
		return
	}

	// Send Response
	response, err := json.Marshal(employee)
	if err != nil {
		http.Error(w, "Error while converting the db response to json. Error: "+err.Error(), http.StatusInternalServerError)
	}
	w.Write(response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) updateEmployeeHandler(w http.ResponseWriter, r *http.Request) {
	// Parse Request
	var employee Employee
	if err := json.NewDecoder(r.Body).Decode(&employee); err != nil {
		http.Error(w, "Request body is invalid", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	err := validateEmployee(employee)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// call DB layer
	err = updateEmployee(h.db, employee)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			http.Error(w, "Employee does not exist.",
				http.StatusNotFound)
			return
		}
		http.Error(w, "Error while updating employee "+
			err.Error(), http.StatusInternalServerError)
		return
	}

	// Send Response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) deleteEmployeeHandler(w http.ResponseWriter, r *http.Request) {
	// Parse Request
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Error parsing the ID, make sure it is an integer. Error: "+err.Error(),
			http.StatusBadRequest)
		return
	}

	// call DB layer
	err = deleteEmployee(h.db, id)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			http.Error(w, "Employee does not exist.",
				http.StatusNotFound)
			return
		}
		http.Error(w, "Error while deleting employee "+
			err.Error(), http.StatusInternalServerError)
		return
	}

	// Send Response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) getEmployeesListHandler(w http.ResponseWriter, r *http.Request) {
	// Parse Request
	// Both are optional fields and if not present we default to page 1, size 20
	pageNum := r.URL.Query().Get("page")
	pageSize := r.URL.Query().Get("size")

	page, err := strconv.Atoi(pageNum)
	if err != nil {
		page = 1 // default page
	}

	size, err := strconv.Atoi(pageSize)
	if err != nil || size < 1 {
		size = 10 // default page size
	}

	offset := (page - 1) * size

	// call DB layer
	employees, err := getEmployeesList(h.db, size, offset)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			http.Error(w, "There are no employees.",
				http.StatusNotFound)
			return
		}
		http.Error(w, "Error while listing employee "+
			err.Error(), http.StatusInternalServerError)
		return
	}

	// Send Response
	json.NewEncoder(w).Encode(employees)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
}

// Validate the employee object to make sure all the fields are present
func validateEmployee(emp Employee) error {
	if emp.ID == 0 {
		return errors.New("Employee ID cannot be 0")
	}
	if emp.Name == "" {
		return errors.New("Employee Name cannot be blank")
	}
	if emp.Position == "" {
		return errors.New("Employee Position cannot be blank")
	}
	if emp.Salary == 0 {
		return errors.New("Employee Salary cannot be 0")
	}

	return nil
}
