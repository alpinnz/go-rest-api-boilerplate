package repositories

import (
	"context"

	"github.com/alpinnz/go-rest-api-boilerplate/internal/domain/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRoleRepositoryImpl struct{}

func NewUserRoleRepository() *UserRoleRepositoryImpl {
	return &UserRoleRepositoryImpl{}
}

// AssignRole assigns a role to a user
func (r *UserRoleRepositoryImpl) AssignRole(ctx context.Context, tx *gorm.DB, userRole *entities.UserRole) error {
	if err := tx.WithContext(ctx).Create(userRole).Error; err != nil {
		return err
	}
	return nil
}

// RemoveRole removes a role from a user
func (r *UserRoleRepositoryImpl) RemoveRole(ctx context.Context, tx *gorm.DB, userID, roleID uuid.UUID) error {
	if err := tx.WithContext(ctx).Where("user_id = ? AND role_id = ?", userID, roleID).Delete(&entities.UserRole{}).Error; err != nil {
		return err
	}
	return nil
}

// GetRolesByUser retrieves all roles for a given user
func (r *UserRoleRepositoryImpl) GetRolesByUser(ctx context.Context, tx *gorm.DB, userID uuid.UUID) ([]*entities.Role, error) {
	var roles []*entities.Role
	if err := tx.WithContext(ctx).
		Joins("JOIN user_roles ur ON ur.role_id = roles.id").
		Where("ur.user_id = ?", userID).
		Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

// GetUsersByRole retrieves all users for a given role
func (r *UserRoleRepositoryImpl) GetUsersByRole(ctx context.Context, tx *gorm.DB, roleID uuid.UUID) ([]*entities.User, error) {
	var users []*entities.User
	if err := tx.WithContext(ctx).
		Joins("JOIN user_roles ur ON ur.user_id = users.id").
		Where("ur.role_id = ?", roleID).
		Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
