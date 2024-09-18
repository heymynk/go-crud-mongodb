package main

import (
	"go-crud-mongodb/config"
	"go-crud-mongodb/routes"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// Initialize MongoDB connection
	client := config.ConnectDB()
	defer client.Disconnect(nil)

	// Initialize Gorilla Mux router
	router := mux.NewRouter()

	// Register routes
	routes.RegisterEmployeeRoutes(router)

	// Start the server
	log.Println("Server started on port 8080")
	if err := http.ListenAndServe(":8081", router); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
