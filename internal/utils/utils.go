package utils

import (
	"fmt"
)

func ErrNotImplemented(feature string) error {
	return fmt.Errorf("feature not implemented: %s", feature)
}