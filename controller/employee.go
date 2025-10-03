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
		return
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
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(models.NewAPIError(http.StatusInternalServerError, "Invalid JSON in file").Jsonify())
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	empID := r.PathValue("id")
	code, err := pkg.UpdateOps(empID, "login")
	if err != nil {
		w.WriteHeader(code)
		w.Write(models.NewAPIError(code, err.Error()).Jsonify())
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Login time updated successfully"))
}

func LogOutHandler(w http.ResponseWriter, r *http.Request) {
	empID := r.PathValue("id")

	// Read employee data from file
	code, err := pkg.UpdateOps(empID, "logout")
	if err != nil {
		w.WriteHeader(code)
		w.Write(models.NewAPIError(code, err.Error()).Jsonify())
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Logout time updated successfully"))
}
func DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	empID := r.PathValue("id")

	if empID == "" {
		// return error here
		w.WriteHeader(http.StatusBadRequest)
		w.Write(models.NewAPIError(http.StatusBadRequest, "Please pass id of the employee").Jsonify())
	}

	statusCode, err := pkg.DeleteEmployee(empID)
	if err != nil {
		log.Println("ERROR", err.Error())
		w.WriteHeader(statusCode)
		w.Write(models.NewAPIError(statusCode, err.Error()).Jsonify())
		return
		// return error
	}
	w.WriteHeader(statusCode)
}
