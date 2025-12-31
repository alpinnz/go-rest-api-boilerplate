# Auth Package

Authentication utilities for JWT token management and password hashing.

## Features

- JWT token generation and validation
- Password hashing with bcrypt
- Secure password comparison

## JWT Token Management

### Setup

Configure JWT secret during application initialization:

```go
import "github.com/alpinnz/go-rest-api-boilerplate/pkg/auth"

func main() {
    auth.SetJWTSecret(os.Getenv("JWT_SECRET"))
}
```

### Token Types

**Access Token** (short-lived):
- Lifetime: 15 minutes
- Used for API requests
- Sent in Authorization header

**Refresh Token** (long-lived):
- Lifetime: 7 days
- Used to obtain new access tokens
- Single-use (rotation pattern)

### Generate Tokens

**Access Token:**
```go
accessToken, err := auth.GenerateAccessToken(userID)
if err != nil {
    return err
}
```

**Refresh Token:**
```go
refreshToken, err := auth.GenerateRefreshToken(userID)
if err != nil {
    return err
}
```

**Legacy (backward compatible):**
```go
token, err := auth.GenerateToken(userID) // Returns access token
```

Token specifications:
- Algorithm: HMAC-SHA256
- Claims: `user_id`, `token_type`, `exp`, `iat`
- Token types: `"access"` or `"refresh"`

### Validate Token

Verify and extract claims from token:

```go
claims, err := auth.ValidateToken(tokenString)
if err != nil {
    return errors.New("invalid or expired token")
}

userID := claims.UserID
tokenType := claims.TokenType // "access" or "refresh"
```

Validation checks:
- Signing method verification
- Signature validation
- Expiration time
- Token type verification

### Token Usage

**API Request:**
```
Authorization: Bearer <access_token>
```

**Refresh Flow:**
```
POST /api/v1/auth/refresh
{
  "refresh_token": "<refresh_token>"
}
```

## Password Management

### Hash Password

Hash plaintext password before storage:

```go
hashedPassword, err := auth.HashPassword("user_password")
if err != nil {
    return err
}
// Store hashedPassword in database
```

Specifications:
- Algorithm: bcrypt
- Cost factor: 14
- Processing time: ~1 second on modern hardware

### Verify Password

Check password against stored hash:

```go
isValid := auth.CheckPasswordHash("user_password", storedHash)
if !isValid {
    return errors.New("invalid credentials")
}
```

Security features:
- Constant-time comparison
- Timing attack prevention

## Security Best Practices

1. **JWT Secret**: Load from environment variables, never hardcode
2. **Password Cost**: Use cost factor 14 or higher for production
3. **Token Storage**: Store tokens securely on client side
4. **Token Transmission**: Always use HTTPS in production
5. **Password Requirements**: Enforce minimum length and complexity

## Example Integration

```go
package main

import (
    "os"
    "github.com/alpinnz/go-rest-api-boilerplate/pkg/auth"
)

func init() {
    auth.SetJWTSecret(os.Getenv("JWT_SECRET"))
}

func register(email, password string) error {
    hashedPassword, err := auth.HashPassword(password)
    if err != nil {
        return err
    }
    
    userID := saveUser(email, hashedPassword)
    token, err := auth.GenerateToken(userID)
    if err != nil {
        return err
    }
    
    return sendToken(token)
}

func login(email, password string) error {
    user := findUserByEmail(email)
    
    isValid := auth.CheckPasswordHash(password, user.Password)
    if !isValid {
        return errors.New("invalid credentials")
    }
    
    token, err := auth.GenerateToken(user.ID)
    if err != nil {
        return err
    }
    
    return sendToken(token)
}
```

