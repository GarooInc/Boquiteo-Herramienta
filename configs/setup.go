package configs

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

// DB Client instance
var DB *mongo.Client

// ConnectDB connects to the database
func ConnectDB() *mongo.Client {
	mongoURI, err := LoadMongoEnv()

	if err != nil {
		log.Fatal(err)
	}

	// Set client options
	clientOptions := options.Client().ApplyURI(mongoURI)

	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)

	if err != nil {
		log.Fatal(err)

	}

	//ping the database
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB")
	return client
}

func LoadSetup() {
	DB = ConnectDB()
}

// GetCollection getting database collections
func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	collection := client.Database("Boquiteo").Collection(collectionName)
	return collection
}
