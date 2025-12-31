package domain

import (
	"context"
	"time"
)

// User represents an authenticated user entity in the system.
// This is a pure domain entity - no JSON tags (serialization is DTO concern).
type User struct {
	ID        int64
	Email     string
	Password  string // bcrypt hashed password
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time // soft delete timestamp
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
