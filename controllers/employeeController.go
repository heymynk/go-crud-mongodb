package controllers

import (
	"context"
	"encoding/json"
	"go-crud-mongodb/config"
	"go-crud-mongodb/models"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var employeeCollection *mongo.Collection

// Initialize the MongoDB client and set the employeeCollection
func init() {
	client := config.ConnectDB()
	employeeCollection = config.GetCollection(client, "employees")
}

// AddEmployee handles creating a new employee
func AddEmployee(w http.ResponseWriter, r *http.Request) {
	var employee models.Employee
	if err := json.NewDecoder(r.Body).Decode(&employee); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	employee.ID = primitive.NewObjectID()
	_, err := employeeCollection.InsertOne(context.Background(), employee)
	if err != nil {
		http.Error(w, "Failed to add employee", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(employee)
}

// GetEmployees handles fetching all employees
func GetEmployees(w http.ResponseWriter, r *http.Request) {
	cursor, err := employeeCollection.Find(context.Background(), bson.M{})
	if err != nil {
		http.Error(w, "Failed to retrieve employees", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	var employees []models.Employee
	if err = cursor.All(context.Background(), &employees); err != nil {
		http.Error(w, "Failed to parse employees", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(employees)
}

// GetEmployeeByID handles fetching a specific employee by ID
func GetEmployeeByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	employeeID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid employee ID", http.StatusBadRequest)
		return
	}

	var employee models.Employee
	err = employeeCollection.FindOne(context.Background(), bson.M{"_id": employeeID}).Decode(&employee)
	if err != nil {
		http.Error(w, "Employee not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(employee)
}

// UpdateEmployee handles updating an existing employee
func UpdateEmployee(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	employeeID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid employee ID", http.StatusBadRequest)
		return
	}

	var updatedData models.Employee
	if err := json.NewDecoder(r.Body).Decode(&updatedData); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	update := bson.M{
		"$set": bson.M{
			"first_name":    updatedData.FirstName,
			"last_name":     updatedData.LastName,
			"position":      updatedData.Position,
			"salary":        updatedData.Salary,
			"full_time":     updatedData.FullTime,
			"department_id": updatedData.DepartmentID,
		},
	}

	_, err = employeeCollection.UpdateOne(context.Background(), bson.M{"_id": employeeID}, update)
	if err != nil {
		http.Error(w, "Failed to update employee", http.StatusInternalServerError)
		return
	}

	var updatedEmployee models.Employee
	err = employeeCollection.FindOne(context.Background(), bson.M{"_id": employeeID}).Decode(&updatedEmployee)
	if err != nil {
		http.Error(w, "Failed to retrieve updated employee", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(updatedEmployee)
}

// DeleteEmployee handles deleting an employee
func DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	employeeID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid employee ID", http.StatusBadRequest)
		return
	}

	_, err = employeeCollection.DeleteOne(context.Background(), bson.M{"_id": employeeID})
	if err != nil {
		http.Error(w, "Failed to delete employee", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
