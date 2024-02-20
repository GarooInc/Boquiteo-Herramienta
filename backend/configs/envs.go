package configs

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func LoadMongoEnv() (string, error) {
	fmt.Println("Loading from system environment")
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI != "" {
		return mongoURI, nil

	}

	fmt.Println("Loading .env file")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	mongoURI = os.Getenv("MONGO_URI")
	if mongoURI == "" {
		return "", fmt.Errorf("MONGO_URI not found")
	}

	return mongoURI, nil
}
