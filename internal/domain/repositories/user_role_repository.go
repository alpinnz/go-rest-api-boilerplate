package repositories

import (
	"context"

	"github.com/alpinnz/go-rest-api-boilerplate/internal/domain/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserRoleRepository handles mapping between users and roles
type UserRoleRepository interface {
	// AssignRole assigns a role to a user
	AssignRole(ctx context.Context, tx *gorm.DB, userRole *entities.UserRole) error

	// RemoveRole removes a role from a user
	RemoveRole(ctx context.Context, tx *gorm.DB, userID, roleID uuid.UUID) error

	// GetRolesByUser retrieves all roles for a given user
	GetRolesByUser(ctx context.Context, tx *gorm.DB, userID uuid.UUID) ([]*entities.Role, error)

	// GetUsersByRole retrieves all users for a given role
	GetUsersByRole(ctx context.Context, tx *gorm.DB, roleID uuid.UUID) ([]*entities.User, error)
}
