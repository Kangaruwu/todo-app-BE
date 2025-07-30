package utils

import (
	"fmt"
	"math/rand"
)

func ErrNotImplemented(feature string) error {
	return fmt.Errorf("feature not implemented: %s", feature)
}

func ErrInvalidCredentials(message string) error {
	return fmt.Errorf("invalid credentials: %s", message)
}

func ErrInternalServerError(message string) error {
	return fmt.Errorf("internal server error: %s", message)
}

func ErrEmailAlreadyExists(message string) error {
	return fmt.Errorf("email already exists: %s", message)
}

func ErrUsernameAlreadyExists(message string) error {
	return fmt.Errorf("username already exists: %s", message)
}

func ErrAccountNotActive(message string) error {
	return fmt.Errorf("account not active: %s", message)
}

func ErrUserNotFound(message string) error {
	return fmt.Errorf("user not found: %s", message)
}

func ErrUnauthorized(message string) error {
	return fmt.Errorf("unauthorized: %s", message)
}

func ErrInvalidInput(message string) error {
	return fmt.Errorf("invalid input: %s", message)
}

func ErrTimeout(message string) error {
	return fmt.Errorf("timeout: %s", message)
}

func RandInRange(min, max int) int {
	if min > max {
		min, max = max, min
	}
	return rand.Intn(max-min+1) + min
}
