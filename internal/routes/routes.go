package routes

import (
	"go-backend-todo/internal/api/handlers"
	"go-backend-todo/internal/api/middlewares"
	"go-backend-todo/internal/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

// Setup all routes for the application
func SetupRoutes(app *fiber.App, cfg *config.Config) {
	// Global middleware
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New(config.GetCORSConfig(cfg)))

	// API routes
	setupAPIRoutes(app)

	// 404 handler
	app.Use(middlewares.NotFound)
}

// Setup API routes
func setupAPIRoutes(app *fiber.App) {
	// API group v1
	api := app.Group("/api/v1")

	// Health check
	api.Get("/health", handlers.Hello)

	// Setup user and todo routes
	setupUserRoutes(api)
	setupTodoRoutes(api)
	// setupAuthRoutes(api) // TODO: Uncomment khi có auth handlers
}

func setupUserRoutes(api fiber.Router) {
	users := api.Group("/users")
	users.Get("/", handlers.GetUsers)
	users.Post("/", handlers.CreateUser)
	users.Get("/:id", handlers.GetUser)
	users.Put("/:id", handlers.UpdateUser)
	users.Delete("/:id", handlers.DeleteUser)
}

func setupTodoRoutes(api fiber.Router) {
	todos := api.Group("/todos")
	todos.Get("/", handlers.GetTodos)
	todos.Post("/", handlers.CreateTodo)
	todos.Get("/:id", handlers.GetTodo)
	todos.Put("/:id", handlers.UpdateTodo)
	todos.Delete("/:id", handlers.DeleteTodo)
}

// func setupAuthRoutes(api fiber.Router) {
//     auth := api.Group("/auth")
//     auth.Post("/login", handlers.Login)
//     auth.Post("/register", handlers.Register)
//     auth.Post("/logout", handlers.Logout)
//     auth.Get("/me", handlers.GetProfile)
// }
