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

	// Initialize repositories
	todoRepo := todo_repository.NewTodoRepository(pool)
	userRepo := user_repository.NewUserRepository(pool)
	authRepo := auth_repository.NewAuthRepository(pool)

	// Initialize JWT manager with userRepo
	jwtManager := middlewares.NewJWTManager(cfg, userRepo)

	// Initialize services
	emailService := service.NewEmailService(cfg)
	todoService := service.NewTodoService(todoRepo)
	userService := service.NewUserService(userRepo)
	authService := service.NewAuthService(userRepo, authRepo, emailService, cfg)

	// Initialize handlers
	todoHandler := handlers.NewTodoHandler(todoService)
	userHandler := handlers.NewUserHandler(userService)
	authHandler := handlers.NewAuthHandler(authService, jwtManager)

	// API routes
	setupAPIRoutes(app, todoHandler, userHandler, authHandler, jwtManager)

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
	jwtManager *middlewares.JWTManager,
) {
	// API group v1
	api := app.Group("/api/v1")

	// Setup routes with dependency injection
	setupTodoRoutes(api, todoHandler, jwtManager)
	setupUserRoutes(api, userHandler, jwtManager)
	setupAuthRoutes(api, authHandler, jwtManager)
}

// setupTodoRoutes sets up todo-related routes with dependency injection
func setupTodoRoutes(api fiber.Router, todoHandler *handlers.TodoHandler, jwtManager *middlewares.JWTManager) {
	todos := api.Group("/todos")

	todos.Use(middlewares.AuthenticateJWT(jwtManager)) 

	todos.Get("/", todoHandler.GetTodos)
	todos.Post("/", todoHandler.CreateTodo)
	todos.Get("/:id", todoHandler.GetTodo)
	todos.Put("/:id", todoHandler.UpdateTodo)
	todos.Delete("/:id", todoHandler.DeleteTodo)
}

// setupUserRoutes sets up user-related routes with dependency injection
func setupUserRoutes(api fiber.Router, userHandler *handlers.UserHandler, jwtManager *middlewares.JWTManager) {
	users := api.Group("/users")

	users.Use(middlewares.AuthenticateJWT(jwtManager)) 

	users.Get("/profile", userHandler.GetUserProfile)
	users.Put("/profile", userHandler.UpdateUserProfile)
	users.Delete("/profile", userHandler.DeleteUserProfile)
	users.Put("/change-password", userHandler.ChangePassword)
}

// setupAuthRoutes sets up authentication-related routes with dependency injection
func setupAuthRoutes(api fiber.Router, authHandler *handlers.AuthHandler, jwtManager *middlewares.JWTManager) {
	auth := api.Group("/auth")
	auth.Post("/login", authHandler.Login)
	auth.Post("/register", authHandler.Register)
	auth.Post("/recover-password", authHandler.RecoverPassword)
	auth.Get("/reset-password/:token", authHandler.ResetPassword)
	auth.Get("/verify-email/:token", authHandler.VerifyEmail)
	auth.Post("/refresh-token", authHandler.RefreshAccessToken)

	auth.Use(middlewares.AuthenticateJWT(jwtManager))  

	auth.Post("/logout", authHandler.Logout)
}
