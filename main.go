package main

import (
	"attandance/controller"
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /employee", controller.CreateEmployeeHandler)
	mux.HandleFunc("GET /employee/{id}", controller.GetEmployeeHandler)
	mux.HandleFunc("GET /employees", controller.GetAllEmployeeHandler)
	mux.HandleFunc("DELETE /employee/{id}", controller.DeleteEmployee)
	mux.HandleFunc("PATCH /employee/login/{id}", controller.LoginHandler)
	mux.HandleFunc("PATCH /employee/logout/{id}", controller.LogOutHandler)
	fmt.Println("Server running on port 8080...")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}
