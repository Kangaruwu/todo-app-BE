package middlewares

import (
	"strings"
	"time"
	"log"
	"fmt"

	"go-backend-todo/internal/config"
	"go-backend-todo/internal/api/responses"
	"go-backend-todo/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// AuthenticateJWT middleware authenticates JWT token
func AuthenticateJWT(c *fiber.Ctx) error {
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

	// Parse and validate token
	claims, err := ParseJWT(tokenString)
	if err != nil {
		log.Printf("JWT Parse Error: %v", err)
		return responses.Unauthorized(c, "Invalid or expired token")
	}

	// Check if token is expired
	if claims.ExpiresAt.Before(time.Now()) {
		return responses.Unauthorized(c, "Token has expired")
	}

	// Check if user account is active
	if claims.EmailValidationStatus != string(models.EmailValidationStatusEnum("confirmed")) {
		return responses.Unauthorized(c, "Email address is not verified")
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

// ParseJWT parse and validate JWT token
func ParseJWT(tokenString string) (*JWTClaims, error) {
    cfg := config.Load() 
    
    token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(cfg.JWT.AccessSecret), nil
    })

    if err != nil {
        log.Printf("JWT Parse Error: %v", err)
        return nil, err
    }

    if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
        return claims, nil
    }

    return nil, fmt.Errorf("invalid token claims")
}

// RequireRole middleware requires a specific role
func RequireRole(requiredRole string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("role")
		if role == nil {
			return responses.Unauthorized(c, "Unauthorized - no role found")
		}

		if role.(string) != requiredRole {
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
	c.Locals("email_validation_status", claims.EmailValidationStatus)
	c.Locals("claims", claims)

	return c.Next()
}

// GetUserFromContext retrieves user information from context
func GetUserFromContext(c *fiber.Ctx) (*JWTClaims, bool) {
	claims := c.Locals("claims")
	if claims == nil {
		return nil, false
	}

	if jwtClaims, ok := claims.(*JWTClaims); ok {
		return jwtClaims, true
	}

	return nil, false
}

// // RefreshToken creates a new token from the old token
// func RefreshToken(oldTokenString string) (string, error) {
// 	claims, err := ParseJWT(oldTokenString)
// 	if err != nil {
// 		return "", err
// 	}

// 	// Create a new token with a new expiration time
// 	userID, _ := uuid.Parse(claims.UserID)
// 	return GenerateAccessToken(userID, claims.Username, claims.Email, claims.Role)
// }

// CheckAccountStatus checks if the user account is active
