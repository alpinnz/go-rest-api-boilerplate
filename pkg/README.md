# pkg/

Public reusable packages for cross-project utilities.

## Directory Structure

```
pkg/
├── auth/           Authentication utilities (JWT, password hashing)
├── response/       HTTP response formatting
└── validator/      Struct field validation
```

## Characteristics

- Reusable across projects
- No internal dependencies
- Well-documented public API
- Independent packages

## Packages

### auth/

Authentication utilities including JWT token management and password hashing.

**Features:**
- JWT token generation and validation
- Bcrypt password hashing
- Secure password verification

[View detailed documentation](./auth/README.md)

### response/

Standardized HTTP response formatting for REST API.

**Features:**
- Consistent response structure
- Success and error responses
- Validation error handling
- Pagination and metadata support

[View detailed documentation](./response/README.md)

### validator/

Struct field validation using reflection and tags.

**Features:**
- Tag-based validation rules
- Required field validation
- Email format validation
- Minimum length validation

[View detailed documentation](./validator/README.md)

## Usage Guidelines

### Import Packages

```go
import (
    "github.com/alpinnz/go-rest-api-boilerplate/pkg/auth"
    "github.com/alpinnz/go-rest-api-boilerplate/pkg/response"
    "github.com/alpinnz/go-rest-api-boilerplate/pkg/validator"
)
```

### Package Independence

Each package is independent and can be used separately:

```go
// Use auth package only
import "github.com/alpinnz/go-rest-api-boilerplate/pkg/auth"

token, err := auth.GenerateToken(userID)
```

### Cross-Package Integration

Packages can be combined for complete functionality:

```go
func (h *Handler) Register(c *gin.Context) {
    var input RegisterInput
    
    // Bind JSON
    if err := c.ShouldBindJSON(&input); err != nil {
        response.BadRequest(c, "INVALID_JSON", "Invalid request body", nil)
        return
    }
    
    // Validate input
    if err := validator.Validate(input); err != nil {
        response.ValidationError(c, buildValidationErrors(err))
        return
    }
    
    // Hash password
    hashedPassword, err := auth.HashPassword(input.Password)
    if err != nil {
        response.InternalServerError(c, "Failed to process password")
        return
    }
    
    // Save user and generate token
    user := saveUser(input.Email, hashedPassword, input.Name)
    token, err := auth.GenerateToken(user.ID)
    if err != nil {
        response.InternalServerError(c, "Failed to generate token")
        return
    }
    
    response.Created(c, gin.H{
        "user":  user,
        "token": token,
    })
}
```

## Examples

### Authentication Flow

```go
// Hash password
hashedPassword, err := auth.HashPassword("password123")

// Verify password
isValid := auth.CheckPasswordHash("password123", hashedPassword)

// Generate JWT token
token, err := auth.GenerateToken(12345)

// Validate token
claims, err := auth.ValidateToken(token)
userID := claims.UserID
```

## response/

HTTP response formatting utilities.

### Files

- `response.go` - Standardized response structure

### Functions

```go
// Success responses
func Success(c *gin.Context, data interface{})
func Created(c *gin.Context, data interface{})
func NoContent(c *gin.Context)

// Error responses
func BadRequest(c *gin.Context, code, message string, errors interface{})
func Unauthorized(c *gin.Context, message string)
func Forbidden(c *gin.Context, resource string)
func NotFound(c *gin.Context, resource string)
func Conflict(c *gin.Context, field, message string)
func InternalServerError(c *gin.Context, message string)
func ValidationError(c *gin.Context, errors ValidationErrors)
```

### Response Structure

```json
{
  "code": "ERROR_CODE",
  "message": "Human readable message",
  "data": {},
  "errors": {},
  "meta": {},
  "trace_id": "abc123"
}
```

### Usage

```go
// Success
response.Success(c, user)

// Error
response.NotFound(c, "user")
response.BadRequest(c, "INVALID_INPUT", "Invalid data", nil)
```

## validator/

Input validation utilities.

### Files

- `validator.go` - Struct field validation

### Tag-Based Validation

```go
type RegisterInput struct {
    Email    string `validate:"required,email"`
    Password string `validate:"required,min=8"`
    Name     string `validate:"required"`
}

// Validate struct
if err := validator.Validate(input); err != nil {
    return err
}
```

### Available Rules

- `required` - Field must not be empty
- `email` - Valid email format
- `min=N` - Minimum string length

## Adding New Package

Create new directory under `pkg/`:

```bash
mkdir pkg/cache
touch pkg/cache/cache.go
```

```go
// pkg/cache/cache.go
package cache

// Cache provides caching interface.
type Cache interface {
    Get(key string) (interface{}, error)
    Set(key string, value interface{}, ttl time.Duration) error
    Delete(key string) error
}
```

## Best Practices

1. **Public API**: Only export what's necessary
2. **Documentation**: Document all exported functions
3. **Testing**: Write comprehensive tests
4. **Examples**: Add usage examples in docs
5. **Independence**: No dependencies on `internal/`

## References

- [Go Package Documentation](https://go.dev/doc/effective_go#commentary)
- [Designing Go Libraries](https://abhinavg.net/2022/12/06/designing-go-libraries/)

