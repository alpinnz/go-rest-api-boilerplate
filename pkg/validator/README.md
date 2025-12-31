# Validator Package

Struct field validation using reflection and custom tags.

## Features

- Tag-based validation
- Required field validation
- Email format validation
- Minimum length validation
- Extensible rule system

## Usage

### Basic Validation

Define validation rules using struct tags:

```go
type RegisterInput struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
    Name     string `json:"name" validate:"required"`
}

func handler(c *gin.Context) {
    var input RegisterInput
    if err := c.ShouldBindJSON(&input); err != nil {
        return err
    }
    
    if err := validator.Validate(input); err != nil {
        return err
    }
}
```

## Validation Rules

### required

Validates that field is not empty or zero value:

```go
type User struct {
    Email string `validate:"required"`
    Age   int    `validate:"required"`
}
```

Checks:
- String: not empty
- Int: not zero
- Bool: not false

### email

Validates email format:

```go
type Input struct {
    Email string `validate:"required,email"`
}
```

Pattern: `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

Examples:
- Valid: `user@example.com`, `test.user@domain.co.uk`
- Invalid: `invalid`, `@example.com`, `user@.com`

### min

Validates minimum string length:

```go
type Input struct {
    Password string `validate:"required,min=8"`
    Username string `validate:"required,min=3"`
}
```

### Multiple Rules

Chain multiple rules with comma:

```go
type Input struct {
    Email    string `validate:"required,email"`
    Password string `validate:"required,min=8"`
}
```

## Error Messages

Validation returns descriptive error messages:

```go
err := validator.Validate(input)
if err != nil {
    // Error examples:
    // "Email is required"
    // "Email must be a valid email"
    // "Password must be at least 8 characters"
}
```

## Integration Example

### With Response Package

```go
import (
    "github.com/alpinnz/go-rest-api-boilerplate/pkg/validator"
    "github.com/alpinnz/go-rest-api-boilerplate/pkg/response"
)

func (h *Handler) Register(c *gin.Context) {
    var input RegisterInput
    if err := c.ShouldBindJSON(&input); err != nil {
        response.BadRequest(c, "INVALID_JSON", "Invalid request body", gin.H{
            "details": err.Error(),
        })
        return
    }
    
    if err := validator.Validate(input); err != nil {
        errors := response.ValidationErrors{
            Body: []response.ValidationItem{
                {
                    Code:    "VALIDATION_FAILED",
                    Field:   "unknown",
                    Message: err.Error(),
                },
            },
        }
        response.ValidationError(c, errors)
        return
    }
    
    // Process valid input
}
```

### Custom Validation

Combine with custom business logic:

```go
func (h *Handler) CreateUser(c *gin.Context) {
    var input CreateUserInput
    
    if err := c.ShouldBindJSON(&input); err != nil {
        return handleError(err)
    }
    
    if err := validator.Validate(input); err != nil {
        return handleError(err)
    }
    
    // Additional custom validation
    if h.userExists(input.Email) {
        response.Conflict(c, "email", "Email already exists")
        return
    }
}
```

## Supported Field Types

- `string`: required, email, min
- `int`, `int8`, `int16`, `int32`, `int64`: required
- `uint`, `uint8`, `uint16`, `uint32`, `uint64`: required
- `float32`, `float64`: required
- `bool`: required

## Limitations

Current implementation:
- Basic validation rules only
- Simple error messages
- No nested struct validation
- No custom error message support

For complex validation scenarios, consider using libraries like:
- `github.com/go-playground/validator/v10`
- `github.com/asaskevich/govalidator`

