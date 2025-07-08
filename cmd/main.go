package main

import (
	"log"

	_ "go-backend-todo/docs" // Import for swagger docs
	"go-backend-todo/internal/config"
	"go-backend-todo/internal/db"
	"go-backend-todo/internal/routes"

	"github.com/gofiber/fiber/v2"
)

// @title Go Backend Todo API
// @version 1.0
// @description This is a sample Todo API server built with Go Fiber
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url https://www.youtube.com/watch?v=KsB99Sf_fX0
// @contact.email https://www.youtube.com/watch?v=KsB99Sf_fX0

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1
// @schemes http https

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

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
