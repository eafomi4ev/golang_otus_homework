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

	Buyer struct {
		Age     int    `validate:"min:18|max:50"`
		Name    string `validate:"len:5"`
		Address string
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
				//ID:     "12345678-12345678-12345678-123456789",
				//Name:   "John",
				//Age:    34,
				//Email:  "test@test.m",
				//Role:   "admin",
				Phones: []string{"899944422331"},
				//meta:   nil,
			},
			nil,
		},
		{
			Buyer{Name: "Bobby", Age: 20, Address: "NY"},
			nil,
		},
		{
			Buyer{Name: "Bobby", Age: 60, Address: "NY"},
			ValidationErrors{
				ValidationError{
					Field: "Age",
					Err:   fmt.Errorf("%w: value=%d is bigger than max=%d", ErrMaxValidation, 60, 50),
				},
			},
		},
		{
			Buyer{Name: "Ann", Age: 60, Address: "Mexico"},
			ValidationErrors{
				ValidationError{
					Field: "Age",
					Err:   fmt.Errorf("%w: value=%d is bigger than max=%d", ErrMaxValidation, 60, 50),
				},
				ValidationError{
					Field: "Name",
					Err:   fmt.Errorf("%w: the length of field's value=%d is not equal to len=%d", ErrLengthValidation, 3, 5),
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
