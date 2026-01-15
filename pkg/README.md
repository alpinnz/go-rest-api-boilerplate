# Package (pkg)

Reusable packages that are safe to import from anywhere (including from other projects).

Guiding rule:
- `pkg/*` must not import `internal/*`.
- Keep these packages small, focused, and easy to test.

## Structure

```
pkg/
├── auth/             # JWT + password hashing helpers
├── context/          # Request-scoped context helpers (request_id, user_id)
├── database/         # Transaction helpers
├── errors/           # Typed/structured error helpers
├── logger/           # Structured logging (Zerolog)
├── migration/        # DB migration runner
├── pagination/       # Pagination helpers
├── response/         # HTTP response helpers (usable by delivery layer)
└── validator/        # Request validation utilities
```

## Packages

### auth

JWT token management and password hashing with separate access/refresh support.

Usage:
```text
// Configure JWT settings (at app initialization)
auth.SetJWTConfig(
    accessSecret, accessExpiration,
    refreshSecret, refreshExpiration,
)

accessToken, err := auth.GenerateAccessToken(userID)
refreshToken, err := auth.GenerateRefreshToken(userID)

claims, err := auth.ValidateToken(token)

hashedPassword, err := auth.HashPassword(password)
match := auth.CheckPasswordHash(password, hashedPassword)
```

### context

Context helpers for request tracing across layers.

Usage:
```text
ctx := context.WithRequestID(ctx, requestID)
requestID := context.GetRequestID(ctx)

ctx = context.WithUserID(ctx, userID)
```

### database

Transaction helpers.

Usage:
```text
err := database.WithTransaction(ctx, db, func(ctx context.Context, tx *sql.Tx) error {
    // ...
    return nil
})
```

### errors

Structured error helpers with optional context.

Usage:
```text
err := errors.NotFound("User not found").WithContext("user_id", userID)
err = errors.Wrap(dbErr, "DB_ERROR", "Database query failed")
```

### logger

Structured logging with zerolog.

Usage:
```text
log := logger.New(logger.Config{Level: "info", Pretty: false})
log.Info().Str("user_id", userID).Msg("created")
```

### migration

Database migration runner.

Typical usage:
- CLI/Makefile runs migrations (recommended)
- Application boot can also run migrations depending on your deployment strategy

Usage:
```text
runner := migration.NewRunner(db, "migrations")
_ = runner.Initialize()
_ = runner.Up()
```

### pagination

Pagination helpers for list endpoints.

### response

HTTP response helpers (commonly used by `internal/delivery/http/*`).

### validator

Input validation using go-playground/validator.

Usage:
```text
validationErrors := validator.Validate(req)
if len(validationErrors) > 0 {
    response.ValidationError(c, validationErrors)
    return
}
```

Built-in validations depend on `go-playground/validator` tags (e.g. `required`, `email`, `min`, `max`, `uuid`, `oneof`).

## Testing

These packages have unit tests under `pkg/*`:
- `pkg/auth/auth_test.go`
- `pkg/errors/errors_test.go`
- `pkg/pagination/pagination_test.go`
- `pkg/validator/validator_test.go`

Run tests using the root Makefile (`make test`) to keep a consistent workflow.
