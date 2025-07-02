package main

import (
	"context"
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

	// Connect to database
	db, err := db.Connect()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close(context.Background())

	// Setup routes with configuration
	routes.SetupRoutes(app, cfg)

	// Start server
	serverAddr := config.GetServerAddress(cfg)
	log.Printf("Server starting on %s...", serverAddr)
	log.Fatal(app.Listen(serverAddr))
}
