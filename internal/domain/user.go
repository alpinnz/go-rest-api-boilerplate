package domain

import (
	"context"
	"time"
)

// User represents an authenticated user entity in the system.
// Password field is never exposed in JSON responses (json:"-" tag).
type User struct {
	ID        int64      `json:"id"`
	Email     string     `json:"email"`
	Password  string     `json:"-"` // bcrypt hashed password
	Name      string     `json:"name"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"` // soft delete timestamp
}

// UserRepository defines data access methods for User entities.
// All methods accept context for cancellation and timeout control.
// Implementations must handle errors explicitly and return domain errors.
type UserRepository interface {
	// Create inserts a new user into the database.
	// Returns ErrConflict if email already exists.
	// Returns error if database operation fails.
	Create(ctx context.Context, user *User) error

	// FindByID retrieves a user by their unique identifier.
	// Returns ErrNotFound if user does not exist.
	// Returns error if database operation fails.
	FindByID(ctx context.Context, id int64) (*User, error)

	// FindByEmail retrieves a user by email address.
	// Email lookup is case-insensitive for better UX.
	// Returns ErrNotFound if user does not exist.
	// Returns error if database operation fails.
	FindByEmail(ctx context.Context, email string) (*User, error)

	// Update modifies an existing user's information.
	// Only updates non-zero fields to support partial updates.
	// Returns ErrNotFound if user does not exist.
	// Returns error if database operation fails.
	Update(ctx context.Context, user *User) error

	// Delete removes a user from the database (hard delete).
	// For soft delete, use Update with DeletedAt timestamp.
	// Returns ErrNotFound if user does not exist.
	// Returns error if database operation fails.
	Delete(ctx context.Context, id int64) error
}

// Session represents an active user session stored in Redis.
// Sessions are used for fast authentication token validation.
type Session struct {
	Token  string // JWT token used as Redis key
	UserID int64  // Associated user ID
}

// SessionRepository defines data access methods for Session storage.
// Implementations use Redis for fast in-memory access with TTL.
// All methods accept context for cancellation and timeout control.
type SessionRepository interface {
	// Set stores a new session in Redis with expiration.
	// Token is used as key, userID as value.
	// Expiration determines TTL (time-to-live) in Redis.
	// Returns error if Redis operation fails.
	Set(ctx context.Context, token string, userID int64, expiration time.Duration) error

	// Get retrieves userID for a given token from Redis.
	// Returns ErrSessionExpired if token doesn't exist (expired or invalid).
	// Returns error if Redis operation fails.
	Get(ctx context.Context, token string) (int64, error)

	// Delete removes a session from Redis (logout).
	// Used for immediate session invalidation.
	// Returns error if Redis operation fails.
	Delete(ctx context.Context, token string) error
}
