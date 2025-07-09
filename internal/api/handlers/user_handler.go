package handlers

import (
	"github.com/gofiber/fiber/v2"
	"go-backend-todo/internal/service"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetUser gets user by ID
func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	id := c.Params("id")

	return c.JSON(fiber.Map{
		"message": "Get user",
		"id":      id,
	})
}

// UpdateUser updates user
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	// TODO: Implement logic to update user
	return c.JSON(fiber.Map{
		"message": "Update user",
		"id":      id,
	})
}

// DeleteUser deletes user
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	// TODO: Implement logic to delete user
	return c.JSON(fiber.Map{
		"message": "Delete user",
		"id":      id,
	})
}
