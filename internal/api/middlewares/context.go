package middlewares

import (
	"github.com/google/uuid"
	"github.com/gofiber/fiber/v2"
)

// GetClaimFromContext retrieves user information from context
func GetClaimFromContext(c *fiber.Ctx) (*JWTClaims, bool) {
	claims := c.Locals("claims")
	if claims == nil {
		return nil, false
	}

	if jwtClaims, ok := claims.(*JWTClaims); ok {
		return jwtClaims, true
	}

	return nil, false
}

// GetUserIDFromContext gets user ID from context
func GetUserIDFromContext(c *fiber.Ctx) (uuid.UUID, error) {
	userIDStr := c.Locals("user_id")
	if userIDStr == nil {
		return uuid.Nil, fiber.NewError(fiber.StatusUnauthorized, "User ID not found in context")
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		return uuid.Nil, fiber.NewError(fiber.StatusBadRequest, "Invalid user ID format")
	}

	return userID, nil
}

// GetRoleFromContext retrieves user role from context
func GetRoleFromContext(c *fiber.Ctx) (string, error) {
	role := c.Locals("role")
	if role == nil {
		return "", fiber.NewError(fiber.StatusUnauthorized, "User role not found in context")
	}

	return role.(string), nil
}