// backend/internal/util/validator_test.go
package util_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/your-org/esms/internal/util"
)

type TestStruct struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
	Age   int    `json:"age" validate:"min=18"`
}

func TestValidator_ValidateStruct(t *testing.T) {
	v := util.NewValidator()

	tests := []struct {
		name    string
		input   TestStruct
		wantErr bool
	}{
		{
			name: "Valid input",
			input: TestStruct{
				Name:  "John Doe",
				Email: "john@example.com",
				Age:   20,
			},
			wantErr: false,
		},
		{
			name: "Missing required field",
			input: TestStruct{
				Email: "john@example.com",
				Age:   20,
			},
			wantErr: true,
		},
		{
			name: "Invalid email",
			input: TestStruct{
				Name:  "John Doe",
				Email: "invalid-email",
				Age:   20,
			},
			wantErr: true,
		},
		{
			name: "Invalid min value",
			input: TestStruct{
				Name:  "John Doe",
				Email: "john@example.com",
				Age:   10,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.ValidateStruct(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestFormatError(t *testing.T) {
	v := util.NewValidator()
	input := TestStruct{Age: 10} // Missing Name, Email, Invalid Age
	err := v.ValidateStruct(input)

	assert.Error(t, err)

	formatted := util.FormatError(err)
	assert.Contains(t, formatted, "name")
	assert.Equal(t, "This field is required", formatted["name"])

	assert.Contains(t, formatted, "email")
	assert.Equal(t, "This field is required", formatted["email"])

	assert.Contains(t, formatted, "age")
	assert.Equal(t, "Must be at least 18 characters long", formatted["age"]) // validatorのminメッセージは文字列長扱いになることがあるが、ここではメッセージ内容の確認のみ
}
