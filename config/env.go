package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func EnvMongoURI() string {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env")
	}

	monguURI := os.Getenv("MONGOURI")
	return monguURI
}
