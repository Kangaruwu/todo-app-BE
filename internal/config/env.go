package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// LoadEnv loads environment variables from .env file
func LoadEnv() {
	// Try to load .env file, but don't fail if it doesn't exist
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}
}

// IsProduction checks if the application is running in production
func IsProduction() bool {
	return os.Getenv("APP_ENV") == "production"
}

// IsDevelopment checks if the application is running in development
func IsDevelopment() bool {
	env := os.Getenv("APP_ENV")
	return env == "" || env == "development"
}
