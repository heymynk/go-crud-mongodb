package routes

import (
	"go-crud-mongodb/controllers"
	"log"

	"github.com/gorilla/mux"
)

func RegisterEmployeeRoutes(router *mux.Router) {
	log.Println("Registering employee routes")

	router.HandleFunc("/employees", controllers.AddEmployee).Methods("POST")
	router.HandleFunc("/employees", controllers.GetEmployees).Methods("GET")
	router.HandleFunc("/employees/{id}", controllers.UpdateEmployee).Methods("PUT")
	router.HandleFunc("/employees/{id}", controllers.DeleteEmployee).Methods("DELETE")

	router.HandleFunc("/departments", controllers.AddDepartment).Methods("POST")
	router.HandleFunc("/departments", controllers.GetDepartments).Methods("GET")
	router.HandleFunc("/departments/{id}", controllers.GetDepartmentByID).Methods("GET")
	router.HandleFunc("/departments/{id}", controllers.UpdateDepartment).Methods("PUT")
	router.HandleFunc("/departments/{id}", controllers.DeleteDepartment).Methods("DELETE")
}
