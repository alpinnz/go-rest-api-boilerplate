package errors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	err := New("TEST_CODE", "Test message")

	assert.NotNil(t, err)
	assert.Equal(t, "TEST_CODE", err.Code)
	assert.Equal(t, "Test message", err.Message)
	assert.Empty(t, err.Context)
}

func TestWrap(t *testing.T) {
	originalErr := errors.New("original error")
	err := Wrap(originalErr, "WRAPPED", "Wrapped message")

	assert.NotNil(t, err)
	assert.Equal(t, "WRAPPED", err.Code)
	assert.Equal(t, "Wrapped message", err.Message)
	assert.Equal(t, originalErr, err.Err)
}

func TestWithContext(t *testing.T) {
	err := New("TEST_CODE", "Test message")
	err.WithContext("user_id", "123")
	err.WithContext("action", "create")

	assert.Len(t, err.Context, 2)
	assert.Equal(t, "123", err.Context["user_id"])
	assert.Equal(t, "create", err.Context["action"])
}

func TestWithContextMap(t *testing.T) {
	err := New("TEST_CODE", "Test message")
	context := map[string]interface{}{
		"user_id": "123",
		"action":  "update",
	}
	err.WithContextMap(context)

	assert.Len(t, err.Context, 2)
	assert.Equal(t, "123", err.Context["user_id"])
	assert.Equal(t, "update", err.Context["action"])
}

func TestError(t *testing.T) {
	err := New("TEST_CODE", "Test message")
	assert.Equal(t, "Test message", err.Error())

	originalErr := errors.New("original error")
	wrappedErr := Wrap(originalErr, "WRAPPED", "Wrapped message")
	assert.Contains(t, wrappedErr.Error(), "Wrapped message")
	assert.Contains(t, wrappedErr.Error(), "original error")
}

func TestUnwrap(t *testing.T) {
	originalErr := errors.New("original error")
	wrappedErr := Wrap(originalErr, "WRAPPED", "Wrapped message")

	assert.Equal(t, originalErr, wrappedErr.Unwrap())
}

func TestIs(t *testing.T) {
	originalErr := errors.New("original error")
	wrappedErr := Wrap(originalErr, "WRAPPED", "Wrapped message")

	assert.True(t, wrappedErr.Is(originalErr))
	assert.False(t, wrappedErr.Is(errors.New("different error")))
}

func TestCommonErrorConstructors(t *testing.T) {
	tests := []struct {
		name    string
		err     *AppError
		code    string
		message string
	}{
		{"NotFound", NotFound("Resource not found"), "NOT_FOUND", "Resource not found"},
		{"BadRequest", BadRequest("Invalid input"), "BAD_REQUEST", "Invalid input"},
		{"Unauthorized", Unauthorized("Auth required"), "UNAUTHORIZED", "Auth required"},
		{"Forbidden", Forbidden("Access denied"), "FORBIDDEN", "Access denied"},
		{"Conflict", Conflict("Duplicate entry"), "CONFLICT", "Duplicate entry"},
		{"InternalServer", InternalServer("Server error"), "INTERNAL_SERVER_ERROR", "Server error"},
		{"ValidationFailed", ValidationFailed("Validation error"), "VALIDATION_FAILED", "Validation error"},
		{"Timeout", Timeout("Request timeout"), "TIMEOUT", "Request timeout"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.code, tt.err.Code)
			assert.Equal(t, tt.message, tt.err.Message)
		})
	}
}

func TestAs(t *testing.T) {
	appErr := New("TEST_CODE", "Test message")

	convertedErr, ok := As(appErr)
	assert.True(t, ok)
	assert.Equal(t, appErr, convertedErr)

	regularErr := errors.New("regular error")
	_, ok = As(regularErr)
	assert.False(t, ok)
}

func TestGetCode(t *testing.T) {
	appErr := New("TEST_CODE", "Test message")
	assert.Equal(t, "TEST_CODE", GetCode(appErr))

	regularErr := errors.New("regular error")
	assert.Equal(t, "UNKNOWN_ERROR", GetCode(regularErr))
}

func TestGetMessage(t *testing.T) {
	appErr := New("TEST_CODE", "Test message")
	assert.Equal(t, "Test message", GetMessage(appErr))

	regularErr := errors.New("regular error")
	assert.Equal(t, "regular error", GetMessage(regularErr))
}

func TestGetContext(t *testing.T) {
	appErr := New("TEST_CODE", "Test message")
	appErr.WithContext("key", "value")

	context := GetContext(appErr)
	assert.NotNil(t, context)
	assert.Equal(t, "value", context["key"])

	regularErr := errors.New("regular error")
	context = GetContext(regularErr)
	assert.Nil(t, context)
}
