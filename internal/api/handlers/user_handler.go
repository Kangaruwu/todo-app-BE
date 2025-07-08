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

// GetUsers gets all users
func (h *UserHandler) GetUsers(c *fiber.Ctx) error {
	// TODO: Implement logic to get users from database

	user, err := h.userService.GetUserByID(c)
	if err != nil {
		return c.JSON(fiber.Map{
			"message": "Get users",
			"data":    "Ko co j dau loi roi",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Get users",
		"data":    user,
	})
}

// CreateUser creates a new user
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	// TODO: Implement logic to create user
	return c.JSON(fiber.Map{
		"message": "Create user",
	})
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
