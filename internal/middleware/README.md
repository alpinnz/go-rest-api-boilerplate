# Middleware Package

HTTP middleware for request processing and cross-cutting concerns.

## Overview

Middleware package provides reusable HTTP middleware for Gin framework. Each middleware handles specific cross-cutting concerns like authentication, CORS, logging, localization, and error recovery.

## Available Middleware

### Logger

Logs HTTP requests with automatic request ID generation for tracing.

**File:** `logger.go`

**Usage:**

```go
r.Use(middleware.Logger())
```

**Behavior:**

1. Automatically generates unique UUID as request_id
2. Sets request_id in context for handlers to access  
3. Logs all requests with structured format including client IP and user agent
4. Available throughout request lifecycle
5. Logs errors separately for better visibility

**Log Output Example:**

```
[request_id=550e8400-e29b-41d4-a716-446655440000] [ip=192.168.1.1] [method=GET] [path=/api/v1/users/me] [status=200] [latency=15ms] [user_agent=Mozilla/5.0]
[request_id=abc12345-6789-0def-ghij-klmnopqrstuv] [ip=192.168.1.2] [method=POST] [path=/api/v1/auth/login] [status=401] [latency=5ms] [user_agent=curl/7.64.1]
[request_id=abc12345-6789-0def-ghij-klmnopqrstuv] [errors] validation failed: email is required
```

**Access in Handler:**

```go
import "github.com/alpinnz/go-rest-api-boilerplate/internal/middleware"

func (h *Handler) MyHandler(c *gin.Context) {
    requestID := middleware.GetRequestID(c)
    log.Printf("[request_id=%s] Processing request", requestID)
}
```

**Response Includes request_id:**

```json
{
  "code": null,
  "message": "Success",
  "data": { ... },
  "request_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

**Benefits:**

- Automatic request tracing
- Structured log format for easy parsing
- Client IP tracking for security audit
- User agent tracking for debugging client issues
- Separate error logging for quick issue identification
- Correlate logs with specific requests
- Debug production issues easily
- No manual ID management needed
- Consistent across all requests

### Language

Extracts Accept-Language header and injects localization dictionary into context.

**File:** `language.go`

**Usage:**

```go
bundle := localization.NewBundle("en")
bundle.Load("en", "internal/localization/lang/en.json")
bundle.Load("id", "internal/localization/lang/id.json")

r.Use(middleware.Language(bundle))
```

**Behavior:**

1. Extracts Accept-Language header from request
2. Determines user's preferred language
3. Loads corresponding dictionary from bundle
4. Injects dictionary into context
5. Falls back to default language if not found

**Usage in Handler:**

```go
func (h *Handler) MyHandler(c *gin.Context) {
    dict := middleware.GetDict(c)
    // Use dict for translations
}
```

### Authentication

Validates JWT token and extracts user information.

**File:** `auth.go`

**Usage:**

```go
authMiddleware := middleware.NewAuthMiddleware(userUseCase)

protected := r.Group("/api/v1")
protected.Use(authMiddleware.Authenticate())
{
    protected.GET("/users/me", userHandler.GetProfile)
}
```

**Behavior:**

1. Extracts token from Authorization header
2. Validates Bearer token format
3. Verifies token signature and expiration
4. Sets userID in context
5. Returns 401 Unauthorized if validation fails

### CORS

Handles Cross-Origin Resource Sharing.

**File:** `cors.go`

**Usage:**

```go
r.Use(middleware.CORS())
```

**Behavior:**

1. Sets CORS headers for all responses
2. Handles preflight OPTIONS requests
3. Allows configured origins, methods, and headers

### Recovery

Recovers from panics and returns proper error response.

**File:** `recovery.go`

**Usage:**

```go
r.Use(middleware.Recovery())
```

**Behavior:**

1. Recovers from panics in handlers
2. Logs panic with stack trace
3. Returns 500 Internal Server Error response
4. Prevents server crash

## Middleware Order

Apply middleware in the correct order for proper execution:

```go
r := gin.New()

// 1. Logger - first to generate request_id for all logs
r.Use(middleware.Logger())

// 2. Language - extract Accept-Language for localized responses
r.Use(middleware.Language(bundle))

// 3. CORS - handle CORS before other processing
r.Use(middleware.CORS())

// 4. Recovery - catch panics from all previous middleware
r.Use(middleware.Recovery())

// 5. Authentication - apply to specific route groups
authMiddleware := middleware.NewAuthMiddleware(userUseCase)

// Public routes
r.POST("/api/v1/auth/register", userHandler.Register)
r.POST("/api/v1/auth/login", userHandler.Login)

// Protected routes
protected := r.Group("/api/v1")
protected.Use(authMiddleware.Authenticate())
{
    protected.GET("/users/me", userHandler.GetProfile)
    protected.GET("/users/:id", userHandler.GetByID)
}
```

**Order Rationale:**

1. **Logger first** - Generates request_id needed for all subsequent logs
2. **Language second** - Extract locale before any response generation
3. **CORS third** - Handle preflight before processing
4. **Recovery last** - Catch panics from all previous middleware

## Custom Middleware

Create custom middleware following the pattern:

```go
func CustomMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Before request processing
        
        c.Next()
        
        // After request processing
    }
}
```

## Best Practices

1. Keep middleware focused on single responsibility
2. Order matters - place middleware strategically
3. Use context to pass data between middleware and handlers
4. Log important events for debugging
5. Handle errors gracefully
6. Don't block request flow unnecessarily

