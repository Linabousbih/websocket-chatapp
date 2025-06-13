package auth

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB() *mongo.Client {
	log.Println("Attempting to connect to MongoDB...")
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("MongoDB is not reachable: %v", err)
	}

	log.Println("Successfully connected to MongoDB!")
	return client
}

var Client *mongo.Client = ConnectDB()

func OpenCollection(collectionName string) *mongo.Collection {

	if Client == nil {
		log.Fatal("MongoDB client is not initialized. Please call ConnectDB first.")
	}
	return Client.Database("chatusersdb").Collection(collectionName)
}