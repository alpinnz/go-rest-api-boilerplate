// Package auth provides authentication utilities including JWT token generation and validation.
package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// jwtSecret holds the secret key used for signing JWT tokens.
// Should be set via SetJWTSecret during application initialization.
// Default value is placeholder and MUST be changed in production.
var jwtSecret = []byte("change-this-secret-key")

// Claims represents JWT token claims including user identity.
// Embeds jwt.RegisteredClaims for standard claims (exp, iat, etc.).
type Claims struct {
	UserID int64 `json:"user_id"` // Authenticated user's unique identifier
	jwt.RegisteredClaims
}

// SetJWTSecret configures the secret key used for signing JWT tokens.
// Must be called during application initialization before generating tokens.
// Secret should be loaded from environment variables, never hardcoded.
//
// Example:
//
//	auth.SetJWTSecret(os.Getenv("JWT_SECRET"))
func SetJWTSecret(secret string) {
	jwtSecret = []byte(secret)
}

// GenerateToken creates a new JWT token for authenticated user.
// Token includes user ID and expires after 24 hours from creation.
// Uses HMAC-SHA256 algorithm for signing.
// Returns signed token string or error if signing fails.
//
// Example:
//
//	token, err := auth.GenerateToken(12345)
//	if err != nil {
//	    return err
//	}
//	// Use token in Authorization header: Bearer <token>
func GenerateToken(userID int64) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ValidateToken verifies JWT token signature and extracts claims.
// Validates signing method, token signature, and expiration.
// Returns parsed claims or error if validation fails.
//
// Possible errors:
//   - "invalid signing method" - Token uses unexpected algorithm
//   - "invalid token" - Signature verification failed
//   - jwt.ErrTokenExpired - Token has expired
//
// Example:
//
//	claims, err := auth.ValidateToken(tokenString)
//	if err != nil {
//	    return errors.New("invalid or expired token")
//	}
//	userID := claims.UserID
func ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
