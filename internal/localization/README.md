# Localization

Multi-language support for REST API responses using code-based translation system.

## Structure

```
localization/
├── types.go       # Dictionary type definition
├── loader.go      # Bundle loader
└── lang/
    ├── en.json    # English translations
    └── id.json    # Indonesian translations
```

## Architecture

### Translation Flow

1. Client sends Accept-Language header
2. Language middleware extracts and injects dictionary into context
3. Validator returns error codes
4. Response layer translates codes to localized messages
5. Client receives response in requested language

### Layer Responsibilities

**Domain/Usecase**: Return error codes only, no language strings
**Validator**: Convert validation errors to error codes
**Localization**: Load and provide translation dictionaries
**Response**: Translate error codes to messages
**Handler**: Extract language and use localized responses

## Error Code Format

Error codes follow dot notation for clear organization:

```
validation.{field}.{rule}    # Field-specific validation errors
user.{action}                # User domain errors
auth.{action}                # Authentication errors
common.{type}                # Common messages
```

### Validation Errors

Pattern: `validation.{field_name}.{rule_name}`

Examples:
- `validation.email.required` → "Email is required"
- `validation.email.email` → "Email format is invalid"
- `validation.password.min` → "Password must be at least {min} characters"
- `validation.name.required` → "Name is required"

**Parameter Interpolation**: Use `{param_name}` in messages, replaced by validator params.

### Domain Errors

Examples:
- `user.not_found` → "User not found"
- `user.email_exists` → "Email already registered"
- `auth.invalid_credentials` → "Invalid email or password"
- `auth.invalid_refresh_token` → "Invalid or expired refresh token"
- `common.error` → "Unexpected error occurred"

## Usage

### Adding Translations

Edit `lang/en.json` or `lang/id.json`:

```json
{
  "fields": {
    "email": "Email",
    "password": "Password",
    "name": "Name",
    "first_name": "First Name",
    "last_name": "Last Name"
  },
  "validation": {
    "required": "{field} is required",
    "email": "{field} must be a valid email address",
    "min": "{field} must be at least {min} characters",
    "max": "{field} must not exceed {max} characters"
  },
  "user": {
    "not_found": "User not found"
  }
}
```

### Field Label Mapping

Technical field names are mapped to human-readable labels:

```
first_name -> "First Name"
last_name  -> "Last Name"
email      -> "Email"
```

Translator automatically:
1. Looks up field label from `fields.{field_name}`
2. Uses label in validation message template
3. Falls back to field name if label not found

### Parameter Interpolation

Use `{key}` syntax for dynamic values:

```json
{
  "validation": {
    "min": "{field} must be at least {min} characters",
    "between": "{field} must be between {min} and {max}"
  }
}
```

Error with params:
```go
domain.AppError{
    Code: "validation.min",
    Field: "password",
    Params: map[string]string{"min": "8"},
}
```

Translation process:
1. Lookup field label: `fields.password` = "Password"
2. Lookup validation message: `validation.min` = "{field} must be at least {min} characters"
3. Replace `{field}` with "Password"
4. Replace `{min}` with "8"
5. Result: "Password must be at least 8 characters"

## Adding New Language

1. Create `lang/{code}.json` with all translation keys
2. Load in router: `bundle.Load("{code}", "internal/localization/lang/{code}.json")`
3. Test with `Accept-Language: {code}` header

## Bundle API

```go
// Create bundle with default language
bundle := localization.NewBundle("en")

// Load language files
bundle.Load("en", "internal/localization/lang/en.json")
bundle.Load("id", "internal/localization/lang/id.json")

// Get dictionary for language
dict := bundle.Get("id")

// Get default language
defaultLang := bundle.DefaultLang()
```

## Dictionary Type

```go
type Dictionary map[string]any
```

Supports nested structure matching JSON format. Access via dot notation path in translator functions.

