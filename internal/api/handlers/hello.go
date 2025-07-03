package handlers

import (
	"github.com/gofiber/fiber/v2"
)

// Hello handler for health check
func Hello(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "Hello, World!",
		"status":  "OK",
		"service": "Go Backend Todo API",
	})
}
