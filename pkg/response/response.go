// Package response provides HTTP response formatting utilities for consistent API responses.
package response

import (
	"github.com/alpinnz/go-rest-api-boilerplate/internal/domain"
	"github.com/alpinnz/go-rest-api-boilerplate/internal/localization"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Response represents standardized API response structure.
type Response struct {
	Code      *string     `json:"code,omitempty"`       // Application-level code or error type
	Message   string      `json:"message"`              // Human-readable message
	Data      interface{} `json:"data,omitempty"`       // Response payload
	Errors    interface{} `json:"errors,omitempty"`     // Validation or contextual error details
	Meta      interface{} `json:"meta,omitempty"`       // Metadata (pagination, sorting, etc.)
	RequestID *string     `json:"request_id,omitempty"` // Request ID generated internally per request
}

// ValidationItem represents a single validation error.
type ValidationItem struct {
	Code    string `json:"code"`    // Error code (e.g., "REQUIRED", "INVALID_FORMAT")
	Field   string `json:"field"`   // Field name that failed validation
	Message string `json:"message"` // Human-readable error message
}

// ValidationErrors represents categorized validation errors.
type ValidationErrors struct {
	Body    []ValidationItem `json:"body"`
	Headers []ValidationItem `json:"headers"`
	Path    []ValidationItem `json:"path"`
	Query   []ValidationItem `json:"query"`
}

// FieldError represents a field validation error with localization support.
type FieldError struct {
	Field   string `json:"field"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Sort represents sorting criteria.
type Sort struct {
	Field string `json:"field"`
	Order string `json:"order"`
}

// Meta represents response metadata.
// Uses pagination.Meta for pagination data.
type Meta struct {
	Pagination interface{} `json:"pagination,omitempty"` // Use pagination.Meta type
	Sort       []Sort      `json:"sort,omitempty"`
}

// NewMeta creates a new Meta instance.
// For pagination, pass pagination.Meta from pkg/pagination.
// Example:
//
//	paginationMeta := pagination.NewMeta(page, perPage, totalData)
//	meta := response.NewMeta(paginationMeta, nil)
func NewMeta(paginationMeta interface{}, sort []Sort) *Meta {
	return &Meta{
		Pagination: paginationMeta,
		Sort:       sort,
	}
}

// getRequestID retrieves request ID from context.
// Request ID is generated internally by RequestID middleware.
func getRequestID(c *gin.Context) *string {
	if requestID, exists := c.Get("request_id"); exists {
		if id, ok := requestID.(string); ok {
			return &id
		}
	}
	// Fallback: generate new UUID if not set (shouldn't happen with middleware)
	id := uuid.New().String()
	return &id
}

// getDict retrieves localization dictionary from context.
// Returns nil if not found, allowing functions to use fallback messages.
func getDict(c *gin.Context) localization.Dictionary {
	if dict, exists := c.Get("lang_dict"); exists {
		if d, ok := dict.(localization.Dictionary); ok {
			return d
		}
	}
	return nil
}

// getMessage gets localized message from dictionary or returns fallback.
func getMessage(c *gin.Context, key string, fallback string) string {
	dict := getDict(c)
	if dict == nil {
		return fallback
	}
	if msg, ok := getByPath(dict, key); ok {
		return msg
	}
	return fallback
}

// Success sends a successful JSON response (200 OK).
// Automatically uses localization from context if available.
func Success(c *gin.Context, data interface{}) {
	c.JSON(200, Response{
		Code:      nil,
		Message:   getMessage(c, "common.success", "Success"),
		Data:      data,
		Errors:    nil,
		Meta:      nil,
		RequestID: getRequestID(c),
	})
}

// SuccessWithMeta sends a successful JSON response with metadata (200 OK).
// Automatically uses localization from context if available.
func SuccessWithMeta(c *gin.Context, data interface{}, meta *Meta) {
	c.JSON(200, Response{
		Code:      nil,
		Message:   getMessage(c, "common.success", "Success"),
		Data:      data,
		Errors:    nil,
		Meta:      meta,
		RequestID: getRequestID(c),
	})
}

// Created sends a resource created response (201 Created).
// Automatically uses localization from context if available.
// If messageKey is empty, uses "common.created" as default.
func Created(c *gin.Context, data interface{}, messageKey ...string) {
	key := "common.created"
	if len(messageKey) > 0 && messageKey[0] != "" {
		key = messageKey[0]
	}

	c.JSON(201, Response{
		Code:      nil,
		Message:   getMessage(c, key, "Resource created successfully"),
		Data:      data,
		Errors:    nil,
		Meta:      nil,
		RequestID: getRequestID(c),
	})
}

// NoContent sends a no content response (204 No Content).
// Automatically uses localization from context if available.
func NoContent(c *gin.Context) {
	c.JSON(204, Response{
		Code:      nil,
		Message:   getMessage(c, "common.no_content", "No Content"),
		Data:      nil,
		Errors:    nil,
		Meta:      nil,
		RequestID: getRequestID(c),
	})
}

// Redirect sends a redirect response (301 Moved Permanently).
// Automatically uses localization from context if available.
func Redirect(c *gin.Context, location string) {
	code := CodeRedirect
	c.JSON(301, Response{
		Code:      &code,
		Message:   getMessage(c, "common.redirect", "Resource moved permanently"),
		Data:      gin.H{"location": location},
		Errors:    nil,
		Meta:      nil,
		RequestID: getRequestID(c),
	})
}

// TemporaryRedirect sends a temporary redirect response (302 Found).
// Automatically uses localization from context if available.
func TemporaryRedirect(c *gin.Context, location string) {
	code := CodeRedirect
	c.JSON(302, Response{
		Code:      &code,
		Message:   getMessage(c, "common.temporary_redirect", "Resource temporarily moved"),
		Data:      gin.H{"location": location},
		Errors:    nil,
		Meta:      nil,
		RequestID: getRequestID(c),
	})
}

// ValidationError sends a validation error response (400 Bad Request).
// Accepts either ValidationErrors struct or []domain.AppError.
// Automatically uses localization from context if available.
func ValidationError(c *gin.Context, errors interface{}) {
	dict := getDict(c)
	msg := getMessage(c, "common.validation_error", "Validation failed")
	code := CodeValidationError

	// Handle different error types
	var errorData interface{}
	if appErrs, ok := errors.([]domain.AppError); ok && dict != nil {
		// Translate AppErrors with localization
		errorData = gin.H{"fields": TranslateValidationErrors(dict, appErrs)}
	} else {
		// Use errors as-is (ValidationErrors struct or other)
		errorData = errors
	}

	c.JSON(400, Response{
		Code:      &code,
		Message:   msg,
		Data:      nil,
		Errors:    errorData,
		Meta:      nil,
		RequestID: getRequestID(c),
	})
}

// BadRequest sends a bad request error response (400 Bad Request).
// Can accept (code, message, errors) for manual usage or (AppError) for localized usage.
// Automatically uses localization from context if available.
func BadRequest(c *gin.Context, args ...interface{}) {

	// Check if first arg is AppError
	if len(args) == 1 {
		if appErr, ok := args[0].(domain.AppError); ok {
			// Use localized error
			Error(c, 400, appErr)
			return
		}
	}

	// Legacy signature: code, message, errors
	var code string
	var message string
	var errors interface{}

	if len(args) >= 1 {
		if c, ok := args[0].(string); ok {
			code = c
		}
	}
	if len(args) >= 2 {
		if m, ok := args[1].(string); ok {
			message = m
		}
	}
	if len(args) >= 3 {
		errors = args[2]
	}

	if code == "" {
		code = CodeBadRequest
	}
	if message == "" {
		message = "Bad request"
	}

	c.JSON(400, Response{
		Code:      &code,
		Message:   message,
		Data:      nil,
		Errors:    errors,
		Meta:      nil,
		RequestID: getRequestID(c),
	})
}

// Unauthorized sends an unauthorized error response (401 Unauthorized).
// Can accept (reason string) for manual usage or (AppError) for localized usage.
// Automatically uses localization from context if available.
func Unauthorized(c *gin.Context, args ...interface{}) {
	// Check if first arg is AppError
	if len(args) == 1 {
		if appErr, ok := args[0].(domain.AppError); ok {
			Error(c, 401, appErr)
			return
		}
	}

	// Legacy signature: reason string
	reason := "Unauthorized"
	if len(args) >= 1 {
		if r, ok := args[0].(string); ok {
			reason = r
		}
	}

	code := CodeUnauthorized
	c.JSON(401, Response{
		Code:      &code,
		Message:   getMessage(c, "auth.unauthorized", "Unauthorized"),
		Data:      nil,
		Errors:    gin.H{"reason": reason},
		Meta:      nil,
		RequestID: getRequestID(c),
	})
}

// Forbidden sends a forbidden error response (403 Forbidden).
// Can accept (permission string) for manual usage or (AppError) for localized usage.
// Automatically uses localization from context if available.
func Forbidden(c *gin.Context, args ...interface{}) {
	// Check if first arg is AppError
	if len(args) == 1 {
		if appErr, ok := args[0].(domain.AppError); ok {
			Error(c, 403, appErr)
			return
		}
	}

	// Legacy signature: permission string
	permission := ""
	if len(args) >= 1 {
		if p, ok := args[0].(string); ok {
			permission = p
		}
	}

	code := CodeForbidden
	c.JSON(403, Response{
		Code:      &code,
		Message:   getMessage(c, "common.forbidden", "Access denied"),
		Data:      nil,
		Errors:    gin.H{"permission": permission},
		Meta:      nil,
		RequestID: getRequestID(c),
	})
}

// NotFound sends a not found error response (404 Not Found).
// Can accept (resource string) for manual usage or (AppError) for localized usage.
// Automatically uses localization from context if available.
func NotFound(c *gin.Context, args ...interface{}) {
	// Check if first arg is AppError
	if len(args) == 1 {
		if appErr, ok := args[0].(domain.AppError); ok {
			Error(c, 404, appErr)
			return
		}
	}

	// Legacy signature: resource string
	resource := ""
	if len(args) >= 1 {
		if r, ok := args[0].(string); ok {
			resource = r
		}
	}

	code := CodeNotFound
	c.JSON(404, Response{
		Code:      &code,
		Message:   getMessage(c, "common.not_found", "Resource not found"),
		Data:      nil,
		Errors:    gin.H{"resource": resource},
		Meta:      nil,
		RequestID: getRequestID(c),
	})
}

// Conflict sends a conflict error response (409 Conflict).
// Can accept (field, message string) for manual usage or (AppError) for localized usage.
// Automatically uses localization from context if available.
func Conflict(c *gin.Context, args ...interface{}) {
	// Check if first arg is AppError
	if len(args) == 1 {
		if appErr, ok := args[0].(domain.AppError); ok {
			Error(c, 409, appErr)
			return
		}
	}

	// Legacy signature: field, message
	var field, message string
	if len(args) >= 1 {
		if f, ok := args[0].(string); ok {
			field = f
		}
	}
	if len(args) >= 2 {
		if m, ok := args[1].(string); ok {
			message = m
		}
	}

	code := CodeConflict
	c.JSON(409, Response{
		Code:      &code,
		Message:   getMessage(c, "common.conflict", "Resource conflict"),
		Data:      nil,
		Errors:    gin.H{"field": field, "message": message},
		Meta:      nil,
		RequestID: getRequestID(c),
	})
}

// InternalServerError sends an internal server error response (500 Internal Server Error).
// Can accept (hint string) for manual usage or (AppError) for localized usage.
// Automatically uses localization from context if available.
func InternalServerError(c *gin.Context, args ...interface{}) {
	// Check if first arg is AppError
	if len(args) == 1 {
		if appErr, ok := args[0].(domain.AppError); ok {
			Error(c, 500, appErr)
			return
		}
	}

	// Legacy signature: hint string
	hint := ""
	if len(args) >= 1 {
		if h, ok := args[0].(string); ok {
			hint = h
		}
	}

	code := CodeInternalServerError
	c.JSON(500, Response{
		Code:      &code,
		Message:   getMessage(c, "common.error", "Unexpected error occurred"),
		Data:      nil,
		Errors:    gin.H{"hint": hint},
		Meta:      nil,
		RequestID: getRequestID(c),
	})
}

// Error sends error response with localization support.
// Automatically uses localization from context if available.
func Error(c *gin.Context, status int, err domain.AppError) {
	dict := getDict(c)
	var msg string

	if dict != nil {
		msg = TranslateAppError(dict, err)
	} else {
		// Fallback to error code if no dictionary
		msg = err.Code
	}

	c.JSON(status, Response{
		Code:      &err.Code,
		Message:   msg,
		Data:      nil,
		Errors:    nil,
		Meta:      nil,
		RequestID: getRequestID(c),
	})
}

// BadGateway sends a bad gateway error response (502 Bad Gateway).
// Can accept (hint string) for manual usage or (AppError) for localized usage.
// Automatically uses localization from context if available.
func BadGateway(c *gin.Context, args ...interface{}) {
	// Check if first arg is AppError
	if len(args) == 1 {
		if appErr, ok := args[0].(domain.AppError); ok {
			Error(c, 502, appErr)
			return
		}
	}

	// Legacy signature: hint string
	hint := ""
	if len(args) >= 1 {
		if h, ok := args[0].(string); ok {
			hint = h
		}
	}

	code := CodeBadGateway
	c.JSON(502, Response{
		Code:      &code,
		Message:   getMessage(c, "common.bad_gateway", "Bad gateway"),
		Data:      nil,
		Errors:    gin.H{"hint": hint},
		Meta:      nil,
		RequestID: getRequestID(c),
	})
}

// ServiceUnavailable sends a service unavailable error response (503 Service Unavailable).
// Can accept (retryAfter int) for manual usage or (AppError) for localized usage.
// Automatically uses localization from context if available.
func ServiceUnavailable(c *gin.Context, args ...interface{}) {
	// Check if first arg is AppError
	if len(args) == 1 {
		if appErr, ok := args[0].(domain.AppError); ok {
			Error(c, 503, appErr)
			return
		}
	}

	// Legacy signature: retryAfter int
	retryAfter := 0
	if len(args) >= 1 {
		if r, ok := args[0].(int); ok {
			retryAfter = r
		}
	}

	code := CodeServiceUnavailable
	c.JSON(503, Response{
		Code:      &code,
		Message:   getMessage(c, "common.service_unavailable", "Service temporarily unavailable"),
		Data:      nil,
		Errors:    gin.H{"retry_after": retryAfter},
		Meta:      nil,
		RequestID: getRequestID(c),
	})
}

// GatewayTimeout sends a gateway timeout error response (504 Gateway Timeout).
// Can accept (hint string) for manual usage or (AppError) for localized usage.
// Automatically uses localization from context if available.
func GatewayTimeout(c *gin.Context, args ...interface{}) {
	// Check if first arg is AppError
	if len(args) == 1 {
		if appErr, ok := args[0].(domain.AppError); ok {
			Error(c, 504, appErr)
			return
		}
	}

	// Legacy signature: hint string
	hint := ""
	if len(args) >= 1 {
		if h, ok := args[0].(string); ok {
			hint = h
		}
	}

	code := CodeGatewayTimeout
	c.JSON(504, Response{
		Code:      &code,
		Message:   getMessage(c, "common.gateway_timeout", "Gateway timeout"),
		Data:      nil,
		Errors:    gin.H{"hint": hint},
		Meta:      nil,
		RequestID: getRequestID(c),
	})
}

// TooManyRequests sends a rate limit exceeded error response (429 Too Many Requests).
// Can accept (message string) for manual usage or (AppError) for localized usage.
// Automatically uses localization from context if available.
func TooManyRequests(c *gin.Context, args ...interface{}) {
	// Check if first arg is AppError
	if len(args) == 1 {
		if appErr, ok := args[0].(domain.AppError); ok {
			Error(c, 429, appErr)
			return
		}
	}

	// Legacy signature: message string
	message := "Too many requests"
	if len(args) >= 1 {
		if m, ok := args[0].(string); ok {
			message = m
		}
	}

	code := CodeTooManyRequests
	c.JSON(429, Response{
		Code:      &code,
		Message:   getMessage(c, "common.too_many_requests", message),
		Data:      nil,
		Errors:    nil,
		Meta:      nil,
		RequestID: getRequestID(c),
	})
}

// RequestTimeout sends a request timeout error response (408 Request Timeout).
// Can accept (message string) for manual usage or (AppError) for localized usage.
// Automatically uses localization from context if available.
func RequestTimeout(c *gin.Context, args ...interface{}) {
	// Check if first arg is AppError
	if len(args) == 1 {
		if appErr, ok := args[0].(domain.AppError); ok {
			Error(c, 408, appErr)
			return
		}
	}

	// Legacy signature: message string
	message := "Request timeout"
	if len(args) >= 1 {
		if m, ok := args[0].(string); ok {
			message = m
		}
	}

	code := CodeRequestTimeout
	c.JSON(408, Response{
		Code:      &code,
		Message:   getMessage(c, "common.request_timeout", message),
		Data:      nil,
		Errors:    nil,
		Meta:      nil,
		RequestID: getRequestID(c),
	})
}
