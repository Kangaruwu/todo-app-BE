package middlewares

import (
	"strings"
	
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"go-backend-todo/internal/models"
	"go-backend-todo/internal/config"
	"go-backend-todo/internal/repository/auth"
)

// AuthenticateJWT middleware authenticates JWT token
func AuthenticateJWT(c *fiber.Ctx) error {
	// Get token from Authorization header
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Missing authorization header",
		})
	}

	// Check format "Bearer <token>"
	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid authorization header format",
		})
	}

	tokenString := tokenParts[1]

	// Parse and validate token
	claims, err := ParseJWT(tokenString)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid or expired token",
		})
	}

	// Save claims to context for use in handler
	c.Locals("user_id", claims.UserID)
	c.Locals("username", claims.Username)
	c.Locals("email", claims.Email)
	c.Locals("role", claims.Role)
	c.Locals("claims", claims)

	return c.Next()
}

// ParseJWT parse and validate JWT token
func ParseJWT(tokenString string) (*models.JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return config.GetEnv("JWT_ACCESS_SECRET",""), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*models.JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}


// RequireRole middleware requires a specific role
func RequireRole(requiredRole string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("role")
		if role == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized - no role found",
			})
		}

		if role.(string) != requiredRole {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Forbidden - insufficient permissions",
			})
		}

		return c.Next()
	}
}

// RequireAdmin middleware requires admin role
func RequireAdmin() fiber.Handler {
	return RequireRole("admin")
}

// OptionalAuth middleware - token is optional
func OptionalAuth(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Next()
	}

	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return c.Next()
	}

	tokenString := tokenParts[1]
	claims, err := ParseJWT(tokenString)
	if err != nil {
		return c.Next()
	}

	// Save claims to context if token is valid
	c.Locals("user_id", claims.UserID)
	c.Locals("username", claims.Username)
	c.Locals("email", claims.Email)
	c.Locals("role", claims.Role)
	c.Locals("claims", claims)

	return c.Next()
}

// GetUserFromContext retrieves user information from context
func GetUserFromContext(c *fiber.Ctx) (*models.JWTClaims, bool) {
	claims := c.Locals("claims")
	if claims == nil {
		return nil, false
	}

	if jwtClaims, ok := claims.(*models.JWTClaims); ok {
		return jwtClaims, true
	}

	return nil, false
}

// RefreshToken creates a new token from the old token
func RefreshToken(oldTokenString string) (string, error) {
	claims, err := ParseJWT(oldTokenString)
	if err != nil {
		return "", err
	}

	// Create a new token with a new expiration time
	userID, _ := uuid.Parse(claims.UserID)
	return auth_repository.GenerateAccessToken(userID, claims.Username, claims.Email, claims.Role)
}
