package repository

import (
	"context"

	"github.com/alpinnz/go-rest-api-boilerplate/internal/domain/entity"
	"github.com/google/uuid"
)

// UserRepository defines data access methods for User entities.
// All methods accept context for cancellation and timeout control.
// Implementations must handle errors explicitly and return domain errors.
type UserRepository interface {
	// Create inserts a new user into the database.
	// Returns ErrConflict if email already exists.
	// Returns error if database operation fails.
	Create(ctx context.Context, user *entity.User) error

	// FindByID retrieves a user by their unique identifier.
	// Returns ErrNotFound if user does not exist.
	// Returns error if database operation fails.
	FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error)

	// FindByEmail retrieves a user by email address.
	// Email lookup is case-insensitive for better UX.
	// Returns ErrNotFound if user does not exist.
	// Returns error if database operation fails.
	FindByEmail(ctx context.Context, email string) (*entity.User, error)

	// FindAll retrieves all users with pagination.
	// Excludes soft-deleted users by default.
	// Returns empty slice if no users exist.
	// Returns error if database operation fails.
	FindAll(ctx context.Context, limit, offset int) ([]*entity.User, error)

	// Update modifies an existing user's information.
	// Only updates non-zero fields to support partial updates.
	// Returns ErrNotFound if user does not exist.
	// Returns error if database operation fails.
	Update(ctx context.Context, user *entity.User) error

	// Delete performs soft delete by setting DeletedAt timestamp.
	// Returns ErrNotFound if user does not exist.
	// Returns error if database operation fails.
	Delete(ctx context.Context, id uuid.UUID) error

	// HardDelete permanently removes a user from the database.
	// Use with caution - this cannot be undone.
	// Returns ErrNotFound if user does not exist.
	// Returns error if database operation fails.
	HardDelete(ctx context.Context, id uuid.UUID) error

	// Count returns the total number of active (non-deleted) users.
	// Returns error if database operation fails.
	Count(ctx context.Context) (int64, error)

	// Restore restores a soft-deleted user.
	// Returns ErrNotFound if user does not exist.
	// Returns error if database operation fails.
	Restore(ctx context.Context, id uuid.UUID) error
}
