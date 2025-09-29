package controller

import (
	"attandance/models"
	"attandance/pkg"

	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
)

func CreateEmployeeHandler(w http.ResponseWriter, r *http.Request) {

	rawBody, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("ERROR", "failed to read raw body", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write(models.NewAPIError(http.StatusBadRequest, err.Error()).Jsonify())
		return
	}
	defer r.Body.Close()

	var emp models.Employee
	err = json.Unmarshal(rawBody, &emp)
	if err != nil {
		log.Println("ERROR", "failed to unmarshall body", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write(models.NewAPIError(http.StatusBadRequest, err.Error()).Jsonify())
		return
	}
	// TODO : Set Login time and logout time zero

	// Add a valid UUID to employee
	emp.ID = uuid.NewString()

	if err := emp.Validate(false); err != nil {
		log.Println("ERROR", "failed to read body", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write(models.NewAPIError(http.StatusBadRequest, err.Error()).Jsonify())
		return
	}

	err = pkg.AddEmployee(&emp)
	if err != nil {
		log.Println("ERROR", "failed to Add employee to file", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(models.NewAPIError(http.StatusInternalServerError, err.Error()).Jsonify())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(emp.Jsonify())
}

func GetEmployeeHandler(w http.ResponseWriter, r *http.Request) {
	ID := r.PathValue("id")
	emp, err := pkg.FindEmployee(ID)
	if err != nil {
		http.Error(w, fmt.Sprintf("User not Found :%v", err), http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(emp.Jsonify())
}
func GetAllEmployeeHandler(w http.ResponseWriter, r *http.Request) {
	data, err := os.ReadFile(pkg.EmployeeFile)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading file: %v", err), http.StatusInternalServerError)
		return
	}

	if !json.Valid(data) {
		http.Error(w, "Invalid JSON in file", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	empID := r.PathValue("id")

	// Read employee data from file
	employeeData, err := os.ReadFile(pkg.EmployeeFile)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading file: %v", err), http.StatusInternalServerError)
		return
	}

	var employees = make([]models.Employee, 0)
	err = json.Unmarshal(employeeData, &employees)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error unmarshalling data: %v", err), http.StatusInternalServerError)
		return
	}

	// Flag to check if employee was found
	found := false

	// Update login time
	for i := range employees {
		if employees[i].ID == empID {
			employees[i].LogInTime = time.Now()
			found = true
			break
		}
	}

	if !found {
		http.Error(w, "Employee not found", http.StatusNotFound)
		return
	}

	// Marshal updated data
	updatedData, err := json.Marshal(employees)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error marshalling updated data: %v", err), http.StatusInternalServerError)
		return
	}

	// Write back to file
	err = os.WriteFile(pkg.EmployeeFile, updatedData, 0666)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error writing to file: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Login time updated successfully"))
}

func LogOutHandler(w http.ResponseWriter, r *http.Request) {
	empID := r.PathValue("id")

	// Read employee data from file
	employeeData, err := os.ReadFile(pkg.EmployeeFile)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading file: %v", err), http.StatusInternalServerError)
		return
	}

	var employees = make([]models.Employee, 0)
	err = json.Unmarshal(employeeData, &employees)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error unmarshalling data: %v", err), http.StatusInternalServerError)
		return
	}

	// Flag to check if employee was found
	found := false

	// Update logout time
	for i := range employees {
		if employees[i].ID == empID {
			employees[i].LogOutTime = time.Now()
			found = true
			break
		}
	}

	if !found {
		http.Error(w, "Employee not found", http.StatusNotFound)
		return
	}

	// Marshal updated data
	updatedData, err := json.MarshalIndent(employees, "", "  ")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error marshalling updated data: %v", err), http.StatusInternalServerError)
		return
	}

	// Write back to file
	err = os.WriteFile(pkg.EmployeeFile, updatedData, 0666)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error writing to file: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Logout time updated successfully"))
}
func DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	empID := r.PathValue("id")
	pkg.DeleteOps(w, empID)
}
