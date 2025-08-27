package repositories

import (
	"context"

	"github.com/alpinnz/go-rest-api-boilerplate/internal/domain/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	// Create inserts a new user into the database
	Create(ctx context.Context, tx *gorm.DB, user *entities.User) (*entities.User, error)

	// GetByID retrieves a user by their ID
	GetByID(ctx context.Context, tx *gorm.DB, id uuid.UUID, withRoles bool) (*entities.User, error)

	// GetByEmail retrieves a user by their email
	GetByEmail(ctx context.Context, tx *gorm.DB, email string) (*entities.User, error)

	// List retrieves a paginated list of users
	List(ctx context.Context, tx *gorm.DB, limit, offset int, withRoles bool) ([]*entities.User, error)

	// Count returns the total number of users
	Count(ctx context.Context, tx *gorm.DB) (int64, error)

	// Update modifies an existing user
	Update(ctx context.Context, tx *gorm.DB, user *entities.User) error

	// Delete removes a user by ID
	Delete(ctx context.Context, tx *gorm.DB, id uuid.UUID) error
}
