package response

// HTTP Response Codes
// These are application-level codes, not HTTP status codes.
// Used in Response.Code field for categorizing errors.
const (
	// Success codes

	CodeSuccess = "" // Empty for success responses

	// Client error codes (4xx equivalent)

	CodeValidationError = "VALIDATION_ERROR"
	CodeBadRequest      = "BAD_REQUEST"
	CodeUnauthorized    = "UNAUTHORIZED"
	CodeForbidden       = "FORBIDDEN"
	CodeNotFound        = "NOT_FOUND"
	CodeConflict        = "CONFLICT"
	CodeInvalidInput    = "INVALID_INPUT"
	CodeInvalidJSON     = "INVALID_JSON"

	// Server error codes (5xx equivalent)

	CodeInternalServerError = "INTERNAL_SERVER_ERROR"
	CodeBadGateway          = "BAD_GATEWAY"
	CodeServiceUnavailable  = "SERVICE_UNAVAILABLE"
	CodeGatewayTimeout      = "GATEWAY_TIMEOUT"

	// Special codes

	CodeRedirect = "REDIRECT"
)

// Note: Domain-specific error codes (e.g., "user.not_found", "auth.invalid_credentials")
// are defined in domain.AppError and used for localization purposes.
// These constants are for HTTP-level response categorization only.
