package response

// Code defines a type for API response codes
type Code string

// Available response codes
const (
	CodeSuccess         Code = "SUCCESS"
	CodeBadRequest      Code = "BAD_REQUEST"
	CodeNotFound        Code = "NOT_FOUND"
	CodeUnauthorized    Code = "UNAUTHORIZED"
	CodeInternalError   Code = "INTERNAL_ERROR"
	CodeDatabaseError   Code = "DATABASE_ERROR"
	CodeValidationError Code = "VALIDATION_ERROR"
	CodeForbidden       Code = "FORBIDDEN"
	CodeConflict        Code = "CONFLICT"
)
