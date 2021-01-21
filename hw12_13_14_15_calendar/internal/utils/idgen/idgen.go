package idgen

import (
	"fmt"
	"strings"

	"github.com/gofrs/uuid"
)

func PrefixedID(prefix string) (string, error) {
	u, err := uuid.NewV4()
	if err != nil {
		return "", fmt.Errorf("error occurred while generating event id: %w", err)
	}

	result := prefix + "-" + strings.ToUpper(u.String()[24:])

	return result, nil
}
