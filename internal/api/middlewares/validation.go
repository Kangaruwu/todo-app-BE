package middlewares

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"encoding/json"

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
		return fmt.Errorf("%s", strings.Join(errors, ", "))
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

// RequestValidation middleware validates request body and query parameters
func RequestValidation(body *[]byte, req interface{}) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if err := json.Unmarshal(*body, req); err != nil {
			return fmt.Errorf("invalid request body: %s", err.Error())
		}

		if err := ValidateStruct(req); err != nil {
			return fmt.Errorf("validation error: %s", err.Error())
		}

		return c.Next()
	}
}
