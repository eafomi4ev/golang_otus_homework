package hw09_struct_validator //nolint:golint,stylecheck

type IntValidator struct{}

func (validator IntValidator) Max(max int, value int) bool {
	return value <= max
}

func (validator IntValidator) Min(min int, value int) bool {
	return value >= min
}

func (validator IntValidator) In(values []int, value int) bool {
	for _, v := range values {
		if v == value {
			return true
		}
	}

	return false
}
