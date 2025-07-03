package main

import (
	"log"

	"go-backend-todo/internal/config"
	"go-backend-todo/internal/db"
	"go-backend-todo/internal/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Create a new Fiber app with configuration
	app := fiber.New(config.GetFiberConfig(cfg))

	// Connect to database using connection pool (RECOMMENDED for production)
	pool, err := db.ConnectPoolWithConfig(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer pool.Close()

	// Setup routes with configuration and database pool
	routes.SetupRoutes(app, cfg, pool)

	// Start server
	serverAddr := config.GetServerAddress(cfg)
	log.Printf("Server starting on %s...", serverAddr)
	log.Fatal(app.Listen(serverAddr))
}
