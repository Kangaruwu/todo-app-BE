package middlewares

import (
	"log"
	"strings"
	"time"

	"go-backend-todo/internal/api/responses"
	"go-backend-todo/internal/models"

	"github.com/gofiber/fiber/v2"
)

// AuthenticateJWT middleware authenticates JWT token
func AuthenticateJWT(jwtManager *JWTManager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get token from Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return responses.Unauthorized(c, "Authorization header is required")
		}

		// Check format "Bearer <token>"
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			return responses.Unauthorized(c, "Invalid authorization format, expected 'Bearer <token>'")
		}

		tokenString := tokenParts[1]

		// Parse and validate token using JWTManager
		claims, err := jwtManager.ParseAccessToken(tokenString)
		if err != nil {
			log.Printf("JWT Parse Error: %v", err)
			return responses.Unauthorized(c, "Invalid or expired token: "+err.Error())
		}

		// Check if token is expired
		if claims.ExpiresAt.Before(time.Now()) {
			return responses.Unauthorized(c, "Token has expired")
		}

		// Check if user account is active
		if claims.EmailValidationStatus != string(models.EmailValidationStatusEnum("confirmed")) {
			return responses.Unauthorized(c, "Email address is not verified or using old token, please verify your email ")
		}

		// Save claims to context for use in handler
		c.Locals("user_id", claims.UserID)
		c.Locals("username", claims.Username)
		c.Locals("email", claims.Email)
		c.Locals("role", claims.Role)
		c.Locals("email_validation_status", claims.EmailValidationStatus)
		c.Locals("claims", claims)

		return c.Next()
	}
}

// RequireRole middleware requires a specific role
func RequireRole(requiredRole string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role, err := GetRoleFromContext(c)
		if err != nil {
			return responses.Unauthorized(c, "Unauthorized - no role found")
		}

		if role != requiredRole {
			return responses.Forbidden(c, "Forbidden - insufficient permissions")
		}

		return c.Next()
	}
}

// RequireAdmin middleware requires admin role
func RequireAdmin() fiber.Handler {
	return RequireRole("admin")
}

// OptionalAuth middleware - token is optional
func OptionalAuth(jwtManager *JWTManager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Next()
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			return c.Next()
		}

		tokenString := tokenParts[1]
		claims, err := jwtManager.ParseAccessToken(tokenString)
		if err != nil {
			return c.Next()
		}

		// Save claims to context if token is valid
		c.Locals("user_id", claims.UserID)
		c.Locals("username", claims.Username)
		c.Locals("email", claims.Email)
		c.Locals("role", claims.Role)
		c.Locals("email_validation_status", claims.EmailValidationStatus)
		c.Locals("claims", claims)

		return c.Next()
	}
}



