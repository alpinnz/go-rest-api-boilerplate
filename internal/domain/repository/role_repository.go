package repository

import (
	"context"

	"github.com/alpinnz/go-rest-api-boilerplate/internal/domain/entity"
	"github.com/google/uuid"
)

// RoleRepository defines data access methods for Role entities.
type RoleRepository interface {
	// Create inserts a new role into the database.
	// Returns ErrConflict if role name already exists.
	// Returns error if database operation fails.
	Create(ctx context.Context, role *entity.Role) error

	// FindByID retrieves a role by its unique identifier.
	// Returns ErrNotFound if role does not exist or is deleted.
	// Returns error if database operation fails.
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Role, error)

	// FindByName retrieves a role by its name.
	// Returns ErrNotFound if role does not exist or is deleted.
	// Returns error if database operation fails.
	FindByName(ctx context.Context, name string) (*entity.Role, error)

	// FindAll retrieves all roles with pagination.
	// Excludes soft-deleted roles by default.
	// Returns empty slice if no roles exist.
	// Returns error if database operation fails.
	FindAll(ctx context.Context, limit, offset int) ([]*entity.Role, error)

	// Update modifies an existing role's information.
	// Returns ErrNotFound if role does not exist.
	// Returns error if database operation fails.
	Update(ctx context.Context, role *entity.Role) error

	// Delete performs soft delete by setting DeletedAt timestamp.
	// Returns ErrNotFound if role does not exist.
	// Returns error if database operation fails.
	Delete(ctx context.Context, id uuid.UUID) error

	// HardDelete permanently removes a role from the database.
	// Use with caution - this cannot be undone.
	// Returns ErrNotFound if role does not exist.
	// Returns error if database operation fails.
	HardDelete(ctx context.Context, id uuid.UUID) error

	// AssignToUser assigns a role to a user.
	// Returns ErrConflict if assignment already exists.
	// Returns error if database operation fails.
	AssignToUser(ctx context.Context, userID, roleID uuid.UUID) error

	// RemoveFromUser removes a role assignment from a user.
	// Returns ErrNotFound if assignment does not exist.
	// Returns error if database operation fails.
	RemoveFromUser(ctx context.Context, userID, roleID uuid.UUID) error

	// FindByUserID retrieves all roles assigned to a user.
	// Excludes soft-deleted roles.
	// Returns empty slice if user has no roles.
	// Returns error if database operation fails.
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Role, error)

	// Count returns the total number of active (non-deleted) roles.
	// Returns error if database operation fails.
	Count(ctx context.Context) (int64, error)

	// Restore restores a soft-deleted role.
	// Returns ErrNotFound if role does not exist.
	// Returns error if database operation fails.
	Restore(ctx context.Context, id uuid.UUID) error
}
