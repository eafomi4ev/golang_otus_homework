package hw09_struct_validator //nolint:golint,stylecheck

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			User{
				ID:     "12345678-12345678-12345678-123456789",
				Name:   "John",
				Age:    34,
				Email:  "test@test.m",
				Role:   "admin",
				Phones: []string{"89994442233"},
				meta:   nil,
			},
			nil,
		},
		{
			User{
				ID:     "12345678-12345678-12345678-123456789",
				Name:   "Ann",
				Age:    25,
				Email:  "ann-WRONG_EMAIL",
				Role:   "admin",
				Phones: []string{"89994442233"},
				meta:   nil,
			},
			ValidationErrors{
				ValidationError{
					Field: "Email",
					Err:   fmt.Errorf("%w. Rule value: %s. Struct value: %s", ErrValidationRegexp, "^\\w+@\\w+\\.\\w+$", "ann-WRONG_EMAIL"),
				},
			},
		},
		{
			User{
				ID:     "",
				Name:   "John",
				Age:    34,
				Email:  "test@adm.ew",
				Role:   "admin",
				Phones: []string{"89994442233", "8999444223310"},
				meta:   nil,
			},
			ValidationErrors{
				ValidationError{
					Field: "ID",
					Err:   fmt.Errorf("%w. Rule value: %d. Struct value: %s", ErrValidationLen, 36, ""),
				},
				ValidationError{
					Field: "Phones",
					Err:   fmt.Errorf("%w. Rule value: %d. Struct value: %s", ErrValidationLen, 11, "8999444223310"),
				},
			},
		},
		{
			App{Version: "1.0.0"},
			nil,
		},
		{
			App{Version: "1.0.0."},
			ValidationErrors{
				ValidationError{
					Field: "Version",
					Err:   fmt.Errorf("%w. Rule value: %d. Struct value: %s", ErrValidationLen, 5, "1.0.0."),
				},
			},
		},
		{
			Token{
				Header:    []byte{0, 0, 0, 50},
				Payload:   []byte{0, 0, 0, 40},
				Signature: []byte{1, 0, 0, 30},
			},
			nil,
		},
		{
			Response{
				Code: 200,
				Body: "Hello",
			},
			nil,
		},
		{
			Response{
				Code: 418,
				Body: "Ooops",
			},
			ValidationErrors{
				ValidationError{
					Field: "Code",
					Err:   fmt.Errorf("%w. Rule values: %s. Struct value: %s", ErrValidationIn, "200,404,500", "418"),
				},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			err := Validate(tt.in)

			require.Equal(t, tt.expectedErr, err)
		})
	}
}
