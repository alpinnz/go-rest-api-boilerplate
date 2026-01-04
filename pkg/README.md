# Package (pkg)

Reusable packages that can be used across the application or in other projects.

## Structure

```
pkg/
├── auth/             # Authentication utilities (JWT, password hashing)
├── context/          # Context helpers for request tracing
├── database/         # Database transaction helpers
├── errors/           # Structured error handling
├── logger/           # Structured logging with Zerolog
├── pagination/       # Pagination utilities
├── response/         # HTTP response formatting
└── validator/        # Input validation utilities
```

## Packages

### auth

JWT token management and password hashing.

Usage:
```go
// Generate JWT token
token, err := auth.GenerateToken(userID, secret, expiration)

// Validate token
claims, err := auth.ValidateToken(token, secret)

// Hash password
hashedPassword, err := auth.HashPassword(password)

// Verify password
err := auth.VerifyPassword(hashedPassword, password)
```

### context

Context helpers for request tracing across layers.

Usage:
```go
// Set request ID
ctx := context.WithRequestID(ctx, requestID)

// Get request ID
requestID := context.GetRequestID(ctx)

// Set user ID
ctx = context.WithUserID(ctx, userID)
```

### database

Transaction management helpers.

Usage:
```go
// Execute operations in transaction
err := database.WithTransaction(ctx, db, func(ctx context.Context, tx *sql.Tx) error {
    // Operations here
    return nil
})
```

Features:
- Automatic rollback on error
- Panic recovery with rollback

### errors

Structured error handling with context.

Usage:
```go
// Create specific errors
err := errors.NotFound("User not found")
err := errors.BadRequest("Invalid email")
err := errors.Unauthorized("Token expired")

// Add context
err := errors.NotFound("User not found").
    WithContext("user_id", userID)
```
- `BadRequest()` - 400 errors
- `Unauthorized()` - 401 errors
- `Forbidden()` - 403 errors
- `Conflict()` - 409 errors
- `InternalServer()` - 500 errors
- `ValidationFailed()` - Validation errors
- `Timeout()` - 408 timeout errors

**Usage:**
```go
// Create error
err := errors.NotFound("User not found")

// Add context
err = errors.NotFound("User not found").WithContext("user_id", userID)

// Wrap existing error
err = errors.Wrap(dbErr, "DB_ERROR", "Database query failed")
```

### logger

Structured logging with Zerolog.

Usage:
```go
// Create logger
log := logger.Development()  // Pretty console output
log := logger.Production()   // JSON structured logs

// With fields
log.WithFields(map[string]interface{}{
    "user_id": userID,
}).Info("User logged in")
```

### pagination

Pagination utilities for list endpoints.

Usage:
```go
// Extract from query params
params := pagination.FromContext(c)

// Query with pagination
users, total, err := repo.FindAll(ctx, params.Offset, params.PerPage)

// Create metadata
meta := pagination.NewMeta(params.Page, params.PerPage, total)

// Return with metadata
response.SuccessWithMeta(c, users, &meta)
```

Query Parameters: `?page=2&per_page=20`

### response

HTTP response formatting with localization support.

Response Structure:
```json
{
  "code": "SUCCESS",
  "message": "Request successful",
  "data": {},
  "meta": {}
}
```

Usage:
```go
// Success responses
response.Success(c, data)
response.Created(c, newUser)
response.SuccessWithMeta(c, users, &meta)

// Error responses
response.NotFound(c, "User not found")
response.BadRequest(c, "INVALID_EMAIL", "Invalid email", errors)
response.Unauthorized(c, "Token expired")
response.ValidationError(c, validationErrors)
```

### validator

Input validation using go-playground/validator.

Usage:
```go
// Define validation in struct
type CreateUserRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Username string `json:"username" validate:"required,min=3,max=50"`
    Password string `json:"password" validate:"required,min=8"`
}

// Validate
validationErrors := validator.Validate(req)
if len(validationErrors) > 0 {
    response.ValidationError(c, validationErrors)
    return
}
```

Built-in Validations: required, email, min, max, len, uuid, oneof

## Testing

Run package tests:
```bash
go test ./pkg/...              # All packages
go test -v ./pkg/pagination/   # Specific package
go test -cover ./pkg/...       # With coverage
```

## Design Principles

Independence:
- No dependencies on internal packages
- Can be used in other projects
- Clear, simple interfaces

Single Responsibility:
- Each package has one well-defined purpose

No Side Effects:
- Predictable behavior
- Easy to test
- No global state


