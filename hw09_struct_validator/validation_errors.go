package hw09_struct_validator

import (
	"fmt"

	"github.com/pkg/errors"
)

var ErrValidation = errors.New("validation error")
var ErrValidationMax = fmt.Errorf("%w. Rule: \"max\"", ErrValidation)
var ErrValidationMin = fmt.Errorf("%w. Rule: \"min\"", ErrValidation)
var ErrValidationIn = fmt.Errorf("%w. Rule: \"in\"", ErrValidation)
var ErrValidationLen = fmt.Errorf("%w. Rule: \"len\"", ErrValidation)
var ErrValidationRegexp = fmt.Errorf("%w. Rule: \"regexp\"", ErrValidation)
