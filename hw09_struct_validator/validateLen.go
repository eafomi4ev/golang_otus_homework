package hw09_struct_validator

import (
	"errors"
	"fmt"
)

var ErrLengthValidation = errors.New("\"len\" validation error")

func ValidateLen(length int64, value interface{}) []error {
	validationErrors := make([]error, 0)

	switch value.(type) {
	case string:
		err := doValidate(length, value.(string))
		if err != nil {
			validationErrors = append(validationErrors, err)
		}
	case []string:
		for _, v := range value.([]string) {
			err := doValidate(length, v)
			if err != nil {
				validationErrors = append(validationErrors, err)
			}
		}
	}

	if len(validationErrors) != 0 {
		return validationErrors
	}

	return nil
}

func doValidate(length int64, value string) error {
	if int64(len(value)) != length {
		return fmt.Errorf("%w: the length of field's value=%d is not equal to len=%d", ErrLengthValidation, len(value), length)
	}

	return nil
}
