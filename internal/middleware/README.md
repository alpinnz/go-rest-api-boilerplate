# Middleware Package

HTTP middleware for request processing and cross-cutting concerns.

## Overview

Middleware package provides reusable HTTP middleware for Gin framework. Each middleware handles specific cross-cutting concerns like authentication, CORS, logging, and error recovery.

## Available Middleware

### Authentication

Validates JWT token and extracts user information.

**File:** `auth.go`

**Usage:**

```go
authMiddleware := middleware.NewAuthMiddleware(userUseCase)

// Protected routes
protected := r.Group("/api/v1")
protected.Use(authMiddleware.Authenticate())
{
    protected.GET("/profile", userHandler.GetProfile)
}
```

**Behavior:**

1. Extracts token from Authorization header
2. Validates Bearer token format
3. Verifies token signature and expiration
4. Sets `userID` in context
5. Returns 401 Unauthorized if validation fails

**Response on Error:**

```json
{
  "code": "UNAUTHORIZED",
  "message": "Unauthorized",
  "data": null,
  "errors": {
    "reason": "Invalid or expired token"
  },
  "meta": null,
  "trace_id": "abc123"
}
```

**Example:**

```go
type AuthMiddleware struct {
    userUseCase *usecase.UserUseCase
}

func NewAuthMiddleware(userUseCase *usecase.UserUseCase) *AuthMiddleware {
    return &AuthMiddleware{userUseCase: userUseCase}
}

func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            response.Unauthorized(c, "Authorization header required")
            c.Abort()
            return
        }

        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            response.Unauthorized(c, "Invalid authorization format")
            c.Abort()
            return
        }

        token := parts[1]
        userID, err := m.userUseCase.ValidateSession(c.Request.Context(), token)
        if err != nil {
            response.Unauthorized(c, "Invalid or expired token")
            c.Abort()
            return
        }

        c.Set("userID", userID)
        c.Next()
    }
}
```

### CORS

Handles Cross-Origin Resource Sharing headers.

**File:** `cors.go`

**Usage:**

```go
r := gin.New()
r.Use(middleware.CORS())
```

**Configuration:**

- **Origin:** `*` (allow all origins)
- **Credentials:** `true`
- **Methods:** GET, POST, PUT, DELETE, PATCH, OPTIONS
- **Headers:** Content-Type, Authorization, X-CSRF-Token, etc.

**Behavior:**

1. Sets CORS headers on all responses
2. Handles OPTIONS preflight requests
3. Allows credentials in requests

**Example:**

```go
func CORS() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }

        c.Next()
    }
}
```

**Production Configuration:**

For production, restrict allowed origins:

```go
func CORS() gin.HandlerFunc {
    return func(c *gin.Context) {
        origin := c.Request.Header.Get("Origin")
        allowedOrigins := []string{
            "https://example.com",
            "https://www.example.com",
        }
        
        for _, allowed := range allowedOrigins {
            if origin == allowed {
                c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
                break
            }
        }
        
        // ...rest of headers
    }
}
```

### Logger

Logs HTTP request details and response time.

**File:** `logger.go`

**Usage:**

```go
r := gin.New()
r.Use(middleware.Logger())
```

**Logged Information:**

- HTTP method
- Request path
- Status code
- Request latency

**Output Format:**

```
GET /api/v1/users 200 15.234ms
POST /api/v1/login 401 8.123ms
```

**Example:**

```go
func Logger() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        path := c.Request.URL.Path
        method := c.Request.Method

        c.Next()

        latency := time.Since(start)
        statusCode := c.Writer.Status()

        log.Printf("%s %s %d %v", method, path, statusCode, latency)
    }
}
```

**Enhanced Logging:**

```go
func Logger() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()

        c.Next()

        log.Printf("[%s] %s %s %d %v %s",
            c.Request.Method,
            c.Request.URL.Path,
            c.ClientIP(),
            c.Writer.Status(),
            time.Since(start),
            c.Errors.String(),
        )
    }
}
```

### Recovery

Recovers from panics and returns consistent error response.

**File:** `recovery.go`

**Usage:**

```go
r := gin.New()
r.Use(middleware.Recovery())
```

**Behavior:**

1. Catches panics in request handlers
2. Logs panic details
3. Returns 500 Internal Server Error
4. Prevents application crash

**Response on Panic:**

```json
{
  "code": "INTERNAL_SERVER_ERROR",
  "message": "Unexpected error occurred",
  "data": null,
  "errors": {
    "hint": "Unexpected error occurred"
  },
  "meta": null,
  "trace_id": "abc123"
}
```

**Example:**

```go
func Recovery() gin.HandlerFunc {
    return func(c *gin.Context) {
        defer func() {
            if err := recover(); err != nil {
                log.Printf("panic recovered: %v", err)
                response.InternalServerError(c, "Unexpected error occurred")
                c.Abort()
            }
        }()
        c.Next()
    }
}
```

## Middleware Order

Apply middleware in the correct order for proper execution:

```go
r := gin.New()

// 1. Recovery - must be first to catch all panics
r.Use(middleware.Recovery())

// 2. CORS - handle CORS before other processing
r.Use(middleware.CORS())

// 3. Logger - log all requests
r.Use(middleware.Logger())

// 4. Authentication - apply to specific route groups
authMiddleware := middleware.NewAuthMiddleware(userUseCase)

// Public routes
r.POST("/api/v1/register", userHandler.Register)
r.POST("/api/v1/login", userHandler.Login)

// Protected routes
protected := r.Group("/api/v1")
protected.Use(authMiddleware.Authenticate())
{
    protected.GET("/profile", userHandler.GetProfile)
    protected.GET("/users/:id", userHandler.GetByID)
}
```

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

### Example: Request ID

```go
func RequestID() gin.HandlerFunc {
    return func(c *gin.Context) {
        requestID := c.GetHeader("X-Request-ID")
        if requestID == "" {
            requestID = uuid.New().String()
        }
        
        c.Set("request_id", requestID)
        c.Writer.Header().Set("X-Request-ID", requestID)
        
        c.Next()
    }
}
```

### Example: Rate Limiting

```go
func RateLimit(limit int) gin.HandlerFunc {
    limiter := rate.NewLimiter(rate.Limit(limit), limit)
    
    return func(c *gin.Context) {
        if !limiter.Allow() {
            response.ServiceUnavailable(c, 60)
            c.Abort()
            return
        }
        
        c.Next()
    }
}
```

## Testing Middleware

```go
func TestAuthMiddleware(t *testing.T) {
    gin.SetMode(gin.TestMode)
    
    r := gin.New()
    authMiddleware := middleware.NewAuthMiddleware(mockUserUseCase)
    r.Use(authMiddleware.Authenticate())
    r.GET("/test", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok"})
    })
    
    // Test without token
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/test", nil)
    r.ServeHTTP(w, req)
    
    assert.Equal(t, 401, w.Code)
    
    // Test with valid token
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("GET", "/test", nil)
    req.Header.Set("Authorization", "Bearer valid-token")
    r.ServeHTTP(w, req)
    
    assert.Equal(t, 200, w.Code)
}
```

