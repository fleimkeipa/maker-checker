package pkg

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func MongoConnect() (*mongo.Database, error) {
	uri := "mongodb://localhost:27017"
	if isStageContainer() {
		fmt.Println("Program is running inside a container.")
		uri = "mongodb://mongodb:27017"
	}

	// Set client options
	clientOptions := options.Client().ApplyURI(uri)

	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}

	// Check the connection
	if err := client.Ping(context.Background(), nil); err != nil {
		return nil, err
	}

	// Set the database and collection variables
	return client.Database("maker-checker"), nil
}

func isStageContainer() bool {
	_, isContainer := os.LookupEnv("HOSTNAME")
	return isContainer
}
