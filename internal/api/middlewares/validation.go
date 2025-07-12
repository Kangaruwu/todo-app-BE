package middlewares

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// Validator instance
var validate = validator.New()

// ValidateStruct validates a struct and returns formatted errors
func ValidateStruct(s interface{}) error {
	if err := validate.Struct(s); err != nil {
		var errors []string
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, fmt.Sprintf("%s: %s", err.Field(), getValidationMessage(err)))
		}
		return fmt.Errorf(strings.Join(errors, ", "))
	}
	return nil
}

// getValidationMessage returns user-friendly validation messages
func getValidationMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "is required"
	case "email":
		return "must be a valid email address"
	case "min":
		return fmt.Sprintf("must be at least %s characters long", err.Param())
	case "max":
		return fmt.Sprintf("must be at most %s characters long", err.Param())
	case "eqfield":
		return fmt.Sprintf("must match %s", err.Param())
	case "oneof":
		return fmt.Sprintf("must be one of: %s", err.Param())
	case "hexcolor":
		return "must be a valid hex color code"
	case "uuid":
		return "must be a valid UUID"
	default:
		return "is invalid"
	}
}

// ValidateBody middleware validates request body
func ValidateBody(model interface{}) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if err := c.BodyParser(model); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Invalid request body format",
				"success": false,
			})
		}

		if err := ValidateStruct(model); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   err.Error(),
				"success": false,
			})
		}

		// Store parsed and validated data in context
		c.Locals("validated_body", model)
		return c.Next()
	}
}

// ValidateQuery validates query parameters
func ValidateQuery(model interface{}) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if err := c.QueryParser(model); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Invalid query parameters",
				"success": false,
			})
		}

		if err := ValidateStruct(model); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   err.Error(),
				"success": false,
			})
		}

		// Store parsed and validated data in context
		c.Locals("validated_query", model)
		return c.Next()
	}
}

// GetValidatedBody retrieves validated body from context
func GetValidatedBody(c *fiber.Ctx) interface{} {
	return c.Locals("validated_body")
}

// GetValidatedQuery retrieves validated query from context
func GetValidatedQuery(c *fiber.Ctx) interface{} {
	return c.Locals("validated_query")
}
