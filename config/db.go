package config

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

// ConnectDB function will initialize the mongoDB client
func ConnectDB() *mongo.Client {

	//load env variables from env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	//get the mongoDB URI from the env variables
	MONGODB_URI := os.Getenv("MONGODB_URI")
	if MONGODB_URI == "" {
		log.Fatal("MONGODB_URI is not set")
	}

	//set MongoDB Client options
	clientOptions := options.Client().ApplyURI(MONGODB_URI)

	// Create a context with timeout for connecting to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//Create and connect the mongodb client
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("failed to create and connect mongoDb client", err)
	}

	//check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Failed to ping mongodb", err)
	}

	log.Println("Connected to mongoDB")
	Client = client
	return client

}

// GetCollection retrieves a MongoDB collection
func GetCollection(Client *mongo.Client, collectionName string) *mongo.Collection {
	return Client.Database("go-crud-mongodb").Collection(collectionName)
}
