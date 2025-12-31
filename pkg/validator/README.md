# Validator Package

Struct field validation using go-playground/validator/v10.

## Features

- go-playground/validator/v10 integration
- 100+ built-in validation rules
- Tag-based validation
- Returns error codes for localization
- Backward compatible API

## Usage

### Basic Validation

Define validation rules using struct tags:

```go
type RegisterInput struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
    Name     string `json:"name" validate:"required,min=2"`
}

func handler(c *gin.Context) {
    var input RegisterInput
    if err := c.ShouldBindJSON(&input); err != nil {
        return err
    }
    
    // New: Returns []domain.AppError for localization
    errs := validator.ValidateStruct(input)
    if len(errs) > 0 {
        // Use with localization
        response.ValidationErrorI18n(c, dict, errs)
        return
    }
}
```

## Validation Rules

### Common Rules

- `required` - Field must not be empty
- `email` - Valid email format
- `min=N` - Minimum length (string) or value (number)
- `max=N` - Maximum length (string) or value (number)
- `len=N` - Exact length
- `eq=value` - Equal to value
- `ne=value` - Not equal to value
- `gt=N` - Greater than
- `gte=N` - Greater than or equal
- `lt=N` - Less than
- `lte=N` - Less than or equal

### String Validation

- `alpha` - Alphabetic characters only
- `alphanum` - Alphanumeric characters only
- `numeric` - Numeric characters only
- `email` - Valid email address
- `url` - Valid URL
- `uri` - Valid URI
- `uuid` - Valid UUID
- `uuid3` - Valid UUID v3
- `uuid4` - Valid UUID v4
- `uuid5` - Valid UUID v5

### Number Validation

- `min=N` - Minimum value
- `max=N` - Maximum value
- `eq=N` - Equal to value
- `ne=N` - Not equal to value
- `gt=N` - Greater than
- `gte=N` - Greater than or equal
- `lt=N` - Less than
- `lte=N` - Less than or equal

### Multiple Rules

Chain multiple rules with comma:

```go
type Input struct {
    Email    string `validate:"required,email"`
    Password string `validate:"required,min=8,max=72"`
    Age      int    `validate:"required,gte=18,lte=120"`
}
```

## Error Handling

### With Localization

ValidateStruct returns error codes that can be translated:

```go
errs := validator.ValidateStruct(input)
if len(errs) > 0 {
    // errs[0].Code = "validation.required"
    // errs[0].Field = "email"
    // errs[0].Params = map[string]string{}
    
    // Translator will:
    // 1. Lookup "fields.email" -> "Email"
    // 2. Lookup "validation.required" -> "{field} is required"
    // 3. Replace {field} with "Email"
    // 4. Result: "Email is required"
    
    dict := middleware.GetDict(c)
    response.ValidationErrorI18n(c, dict, errs)
    return
}
```

### Backward Compatible

Old Validate function still works:

```go
err := validator.Validate(input)
if err != nil {
    // Returns first error as error interface
    return err
}
```

## Examples

### Complete Handler with Localization

```go
func (h *Handler) Register(c *gin.Context) {
    dict := middleware.GetDict(c)
    
    var input RegisterInput
    if err := c.ShouldBindJSON(&input); err != nil {
        response.BadRequestI18n(c, dict, domain.NewAppError("request.invalid_json"))
        return
    }
    
    errs := validator.ValidateStruct(input)
    if len(errs) > 0 {
        response.ValidationErrorI18n(c, dict, errs)
        return
    }
    
    // Process valid input
    user, err := h.userUseCase.Register(c.Request.Context(), input)
    if err != nil {
        response.InternalServerErrorI18n(c, dict, domain.NewAppError("common.error"))
        return
    }
    
    response.CreatedI18n(c, dict, "user.created", user)
}
```

## Error Code Format

Validation errors are converted to generic error codes with field label lookup:

```
validation.{rule}
```

Examples:
- `validation.required` + field="email" → "Email is required"
- `validation.min` + field="password" + params={"min":"8"} → "Password must be at least 8 characters"
- `validation.email` + field="email" → "Email must be a valid email address"

The translator automatically:
1. Gets field label from `fields.{field_name}` in translation JSON
2. Injects field label as `{field}` parameter
3. Replaces all parameters in validation message template

This approach makes translations scalable - add 100 fields without adding 100 translations.

## Supported Field Types

- String
- Int, Int8, Int16, Int32, Int64
- Uint, Uint8, Uint16, Uint32, Uint64
- Float32, Float64
- Bool
- Struct (nested validation)
- Slice, Array
- Map
- Pointer

## Full Documentation

For complete list of validation rules and features, see:
https://pkg.go.dev/github.com/go-playground/validator/v10

