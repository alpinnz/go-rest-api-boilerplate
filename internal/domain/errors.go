package domain

import "errors"

// Domain-specific errors for consistent error handling across layers.
// These errors should be mapped to appropriate HTTP status codes in handlers.
//
// Usage:
//
//	if user == nil {
//	    return ErrNotFound
//	}
var (
	// ErrNotFound indicates requested resource does not exist.
	// Maps to HTTP 404 Not Found.
	ErrNotFound = errors.New("resource not found")

	// ErrConflict indicates resource already exists (duplicate).
	// Commonly used for unique constraint violations (e.g., duplicate email).
	// Maps to HTTP 409 Conflict.
	ErrConflict = errors.New("resource already exists")

	// ErrInvalidInput indicates input validation failed.
	// Used for business logic validation errors.
	// Maps to HTTP 400 Bad Request.
	ErrInvalidInput = errors.New("invalid input")

	// ErrUnauthorized indicates missing or invalid authentication.
	// Used when user is not authenticated.
	// Maps to HTTP 401 Unauthorized.
	ErrUnauthorized = errors.New("unauthorized")

	// ErrForbidden indicates user lacks permission for requested action.
	// Used when user is authenticated but not authorized.
	// Maps to HTTP 403 Forbidden.
	ErrForbidden = errors.New("forbidden")

	// ErrInternalServer indicates unexpected server error.
	// Used for unrecoverable errors that shouldn't be exposed to client.
	// Maps to HTTP 500 Internal Server Error.
	ErrInternalServer = errors.New("internal server error")

	// ErrBadRequest indicates malformed request.
	// Used for syntax errors in request (e.g., invalid JSON).
	// Maps to HTTP 400 Bad Request.
	ErrBadRequest = errors.New("bad request")

	// ErrInvalidCredentials indicates authentication failed.
	// Used when username/password combination is incorrect.
	// Maps to HTTP 401 Unauthorized.
	ErrInvalidCredentials = errors.New("invalid credentials")

	// ErrSessionExpired indicates session token expired or invalid.
	// Used when token doesn't exist in Redis (expired or never created).
	// Maps to HTTP 401 Unauthorized.
	ErrSessionExpired = errors.New("session expired")
)
