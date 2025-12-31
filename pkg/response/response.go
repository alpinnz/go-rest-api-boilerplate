// Package response provides HTTP response formatting utilities for consistent API responses.
package response

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Response represents standardized API response structure.
type Response struct {
	Code    *string     `json:"code"`     // Application-level code or error type
	Message string      `json:"message"`  // Human-readable message
	Data    interface{} `json:"data"`     // Response payload
	Errors  interface{} `json:"errors"`   // Validation or contextual error details
	Meta    interface{} `json:"meta"`     // Metadata (pagination, sorting, etc.)
	TraceID *string     `json:"trace_id"` // Optional request trace identifier
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

// Pagination represents pagination metadata.
type Pagination struct {
	Page      int `json:"page"`
	PerPage   int `json:"per_page"`
	TotalData int `json:"total_data"`
	TotalPage int `json:"total_page"`
}

// Sort represents sorting criteria.
type Sort struct {
	Field string `json:"field"`
	Order string `json:"order"`
}

// Meta represents response metadata.
type Meta struct {
	Pagination *Pagination `json:"pagination,omitempty"`
	Sort       []Sort      `json:"sort,omitempty"`
}

// getTraceID retrieves or generates trace ID from context.
func getTraceID(c *gin.Context) *string {
	if traceID, exists := c.Get("trace_id"); exists {
		if id, ok := traceID.(string); ok {
			return &id
		}
	}
	id := uuid.New().String()[:6]
	return &id
}

// Success sends a successful JSON response (200 OK).
func Success(c *gin.Context, data interface{}) {
	c.JSON(200, Response{
		Code:    nil,
		Message: "Success",
		Data:    data,
		Errors:  nil,
		Meta:    nil,
		TraceID: getTraceID(c),
	})
}

// SuccessWithMeta sends a successful JSON response with metadata (200 OK).
func SuccessWithMeta(c *gin.Context, data interface{}, meta *Meta) {
	c.JSON(200, Response{
		Code:    nil,
		Message: "Success",
		Data:    data,
		Errors:  nil,
		Meta:    meta,
		TraceID: getTraceID(c),
	})
}

// Created sends a resource created response (201 Created).
func Created(c *gin.Context, data interface{}) {
	c.JSON(201, Response{
		Code:    nil,
		Message: "Resource created successfully",
		Data:    data,
		Errors:  nil,
		Meta:    nil,
		TraceID: getTraceID(c),
	})
}

// NoContent sends a no content response (204 No Content).
func NoContent(c *gin.Context) {
	c.JSON(204, Response{
		Code:    nil,
		Message: "No Content",
		Data:    nil,
		Errors:  nil,
		Meta:    nil,
		TraceID: getTraceID(c),
	})
}

// Redirect sends a redirect response (301 Moved Permanently).
func Redirect(c *gin.Context, location string) {
	code := "REDIRECT"
	c.JSON(301, Response{
		Code:    &code,
		Message: "Resource moved permanently",
		Data:    gin.H{"location": location},
		Errors:  nil,
		Meta:    nil,
		TraceID: getTraceID(c),
	})
}

// TemporaryRedirect sends a temporary redirect response (302 Found).
func TemporaryRedirect(c *gin.Context, location string) {
	code := "REDIRECT"
	c.JSON(302, Response{
		Code:    &code,
		Message: "Resource temporarily moved",
		Data:    gin.H{"location": location},
		Errors:  nil,
		Meta:    nil,
		TraceID: getTraceID(c),
	})
}

// ValidationError sends a validation error response (400 Bad Request).
func ValidationError(c *gin.Context, errors ValidationErrors) {
	code := "VALIDATION_ERROR"
	c.JSON(400, Response{
		Code:    &code,
		Message: "Validation failed",
		Data:    nil,
		Errors:  errors,
		Meta:    nil,
		TraceID: getTraceID(c),
	})
}

// BadRequest sends a bad request error response (400 Bad Request).
func BadRequest(c *gin.Context, code string, message string, errors interface{}) {
	c.JSON(400, Response{
		Code:    &code,
		Message: message,
		Data:    nil,
		Errors:  errors,
		Meta:    nil,
		TraceID: getTraceID(c),
	})
}

// Unauthorized sends an unauthorized error response (401 Unauthorized).
func Unauthorized(c *gin.Context, reason string) {
	code := "UNAUTHORIZED"
	c.JSON(401, Response{
		Code:    &code,
		Message: "Unauthorized",
		Data:    nil,
		Errors:  gin.H{"reason": reason},
		Meta:    nil,
		TraceID: getTraceID(c),
	})
}

// Forbidden sends a forbidden error response (403 Forbidden).
func Forbidden(c *gin.Context, permission string) {
	code := "FORBIDDEN"
	c.JSON(403, Response{
		Code:    &code,
		Message: "Access denied",
		Data:    nil,
		Errors:  gin.H{"permission": permission},
		Meta:    nil,
		TraceID: getTraceID(c),
	})
}

// NotFound sends a not found error response (404 Not Found).
func NotFound(c *gin.Context, resource string) {
	code := "NOT_FOUND"
	c.JSON(404, Response{
		Code:    &code,
		Message: "Resource not found",
		Data:    nil,
		Errors:  gin.H{"resource": resource},
		Meta:    nil,
		TraceID: getTraceID(c),
	})
}

// Conflict sends a conflict error response (409 Conflict).
func Conflict(c *gin.Context, field string, message string) {
	code := "CONFLICT"
	c.JSON(409, Response{
		Code:    &code,
		Message: "Resource conflict",
		Data:    nil,
		Errors:  gin.H{"field": field, "message": message},
		Meta:    nil,
		TraceID: getTraceID(c),
	})
}

// InternalServerError sends an internal server error response (500 Internal Server Error).
func InternalServerError(c *gin.Context, hint string) {
	code := "INTERNAL_SERVER_ERROR"
	c.JSON(500, Response{
		Code:    &code,
		Message: "Unexpected error occurred",
		Data:    nil,
		Errors:  gin.H{"hint": hint},
		Meta:    nil,
		TraceID: getTraceID(c),
	})
}

// BadGateway sends a bad gateway error response (502 Bad Gateway).
func BadGateway(c *gin.Context, hint string) {
	code := "BAD_GATEWAY"
	c.JSON(502, Response{
		Code:    &code,
		Message: "Bad gateway",
		Data:    nil,
		Errors:  gin.H{"hint": hint},
		Meta:    nil,
		TraceID: getTraceID(c),
	})
}

// ServiceUnavailable sends a service unavailable error response (503 Service Unavailable).
func ServiceUnavailable(c *gin.Context, retryAfter int) {
	code := "SERVICE_UNAVAILABLE"
	c.JSON(503, Response{
		Code:    &code,
		Message: "Service temporarily unavailable",
		Data:    nil,
		Errors:  gin.H{"retry_after": retryAfter},
		Meta:    nil,
		TraceID: getTraceID(c),
	})
}

// GatewayTimeout sends a gateway timeout error response (504 Gateway Timeout).
func GatewayTimeout(c *gin.Context, hint string) {
	code := "GATEWAY_TIMEOUT"
	c.JSON(504, Response{
		Code:    &code,
		Message: "Gateway timeout",
		Data:    nil,
		Errors:  gin.H{"hint": hint},
		Meta:    nil,
		TraceID: getTraceID(c),
	})
}
