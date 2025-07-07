package routes

import (
	"go-backend-todo/internal/api/handlers"
	"go-backend-todo/internal/api/middlewares"
	"go-backend-todo/internal/config"
	"go-backend-todo/internal/repository/todo"
	"go-backend-todo/internal/repository/user"
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
	todoRepo := todo_repository.NewTodoRepository(pool)
	userRepo := user_repository.NewUserRepository(pool) 

	// Initialize services
	todoService := service.NewTodoService(todoRepo)
	userService := service.NewUserService(userRepo) 

	// Initialize handlers
	todoHandler := handlers.NewTodoHandler(todoService)
	userHandler := handlers.NewUserHandler(userService) 

	// API routes
	setupAPIRoutes(app, todoHandler, userHandler)

	// 404 handler
	app.Use(middlewares.NotFound)
}

// setupAPIRoutes sets up API routes with handlers
func setupAPIRoutes(
	app *fiber.App, 
	todoHandler *handlers.TodoHandler,
	userHandler *handlers.UserHandler) {
	// API group v1
	api := app.Group("/api/v1")

	// Health check
	api.Get("/health", handlers.Hello)

	// Setup routes with dependency injection
	setupTodoRoutes(api, todoHandler)
	setupUserRoutes(api, userHandler)

	// setupAuthRoutes(api) // TODO: 
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
func setupUserRoutes(api fiber.Router, userHandler *handlers.UserHandler) {
	users := api.Group("/users")
	users.Get("/", userHandler.GetUsers)
	users.Post("/", userHandler.CreateUser)
	users.Get("/:id", userHandler.GetUser)
	users.Put("/:id", userHandler.UpdateUser)
	users.Delete("/:id", userHandler.DeleteUser)
}


// setupAuthRoutes sets up authentication-related routes
// func setupAuthRoutes(api fiber.Router) {
//     auth := api.Group("/auth")
//     auth.Post("/login", handlers.Login)
//     auth.Post("/register", handlers.Register)
//     auth.Post("/logout", handlers.Logout)
//     auth.Get("/me", handlers.GetProfile)
// }
