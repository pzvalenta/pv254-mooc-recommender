package db

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DB ...
var DB *mongo.Database

// GetClient ...
func GetClient() *mongo.Client {
	clientOptions := options.Client().ApplyURI("mongodb://mongo_dev:27017")
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Connect(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	DB = client.Database("mydb")
	return client
}

// GetDB ...
func GetDB() *mongo.Database {
	return DB
}
