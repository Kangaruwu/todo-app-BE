package handlers

import (
	"github.com/gofiber/fiber/v2"
)

// User handlers - Legacy handlers (no dependency injection yet)
// TODO: Refactor to use dependency injection like TodoHandler


// GetUsers gets all users
func GetUsers(c *fiber.Ctx) error {
	// TODO: Implement logic to get users from database

	

	return c.JSON(fiber.Map{
		"message": "Get users",
		"data":    []interface{}{},
	})
}

// CreateUser creates a new user
func CreateUser(c *fiber.Ctx) error {
	// TODO: Implement logic to create user
	return c.JSON(fiber.Map{
		"message": "Create user",
	})
}

// GetUser gets user by ID
func GetUser(c *fiber.Ctx) error {
	id := c.Params("id")

	return c.JSON(fiber.Map{
		"message": "Get user",
		"id":      id,
	})
}

// UpdateUser updates user
func UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	// TODO: Implement logic to update user
	return c.JSON(fiber.Map{
		"message": "Update user",
		"id":      id,
	})
}

// DeleteUser deletes user
func DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	// TODO: Implement logic to delete user
	return c.JSON(fiber.Map{
		"message": "Delete user",
		"id":      id,
	})
}
