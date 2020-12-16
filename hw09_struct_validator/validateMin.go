package hw09_struct_validator

import (
	"fmt"
)

func ValidateMin(min int64, value int64) error {
	if value < min {
		return fmt.Errorf("\"min\" validation error: value=%d is less than min=%d", value, min)
	}

	return nil
}
