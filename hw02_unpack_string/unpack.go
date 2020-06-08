package hw02_unpack_string //nolint:golint,stylecheck

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	runes := []rune(str)

	if len(runes) == 0 {
		return "", nil
	}

	if unicode.IsDigit(runes[0]) {
		return "", ErrInvalidString
	}

	var builder strings.Builder
	var char string
	for i := 0; i < len(runes)-1; i++ {
		if unicode.IsDigit(runes[i]) {
			if unicode.IsDigit(runes[i+1]) {
				return "", ErrInvalidString
			}
			continue
		}

		char = string(runes[i])

		if unicode.IsDigit(runes[i+1]) {

			countStr := string(runes[i+1])
			count, _ := strconv.Atoi(countStr)

			if count != 0 {
				builder.WriteString(strings.Repeat(char, count))
			}
		} else {
			builder.WriteString(char)
		}
	}

	indexOfLastRune := len(runes) - 1
	if unicode.IsDigit(runes[indexOfLastRune]) == false {
		char = string(runes[indexOfLastRune])
		builder.WriteString(char)
	}

	return builder.String(), nil
}
