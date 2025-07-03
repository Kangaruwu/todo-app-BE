package routes

import (
	"go-backend-todo/internal/api/handlers"
	"go-backend-todo/internal/api/middlewares"
	"go-backend-todo/internal/config"
	"go-backend-todo/internal/repository"
	"go-backend-todo/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/jackc/pgx/v5/pgxpool"
)

// SetupRoutes sets up all routes for the application
func SetupRoutes(app *fiber.App, cfg *config.Config, pool *pgxpool.Pool) {
	// Global middleware
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New(config.GetCORSConfig(cfg)))

	// Initialize repositories
	todoRepo := repository.NewTodoRepository(pool)
	// userRepo := repository.NewUserRepository(pool) // TODO: implement

	// Initialize services
	todoService := service.NewTodoService(todoRepo)
	// userService := service.NewUserService(userRepo) // TODO: implement

	// Initialize handlers
	todoHandler := handlers.NewTodoHandler(todoService)
	// userHandler := handlers.NewUserHandler(userService) // TODO: implement

	// API routes
	setupAPIRoutes(app, todoHandler)

	// 404 handler
	app.Use(middlewares.NotFound)
}

// setupAPIRoutes sets up API routes with handlers
func setupAPIRoutes(app *fiber.App, todoHandler *handlers.TodoHandler) {
	// API group v1
	api := app.Group("/api/v1")

	// Health check
	api.Get("/health", handlers.Hello)

	// Setup todo routes with dependency injection
	setupTodoRoutes(api, todoHandler)

	// Setup user routes (legacy, no DI)
	setupUserRoutes(api)

	// setupAuthRoutes(api) // TODO: Uncomment when finish auth handlers
}

// setupTodoRoutes sets up todo-related routes with dependency injection
func setupTodoRoutes(api fiber.Router, todoHandler *handlers.TodoHandler) {
	todos := api.Group("/todos")
	todos.Get("/", todoHandler.GetTodos)
	todos.Post("/", todoHandler.CreateTodo)
	todos.Get("/:id", todoHandler.GetTodo)
	todos.Put("/:id", todoHandler.UpdateTodo)
	todos.Delete("/:id", todoHandler.DeleteTodo)
}

// setupUserRoutes sets up user-related routes (legacy implementation)
func setupUserRoutes(api fiber.Router) {
	users := api.Group("/users")
	users.Get("/", handlers.GetUsers)
	users.Post("/", handlers.CreateUser)
	users.Get("/:id", handlers.GetUser)
	users.Put("/:id", handlers.UpdateUser)
	users.Delete("/:id", handlers.DeleteUser)
}


// setupAuthRoutes sets up authentication-related routes
// func setupAuthRoutes(api fiber.Router) {
//     auth := api.Group("/auth")
//     auth.Post("/login", handlers.Login)
//     auth.Post("/register", handlers.Register)
//     auth.Post("/logout", handlers.Logout)
//     auth.Get("/me", handlers.GetProfile)
// }
