package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required,min=3,max=20"`
	Age      int    `json:"age" validate:"required,gte=18,lte=100"`
}

func TestValidate_Success(t *testing.T) {
	data := TestStruct{
		Email:    "test@example.com",
		Username: "testuser",
		Age:      25,
	}

	errors := Validate(data)
	assert.Empty(t, errors)
}

func TestValidate_RequiredField(t *testing.T) {
	data := TestStruct{
		Email:    "",
		Username: "testuser",
		Age:      25,
	}

	errors := Validate(data)
	assert.NotEmpty(t, errors)
	assert.Equal(t, "email", errors[0].Field)
}

func TestValidate_EmailFormat(t *testing.T) {
	data := TestStruct{
		Email:    "invalid-email",
		Username: "testuser",
		Age:      25,
	}

	errors := Validate(data)
	assert.NotEmpty(t, errors)
	assert.Equal(t, "email", errors[0].Field)
}

func TestValidate_MinLength(t *testing.T) {
	data := TestStruct{
		Email:    "test@example.com",
		Username: "ab",
		Age:      25,
	}

	errors := Validate(data)
	assert.NotEmpty(t, errors)
	assert.Equal(t, "username", errors[0].Field)
}

func TestValidate_NumberRange(t *testing.T) {
	data := TestStruct{
		Email:    "test@example.com",
		Username: "testuser",
		Age:      15,
	}

	errors := Validate(data)
	assert.NotEmpty(t, errors)
	assert.Equal(t, "age", errors[0].Field)
}
