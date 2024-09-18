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

var departmentCollection *mongo.Collection

// Initialize the MongoDB client and set the departmentCollection
func init() {
	client := config.ConnectDB()
	departmentCollection = config.GetCollection(client, "departments")
}

// AddDepartment handles creating a new department
func AddDepartment(w http.ResponseWriter, r *http.Request) {
	var department models.Department
	if err := json.NewDecoder(r.Body).Decode(&department); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	department.ID = primitive.NewObjectID()
	_, err := departmentCollection.InsertOne(context.Background(), department)
	if err != nil {
		http.Error(w, "Failed to add department", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(department)
}

// GetDepartments handles fetching all departments
func GetDepartments(w http.ResponseWriter, r *http.Request) {
	cursor, err := departmentCollection.Find(context.Background(), bson.M{})
	if err != nil {
		http.Error(w, "Failed to retrieve departments", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	var departments []models.Department
	if err = cursor.All(context.Background(), &departments); err != nil {
		http.Error(w, "Failed to parse departments", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(departments)
}

// GetDepartmentByID handles fetching a specific department by ID
func GetDepartmentByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	departmentID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid department ID", http.StatusBadRequest)
		return
	}

	var department models.Department
	err = departmentCollection.FindOne(context.Background(), bson.M{"_id": departmentID}).Decode(&department)
	if err != nil {
		http.Error(w, "Department not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(department)
}

// UpdateDepartment handles updating an existing department
func UpdateDepartment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	departmentID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid department ID", http.StatusBadRequest)
		return
	}

	var updatedData models.Department
	if err := json.NewDecoder(r.Body).Decode(&updatedData); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	update := bson.M{
		"$set": bson.M{
			"name":     updatedData.Name,
			"location": updatedData.Location,
		},
	}

	_, err = departmentCollection.UpdateOne(context.Background(), bson.M{"_id": departmentID}, update)
	if err != nil {
		http.Error(w, "Failed to update department", http.StatusInternalServerError)
		return
	}

	var updatedDepartment models.Department
	err = departmentCollection.FindOne(context.Background(), bson.M{"_id": departmentID}).Decode(&updatedDepartment)
	if err != nil {
		http.Error(w, "Failed to retrieve updated department", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(updatedDepartment)
}

// DeleteDepartment handles deleting a department
func DeleteDepartment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	departmentID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid department ID", http.StatusBadRequest)
		return
	}

	_, err = departmentCollection.DeleteOne(context.Background(), bson.M{"_id": departmentID})
	if err != nil {
		http.Error(w, "Failed to delete department", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
