package config

import (
	"os"
	"strconv"
)

// Config holds all configuration for our application
type Config struct {
	App      AppConfig
	Database DatabaseConfig
	Server   ServerConfig
}

// AppConfig holds application-specific configuration
type AppConfig struct {
	Name        string
	Environment string
	Debug       bool
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
	ChannelBinding string 
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Host string
	Port string
}

// Load loads configuration from environment variables
func Load() *Config {
	// Load environment variables from .env file
	LoadEnv()

	return &Config{
		App: AppConfig{
			Name:        GetEnv("APP_NAME", "Go Backend Todo API"),
			Environment: GetEnv("APP_ENV", "development"),
			Debug:       getEnvAsBool("APP_DEBUG", true),
		},
		Database: DatabaseConfig{
			Host:     		GetEnv("DB_HOST", "localhost"),
			Port:     		getEnvAsInt("DB_PORT", 5432),
			User:     		GetEnv("DB_USER", "postgres"),
			Password: 		GetEnv("DB_PASSWORD", "password"),
			DBName:  		GetEnv("DB_NAME", "todo_db"),
			SSLMode:  		GetEnv("DB_SSLMODE", "disable"),
			ChannelBinding: GetEnv("DB_CHANNEL_BINDING", "prefer"),
		},
		Server: ServerConfig{
			Host: GetEnv("SERVER_HOST", "localhost"),
			Port: GetEnv("SERVER_PORT", "8080"),
		},
	}
}

// Helper functions
func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}
