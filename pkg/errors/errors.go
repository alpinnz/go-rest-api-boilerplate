package errors

import (
	"errors"
	"fmt"
)

// AppError represents an application-level error with context
type AppError struct {
	Err     error                  // Original error
	Code    string                 // Error code for API responses
	Message string                 // User-friendly message
	Context map[string]interface{} // Additional context
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap implements the errors.Unwrap interface
func (e *AppError) Unwrap() error {
	return e.Err
}

// New creates a new AppError
func New(code, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Context: make(map[string]interface{}),
	}
}

// Wrap wraps an existing error with context
func Wrap(err error, code, message string) *AppError {
	return &AppError{
		Err:     err,
		Code:    code,
		Message: message,
		Context: make(map[string]interface{}),
	}
}

// WithContext adds context to the error
func (e *AppError) WithContext(key string, value interface{}) *AppError {
	e.Context[key] = value
	return e
}

// WithContextMap adds multiple context values
func (e *AppError) WithContextMap(context map[string]interface{}) *AppError {
	for k, v := range context {
		e.Context[k] = v
	}
	return e
}

// Is checks if the error matches a specific error
func (e *AppError) Is(target error) bool {
	if e.Err == nil {
		return false
	}
	return errors.Is(e.Err, target)
}

// Common error constructors

// NotFound creates a not found error
func NotFound(message string) *AppError {
	return New("NOT_FOUND", message)
}

// BadRequest creates a bad request error
func BadRequest(message string) *AppError {
	return New("BAD_REQUEST", message)
}

// Unauthorized creates an unauthorized error
func Unauthorized(message string) *AppError {
	return New("UNAUTHORIZED", message)
}

// Forbidden creates a forbidden error
func Forbidden(message string) *AppError {
	return New("FORBIDDEN", message)
}

// Conflict creates a conflict error
func Conflict(message string) *AppError {
	return New("CONFLICT", message)
}

// InternalServer creates an internal server error
func InternalServer(message string) *AppError {
	return New("INTERNAL_SERVER_ERROR", message)
}

// ValidationFailed creates a validation error
func ValidationFailed(message string) *AppError {
	return New("VALIDATION_FAILED", message)
}

// Timeout creates a timeout error
func Timeout(message string) *AppError {
	return New("TIMEOUT", message)
}

// As attempts to convert error to AppError
func As(err error) (*AppError, bool) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr, true
	}
	return nil, false
}

// GetCode safely gets error code from any error
func GetCode(err error) string {
	if appErr, ok := As(err); ok {
		return appErr.Code
	}
	return "UNKNOWN_ERROR"
}

// GetMessage safely gets error message from any error
func GetMessage(err error) string {
	if appErr, ok := As(err); ok {
		return appErr.Message
	}
	return err.Error()
}

// GetContext safely gets error context
func GetContext(err error) map[string]interface{} {
	if appErr, ok := As(err); ok {
		return appErr.Context
	}
	return nil
}
