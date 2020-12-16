package hw09_struct_validator

import (
	"errors"
	"fmt"
)

var ErrMaxValidation = errors.New("\"max\" validation error")

func ValidateMax(max int64, value int64) error {
	if value > max {
		return fmt.Errorf("%w: value=%d is bigger than max=%d", ErrMaxValidation, value, max)
	}

	return nil
}
