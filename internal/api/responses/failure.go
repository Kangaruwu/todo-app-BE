package responses

import (
	"github.com/gofiber/fiber/v2"
)

// ErrorResponse represents an error API response
type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// BadRequest returns a 400 Bad Request error
func BadRequest(c *fiber.Ctx, message string) error {
	if message == "" {
		message = "Bad Request"
	}
	return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
		Success: false,
		Message: message,
	})
}

// BadRequestWithError returns a 400 Bad Request error with error details
func BadRequestWithError(c *fiber.Ctx, message string, err error) error {
	if message == "" {
		message = "Bad Request"
	}
	return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
		Success: false,
		Message: message,
		Error:   err.Error(),
	})
}

// Unauthorized returns a 401 Unauthorized error
func Unauthorized(c *fiber.Ctx, message string) error {
	if message == "" {
		message = "Unauthorized"
	}
	return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
		Success: false,
		Message: message,
	})
}

// Forbidden returns a 403 Forbidden error
func Forbidden(c *fiber.Ctx, message string) error {
	if message == "" {
		message = "Forbidden"
	}
	return c.Status(fiber.StatusForbidden).JSON(ErrorResponse{
		Success: false,
		Message: message,
	})
}

// NotFound returns a 404 Not Found error
func NotFound(c *fiber.Ctx, message string) error {
	if message == "" {
		message = "Not Found"
	}
	return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
		Success: false,
		Message: message,
	})
}

// Conflict returns a 409 Conflict error
func Conflict(c *fiber.Ctx, message string) error {
	if message == "" {
		message = "Conflict"
	}
	return c.Status(fiber.StatusConflict).JSON(ErrorResponse{
		Success: false,
		Message: message,
	})
}

// InternalServerError returns a 500 Internal Server Error
func InternalServerError(c *fiber.Ctx, message string) error {
	if message == "" {
		message = "Internal Server Error"
	}
	return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
		Success: false,
		Message: message,
	})
}

// InternalServerErrorWithError returns a 500 error with error details
func InternalServerErrorWithError(c *fiber.Ctx, message string, err error) error {
	if message == "" {
		message = "Internal Server Error"
	}
	return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
		Success: false,
		Message: message,
		Error:   err.Error(),
	})
}
