package hw02_unpack_string //nolint:golint,stylecheck

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func parseDigitRune(r rune) (int, error) {
	stringRepresentation := string(r)
	return strconv.Atoi(stringRepresentation)
}

func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

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
	count := -1

	for i := len(runes) - 1; i >= 0; i-- {
		if unicode.IsDigit(runes[i]) {
			if count != -1 {
				return "", ErrInvalidString
			}

			count, _ = parseDigitRune(runes[i])
			continue
		}

		if count == 0 {
			count = -1
			continue
		}

		char = string(runes[i])

		if count == -1 {
			builder.WriteString(char)
		} else {
			tmp := strings.Repeat(char, count)
			builder.WriteString(tmp)
			count = -1
		}
	}

	result := reverse(builder.String())

	return result, nil
}
