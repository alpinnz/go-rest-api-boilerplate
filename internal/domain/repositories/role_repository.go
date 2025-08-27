package repositories

import (
	"context"

	"github.com/alpinnz/go-rest-api-boilerplate/internal/domain/entities"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoleRepository interface {
	// Create inserts a new role into the database
	Create(ctx context.Context, tx *gorm.DB, role *entities.Role) (*entities.Role, error)

	// GetByID retrieves a role by their ID
	GetByID(ctx context.Context, tx *gorm.DB, id uuid.UUID) (*entities.Role, error)

	// GetByName retrieves a role by their name
	GetByName(ctx context.Context, tx *gorm.DB, name string) (*entities.Role, error)

	// List retrieves a paginated list of roles
	List(ctx context.Context, tx *gorm.DB, limit, offset int) ([]*entities.Role, error)

	// Count returns the total number of roles
	Count(ctx context.Context, tx *gorm.DB) (int64, error)

	// Update modifies an existing role
	Update(ctx context.Context, tx *gorm.DB, role *entities.Role) error

	// Delete removes a role by ID
	Delete(ctx context.Context, tx *gorm.DB, id uuid.UUID) error
}
