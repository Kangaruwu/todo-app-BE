package routes

import (
	"go-backend-todo/internal/api/handlers"
	"go-backend-todo/internal/api/middlewares"
	"go-backend-todo/internal/config"
	auth_repository "go-backend-todo/internal/repository/auth"
	todo_repository "go-backend-todo/internal/repository/todo"
	user_repository "go-backend-todo/internal/repository/user"
	"go-backend-todo/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/jackc/pgx/v5/pgxpool"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

// SetupRoutes sets up all routes for the application
func SetupRoutes(app *fiber.App, cfg *config.Config, pool *pgxpool.Pool) {
	// Global middleware
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New(config.GetCORSConfig(cfg)))
	jwtManager := middlewares.NewJWTManager(cfg)

	// Initialize repositories
	todoRepo := todo_repository.NewTodoRepository(pool)
	userRepo := user_repository.NewUserRepository(pool)
	authRepo := auth_repository.NewAuthRepository(pool)

	// Initialize services
	todoService := service.NewTodoService(todoRepo)
	userService := service.NewUserService(userRepo)
	authService := service.NewAuthService(userRepo, authRepo)

	// Initialize handlers
	todoHandler := handlers.NewTodoHandler(todoService)
	userHandler := handlers.NewUserHandler(userService)
	authHandler := handlers.NewAuthHandler(authService, jwtManager)

	// API routes
	setupAPIRoutes(app, todoHandler, userHandler, authHandler)

	// Swagger documentation
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	// 404 handler
	app.Use(middlewares.NotFound)
}

// setupAPIRoutes sets up API routes with handlers
func setupAPIRoutes(
	app *fiber.App,
	todoHandler *handlers.TodoHandler,
	userHandler *handlers.UserHandler,
	authHandler *handlers.AuthHandler,
) {
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Server is running",
		})
	})
	// API group v1
	api := app.Group("/api/v1")

	// Setup routes with dependency injection
	setupTodoRoutes(api, todoHandler)
	setupUserRoutes(api, userHandler)
	setupAuthRoutes(api, authHandler)
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

// setupUserRoutes sets up user-related routes with dependency injection
func setupUserRoutes(api fiber.Router, userHandler *handlers.UserHandler) {
	users := api.Group("/users")
	users.Get("/profile", userHandler.GetUser)
	users.Put("/profile", userHandler.UpdateUser)
	users.Delete("/profile", userHandler.DeleteUser)
}

// setupAuthRoutes sets up authentication-related routes with dependency injection
func setupAuthRoutes(api fiber.Router, authHandler *handlers.AuthHandler) {
	auth := api.Group("/auth")
	auth.Post("/login", authHandler.Login)
	auth.Post("/register", authHandler.Register)
	auth.Post("/recover-password", authHandler.RecoverPassword)
	auth.Get("/reset-password", authHandler.ResetPassword)
	auth.Get("/confirm-email/:token", authHandler.ConfirmEmail)
}
