// Package auth provides authentication utilities including JWT token generation and validation.
package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Token types
const (
	TokenTypeAccess  = "access"
	TokenTypeRefresh = "refresh"
)

// JWT secrets and expiration durations
// Should be set via SetJWTConfig during application initialization.
var (
	accessTokenSecret      = []byte("change-this-access-token-secret-key")
	accessTokenExpiration  = 15 * time.Minute // 15 minutes
	refreshTokenSecret     = []byte("change-this-refresh-token-secret-key")
	refreshTokenExpiration = 7 * 24 * time.Hour // 7 days
)

// Claims represents JWT token claims including user identity and token type.
// Embeds jwt.RegisteredClaims for standard claims (exp, iat, etc.).
type Claims struct {
	UserID    string `json:"user_id"`    // Authenticated user's unique identifier (UUID string)
	TokenType string `json:"token_type"` // Token type: "access" or "refresh"
	jwt.RegisteredClaims
}

// SetJWTConfig configures the secrets and expirations for JWT tokens.
// Must be called during application initialization before generating tokens.
// Secrets should be loaded from environment variables, never hardcoded.
//
// Example:
//
//	auth.SetJWTConfig(
//	    cfg.JWT.AccessTokenSecret, cfg.JWT.AccessTokenExpiration,
//	    cfg.JWT.RefreshTokenSecret, cfg.JWT.RefreshTokenExpiration,
//	)
func SetJWTConfig(accessSecret string, accessExp time.Duration, refreshSecret string, refreshExp time.Duration) {
	accessTokenSecret = []byte(accessSecret)
	accessTokenExpiration = accessExp
	refreshTokenSecret = []byte(refreshSecret)
	refreshTokenExpiration = refreshExp
}

// GetAccessTokenExpiration returns the configured access token expiration duration.
func GetAccessTokenExpiration() time.Duration {
	return accessTokenExpiration
}

// GetRefreshTokenExpiration returns the configured refresh token expiration duration.
func GetRefreshTokenExpiration() time.Duration {
	return refreshTokenExpiration
}

// GenerateToken creates a new JWT access token for authenticated user.
// Token expires based on configured ACCESS_TOKEN_EXPIRATION.
// Uses HMAC-SHA256 algorithm for signing.
// Returns signed token string or error if signing fails.
//
// Example:
//
//	token, err := auth.GenerateToken(userID)
//	if err != nil {
//	    return err
//	}
//	// Use token in Authorization header: Bearer <token>
func GenerateToken(userID uuid.UUID) (string, error) {
	return generateTokenWithType(userID, TokenTypeAccess, accessTokenExpiration, accessTokenSecret)
}

// GenerateAccessToken creates a new JWT access token (short-lived).
// Token expires based on configured ACCESS_TOKEN_EXPIRATION.
func GenerateAccessToken(userID uuid.UUID) (string, error) {
	return generateTokenWithType(userID, TokenTypeAccess, accessTokenExpiration, accessTokenSecret)
}

// GenerateRefreshToken creates a new JWT refresh token (long-lived).
// Token expires based on configured REFRESH_TOKEN_EXPIRATION.
func GenerateRefreshToken(userID uuid.UUID) (string, error) {
	return generateTokenWithType(userID, TokenTypeRefresh, refreshTokenExpiration, refreshTokenSecret)
}

// generateTokenWithType creates a JWT token with specified type, expiration, and secret.
func generateTokenWithType(userID uuid.UUID, tokenType string, expiration time.Duration, secret []byte) (string, error) {
	claims := Claims{
		UserID:    userID.String(),
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

// ValidateToken verifies JWT token signature and extracts claims.
// Validates signing method, token signature, and expiration.
// Uses appropriate secret (access or refresh) based on token type.
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
	// First parse without validation to get token type
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		// Get claims to determine token type
		claims, ok := token.Claims.(*Claims)
		if !ok {
			return nil, errors.New("invalid token claims")
		}

		// Use appropriate secret based on token type
		if claims.TokenType == TokenTypeRefresh {
			return refreshTokenSecret, nil
		}
		return accessTokenSecret, nil
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
