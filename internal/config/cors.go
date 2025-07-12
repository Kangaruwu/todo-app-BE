package config

import (
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// GetCORSConfigFromConfig returns CORS configuration for Fiber from Config struct
func GetCORSConfigFromConfig(cfg *Config) cors.Config {
	return cors.Config{
		AllowOrigins:     cfg.CORS.AllowOrigins,
		AllowMethods:     cfg.CORS.AllowMethods,
		AllowHeaders:     cfg.CORS.AllowHeaders,
		AllowCredentials: cfg.CORS.AllowCredentials,
		ExposeHeaders:    "Content-Length,Content-Type",
	}
}
