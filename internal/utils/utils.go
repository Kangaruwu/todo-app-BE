package utils

import (
	"fmt"
	"math/rand"
)

func ErrNotImplemented(feature string) error {
	return fmt.Errorf("feature not implemented: %s", feature)
}

func RandInRange(min, max int) int {
	if min > max {
		min, max = max, min
	}
	return rand.Intn(max-min+1) + min
}
