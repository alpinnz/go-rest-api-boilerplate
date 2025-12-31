package domain

import (
	"context"
	"time"
)

// SessionRepository defines data access methods for Session storage.
// Implementations use Redis for fast in-memory access with TTL.
// All methods accept context for cancellation and timeout control.
type SessionRepository interface {
	// SetAccessToken stores a new session in Redis with expiration.
	// Token is used as key, userID as value.
	// Expiration determines TTL (time-to-live) in Redis.
	// Returns error if Redis operation fails.
	SetAccessToken(ctx context.Context, token string, userID int64, expiration time.Duration) error

	// GetAccessToken retrieves userID for a given token from Redis.
	// Returns ErrSessionExpired if token doesn't exist (expired or invalid).
	// Returns error if Redis operation fails.
	GetAccessToken(ctx context.Context, token string) (int64, error)

	// DeleteAccessToken removes a session from Redis (logout).
	// Used for immediate session invalidation.
	// Returns error if Redis operation fails.
	DeleteAccessToken(ctx context.Context, token string) error

	// SetRefreshToken stores refresh token mapping to access token in Redis.
	// Used to validate and track refresh token usage.
	// Returns error if Redis operation fails.
	SetRefreshToken(ctx context.Context, refreshToken string, accessToken string, expiration time.Duration) error

	// GetByRefreshToken retrieves access token for a given refresh token.
	// Returns ErrSessionExpired if refresh token doesn't exist or expired.
	// Returns error if Redis operation fails.
	GetByRefreshToken(ctx context.Context, refreshToken string) (string, error)

	// DeleteRefreshToken removes refresh token from Redis.
	// Used when refresh token is used or during logout.
	// Returns error if Redis operation fails.
	DeleteRefreshToken(ctx context.Context, refreshToken string) error
}
