package config

import (
	"fmt"

	"github.com/gofiber/fiber/v2/middleware/cors"
)

// GetCORSConfig returns CORS configuration
func GetCORSConfig(cfg *Config) cors.Config {
	if cfg.App.Environment == "production" {
		return cors.Config{
			AllowOrigins:     "https://yourdomain.com,https://www.yourdomain.com",
			AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
			AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
			AllowCredentials: true,
		}
	}

	return cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}
}

// GetDatabaseURL returns formatted database connection string
func GetDatabaseURL(cfg *Config) string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", //" channel_binding=%s"
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
		//cfg.Database.ChannelBinding,
	)
}

// GetServerAddress returns formatted server address
func GetServerAddress(cfg *Config) string {
	if cfg.Server.Host == "" || cfg.Server.Host == "localhost" {
		return ":" + cfg.Server.Port
	}
	return cfg.Server.Host + ":" + cfg.Server.Port
}
