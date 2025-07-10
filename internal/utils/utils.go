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

func RandInRange(min, max int) int {
	if min > max {
		min, max = max, min
	}
	return rand.Intn(max-min+1) + min
}
