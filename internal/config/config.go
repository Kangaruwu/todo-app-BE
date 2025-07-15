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
	JWT      JWTConfig
	Email    EmailConfig
	CORS     CORSConfig
	Token    TokenConfig
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

// JWTConfig holds JWT configuration
type JWTConfig struct {
	AccessSecret     	string
	RefreshSecret    	string
	VerificationSecret 	string
	RecoverySecret     	string
	AccessExpiryHour 	int
	RefreshExpiryDay 	int
}

// EmailConfig holds email configuration
type EmailConfig struct {
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
	FromEmail    string
	FromName     string
}

// CORSConfig holds CORS configuration
type CORSConfig struct {
	AllowOrigins     string
	AllowMethods     string
	AllowHeaders     string
	AllowCredentials bool
}

// TokenConfig holds token configuration
type TokenConfig struct {
	VerifyEmailTokenSecret string
	VerifyEmailTokenTTL    int // in minutes
	RecoverPasswordTokenSecret string
	RecoverPasswordTokenTTL    int // in minutes
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
		JWT: JWTConfig{
			AccessSecret:      	GetEnv("JWT_ACCESS_SECRET", "your-super-secret-access-key"),
			RefreshSecret:     	GetEnv("JWT_REFRESH_SECRET", "your-super-secret-refresh-key"),
			VerificationSecret: GetEnv("JWT_VERIFICATION_SECRET", "your-super-secret-verification-key"),
			RecoverySecret:     GetEnv("JWT_RECOVERY_SECRET", "your-super-secret-recovery-key"),
			AccessExpiryHour: getEnvAsInt("JWT_ACCESS_EXPIRY_HOUR", 24),
			RefreshExpiryDay: getEnvAsInt("JWT_REFRESH_EXPIRY_DAY", 7),
		},
		Email: EmailConfig{
			SMTPHost:     GetEnv("SMTP_HOST", "smtp.gmail.com"),
			SMTPPort:     getEnvAsInt("SMTP_PORT", 587),
			SMTPUsername: GetEnv("SMTP_USERNAME", ""),
			SMTPPassword: GetEnv("SMTP_PASSWORD", ""),
			FromEmail:    GetEnv("FROM_EMAIL", "noreply@todoapp.com"),
			FromName:     GetEnv("FROM_NAME", "Todo App"),
		},
		CORS: CORSConfig{
			AllowOrigins:     GetEnv("CORS_ALLOW_ORIGINS", "*"), // Allow all origins 
			AllowMethods:     GetEnv("CORS_ALLOW_METHODS", "GET,POST,PUT,DELETE,OPTIONS,PATCH"),
			AllowHeaders:     GetEnv("CORS_ALLOW_HEADERS", "Content-Type,Authorization,Accept,Origin,X-Requested-With"),
			AllowCredentials: getEnvAsBool("CORS_ALLOW_CREDENTIALS", false), // Must be false when AllowOrigins is "*"
		},
		Token: TokenConfig{
			VerifyEmailTokenSecret: GetEnv("VERIFY_EMAIL_TOKEN_SECRET", "your-super-secret-verify-email-key"),
			VerifyEmailTokenTTL:    getEnvAsInt("VERIFY_EMAIL_TOKEN_TTL_MINUTES", 30), // Default 30 minutes
			RecoverPasswordTokenSecret: GetEnv("RECOVER_PASSWORD_TOKEN_SECRET", "your-super-secret-recover-password-key"),
			RecoverPasswordTokenTTL:    getEnvAsInt("RECOVER_PASSWORD_TOKEN_TTL_MINUTES", 30), // Default 30 minutes
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
