package hw09_struct_validator

import (
	"errors"
)

var ErrMaxValidation = errors.New("\"max\" validation error")

type IntValidator struct {
}

func (validator IntValidator) Max(max int, value int) bool {
	return value <= max
}

func (validator IntValidator) Min(min int, value int) bool {
	return value >= min
}

func (validator IntValidator) In(values []int, value int) bool {
	panic("implement me")
}

func Validate() {

}
