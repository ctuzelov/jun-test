package repository

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Function that handles mongoDB connection
func DBinstance() *mongo.Client {
	err := godotenv.Load("./cmd/.env")
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		return nil
	}
	MongoDb := os.Getenv("MONGODB_URL")

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

	defer cancel()
	opts := options.Client().ApplyURI(MongoDb)

	client, err := mongo.Connect(ctx, opts)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")
	return client
}

var Client *mongo.Client = DBinstance()

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	var collection *mongo.Collection = client.Database("test").Collection(collectionName)
	return collection
}
