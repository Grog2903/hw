package hw09structvalidator

import (
	"encoding/json"
	"errors"
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
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
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
			in: User{
				ID:     "123456789012345678901234567890123456",
				Age:    20,
				Email:  "test@test.ru",
				Role:   UserRole("admin"),
				Phones: []string{"12345678912", "12345678912"},
			},
			expectedErr: nil,
		},
		{
			in: User{
				ID:     "123456789012345678901234567890123456",
				Age:    20,
				Email:  "test@test.ru",
				Role:   UserRole("admin"),
				Phones: []string{"12345678912", "1234567891212323"},
			},
			expectedErr: ErrIncorrectLength,
		},
		{
			in: App{
				Version: "0.0.12",
			},
			expectedErr: ErrIncorrectLength,
		},
		{
			in:          Token{},
			expectedErr: nil,
		},
		{
			in: Response{
				Code: 200,
			},
			expectedErr: nil,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			validateError, _ := Validate(tt.in)
			if tt.expectedErr != nil {
				require.Error(t, validateError)

				require.True(t, errors.Is(validateError, tt.expectedErr))
			} else {
				require.True(t, validateError.Error() == "")
			}
		})
	}
}
