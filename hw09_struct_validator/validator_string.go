package hw09_struct_validator

import (
	"fmt"
	"regexp"
)

type StringValidator struct{}

func (validator StringValidator) Regexp(pattern string, value string) (bool, error) {
	ok, err := regexp.MatchString(pattern, value)
	if err != nil {
		return false, fmt.Errorf("%w, because of %v", ErrValidationRegexp, err)
	}

	return ok, nil
}

func (validator StringValidator) Len(expectedLen int, value string) bool {
	return len(value) == expectedLen
}

func (validator StringValidator) In(values []string, value string) bool {
	for _, v := range values {
		if v == value {
			return true
		}
	}

	return false
}
