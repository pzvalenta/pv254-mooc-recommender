package db

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DB ...
var DB *mongo.Database

// GetClient ...
func GetClient() *mongo.Client {
	host := "localhost" // default host, used oudside of container
	if envHost, ok := os.LookupEnv("DB_HOST"); ok {
		host = envHost
	}

	clientOptions := options.Client().ApplyURI("mongodb://" + host + ":27017")
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
