package pkg

import (
	"attandance/models"
	"errors"
	"net/http"
	"time"

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
func DeleteEmployee(id string) (int, error) {

	fileLock.Lock()
	defer fileLock.Unlock()

	/* logic :
	1.

	*/
	employeesData, err := os.ReadFile(EmployeeFile)
	if err != nil {
		return http.StatusInternalServerError, errors.New("Error reading file")
	}
	var employees = make([]models.Employee, 0)
	if err := json.Unmarshal(employeesData, &employees); err != nil {
		return http.StatusInternalServerError, errors.New("Error parsing employee data")
	}
	updatedEmployees := make([]models.Employee, 0)
	found := false
	for _, emp := range employees {
		if emp.ID != id {
			updatedEmployees = append(updatedEmployees, emp)
		} else {
			found = true
		}
	}
	if !found {
		return http.StatusNotFound, errors.New("employee not found")
	}
	updatedData, err := json.Marshal(updatedEmployees)
	if err != nil {
		return http.StatusInternalServerError, errors.New("error writing updated data")
	}
	if err := os.WriteFile(EmployeeFile, updatedData, 0644); err != nil {

		return http.StatusInternalServerError, errors.New("error saving updated file")
	}
	return http.StatusOK, nil
}

// Update employee feilds
func UpdateOps(empID string, incomingType string) (error, int) {
	fileLock.Lock()
	defer fileLock.Unlock()
	if incomingType != "login" && incomingType != "logout" {
		return errors.New("icoming type is not correct"), http.StatusInternalServerError
	}
	if incomingType == "login" {
		employeeData, err := os.ReadFile(EmployeeFile)
		if err != nil {
			return errors.New("error reading file"), http.StatusInternalServerError
		}
		var employees = make([]models.Employee, 0)
		err = json.Unmarshal(employeeData, &employees)
		if err != nil {
			return errors.New("Error unmarshalling data"), http.StatusInternalServerError
		}
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
			return errors.New("Employee not found"), http.StatusNotFound
		}
		updatedData, err := json.Marshal(employees)
		if err != nil {
			return errors.New("Error marshalling updated data "), http.StatusInternalServerError
		}

		// Write back to file
		err = os.WriteFile(EmployeeFile, updatedData, 0666)
		if err != nil {
			return errors.New("Error writing to file"), http.StatusInternalServerError
		}
	} else {
		employeeData, err := os.ReadFile(EmployeeFile)
		if err != nil {
			return errors.New("error reading file"), http.StatusInternalServerError
		}
		var employees = make([]models.Employee, 0)
		err = json.Unmarshal(employeeData, &employees)
		if err != nil {
			return errors.New("Error unmarshalling data"), http.StatusInternalServerError
		}

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
			return errors.New("Employee not found"), http.StatusNotFound
		}
		updatedData, err := json.Marshal(employees)
		if err != nil {
			return errors.New("Error marshalling updated data "), http.StatusInternalServerError
		}

		// Write back to file
		err = os.WriteFile(EmployeeFile, updatedData, 0666)
		if err != nil {
			return errors.New("Error writing to file"), http.StatusInternalServerError
		}
	}
	return nil, http.StatusOK
}
