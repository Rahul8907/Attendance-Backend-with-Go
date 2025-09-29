package pkg

import (
	"attandance/models"
	"log"
	"net/http"

	"encoding/json"
	"fmt"
	"os"
	"sync"
)

const (
	EmployeeFile string = "employees.json"
)

var (
	fileLock sync.RWMutex = sync.RWMutex{}
)

/*
Adds employee to file.
*/
func AddEmployee(emp *models.Employee) error {
	fileLock.Lock()
	defer fileLock.Unlock()
	rawContents, err := os.ReadFile(EmployeeFile)
	if err != nil {
		return fmt.Errorf("error while reading employee file, err : %s", err.Error())
	}
	emps := make([]models.Employee, 0)
	err = json.Unmarshal(rawContents, &emps)
	if err != nil {
		return fmt.Errorf("error while unmarshalling employee file contents, err : %s", err.Error())
	}
	emps = append(emps, *emp)

	fc, err := json.Marshal(emps)
	if err != nil {
		return fmt.Errorf("error while marshalling employee file contents, err : %s", err.Error())
	}
	err = os.WriteFile(EmployeeFile, fc, os.FileMode(0666))
	if err != nil {
		return fmt.Errorf("error while writing employee file contents, err : %s", err.Error())
	}
	return nil
}

// find employee by Id
func FindEmployee(id string) (models.Employee, error) {
	fileLock.RLock()
	defer fileLock.RUnlock()
	rawContents, err := os.ReadFile(EmployeeFile)
	if err != nil {
		return models.Employee{}, fmt.Errorf("error while reading employee file: %s", err.Error())
	}

	emps := make([]models.Employee, 0)
	err = json.Unmarshal(rawContents, &emps)
	if err != nil {
		return models.Employee{}, fmt.Errorf("error while unmarshalling employee file contents: %s", err.Error())
	}

	for _, emp := range emps {
		if emp.ID == id {
			return emp, nil
		}
	}

	return models.Employee{}, fmt.Errorf("employee with id %s not found", id)
}

// Delete employee
func DeleteOps(res http.ResponseWriter, empID string) {
	fileLock.Lock()
	defer fileLock.Unlock()
	employeesData, err := os.ReadFile(EmployeeFile)
	if err != nil {
		log.Println("ERROR", "Error reading file", err.Error())
		res.WriteHeader(http.StatusInternalServerError)
		res.Write(models.NewAPIError(http.StatusInternalServerError, "Error reading file").Jsonify())
		return
	}
	var employees = make([]models.Employee, 0)

	if err := json.Unmarshal(employeesData, &employees); err != nil {
		log.Println("ERROR", "Error parsing employee data", err.Error())
		res.WriteHeader(http.StatusInternalServerError)
		res.Write(models.NewAPIError(http.StatusInternalServerError, err.Error()).Jsonify())
		return
	}
	updatedEmployees := make([]models.Employee, 0)
	found := false
	for _, emp := range employees {
		if emp.ID != empID {
			updatedEmployees = append(updatedEmployees, emp)
		} else {
			found = true
		}
	}
	if !found {
		log.Println("ERROR", "Employee not found", err.Error())
		res.WriteHeader(http.StatusNotFound)
		res.Write(models.NewAPIError(http.StatusNotFound, err.Error()).Jsonify())
		return
	}
	updatedData, err := json.Marshal(updatedEmployees)
	if err != nil {
		log.Println("ERROR", "Error writing updated data", err.Error())
		res.WriteHeader(http.StatusInternalServerError)
		res.Write(models.NewAPIError(http.StatusInternalServerError, err.Error()).Jsonify())
		return
	}
	if err := os.WriteFile(EmployeeFile, updatedData, 0644); err != nil {
		log.Println("ERROR", "Error saving updated file", err.Error())
		res.WriteHeader(http.StatusInternalServerError)
		res.Write(models.NewAPIError(http.StatusInternalServerError, err.Error()).Jsonify())
		return
	}
	res.WriteHeader(http.StatusOK)
}

//func UpdateOps()
