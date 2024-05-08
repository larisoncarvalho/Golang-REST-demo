package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

// CREATE TABLE IF NOT EXISTS employees (
//
//	ID INTEGER PRIMARY KEY,
//	Name TEXT,
//	Position TEXT,
//	Salary REAL
//
// );
// Employee Struct:
type Employee struct {
	//Unique identifier for the employee.
	ID int `json:"id"`
	//Name of the employee.
	Name string `json:"name"`
	//Position/title of the employee.
	Position string `json:"position"`
	//Salary of the employee.
	Salary float64 `json:"salary"`
}

// Insert the employee
func createEmployee(db *sql.DB, emp Employee) error {
	_, err := db.Exec("INSERT INTO employees (id, name, position, salary) VALUES (?, ?, ?, ?)",
		emp.ID, emp.Name, emp.Position, emp.Salary)
	return err
}

// Update the employee
func updateEmployee(db *sql.DB, emp Employee) error {
	// Check if employee with this ID exists
	_, err := getEmployeeById(db, emp.ID)
	if err != nil {
		return err
	}
	_, err = db.Exec("UPDATE employees set name = ?, position = ?, salary = ? where id = ?",
		emp.Name, emp.Position, emp.Salary, emp.ID)
	return err
}

// Delete the employee
func deleteEmployee(db *sql.DB, id int) error {
	// Check if employee with this ID exists
	_, err := getEmployeeById(db, id)
	if err != nil {
		return err
	}
	_, err = db.Exec("DELETE from employees where id = ?", id)
	return err
}

// Get employee by Id
func getEmployeeById(db *sql.DB, id int) (Employee, error) {
	var employee Employee
	row := db.QueryRow("SELECT id, name, position, salary from employees where id = ?", id)
	err := row.Scan(&employee.ID, &employee.Name, &employee.Position, &employee.Salary)
	if err != nil {
		return Employee{}, err
	}
	return employee, nil
}

// List the employees
func getEmployeesList(db *sql.DB, size int, offset int) ([]Employee, error) {
	var employees []Employee
	rows, err := db.Query("SELECT id, name, position, salary FROM employees ORDER BY ID asc LIMIT ? OFFSET ? ", size, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var employee Employee
		if err := rows.Scan(&employee.ID, &employee.Name, &employee.Position, &employee.Salary); err != nil {
			return nil, err
		}
		employees = append(employees, employee)
	}

	return employees, nil
}
