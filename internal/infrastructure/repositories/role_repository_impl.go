package repositories

import (
	"context"

	"github.com/alpinnz/go-rest-api-boilerplate/internal/domain/entities"
	"github.com/alpinnz/go-rest-api-boilerplate/pkg/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RoleRepositoryImpl provides the concrete implementation of the Role repository.
// It encapsulates database operations related to the Role entity using GORM.
type RoleRepositoryImpl struct{}

// NewRoleRepository creates and returns a new instance of RoleRepositoryImpl.
func NewRoleRepository() *RoleRepositoryImpl {
	return &RoleRepositoryImpl{}
}

// Create inserts a new role into the database.
//
// Parameters:
//   - ctx: context for cancellation and deadlines.
//   - tx:  active GORM transaction.
//   - role: pointer to the Role entity to be created.
//
// Returns:
//   - The created Role entity with ID populated.
//   - An error if the insert fails (e.g., duplicate constraint, DB error).
func (r *RoleRepositoryImpl) Create(ctx context.Context, tx *gorm.DB, role *entities.Role) (*entities.Role, error) {
	if err := tx.WithContext(ctx).Create(role).Error; err != nil {
		_, appErr := errors.HandleRepoError(*role, err)
		return role, appErr
	}
	return role, nil
}

// GetByID retrieves a role by its UUID primary key.
//
// Parameters:
//   - ctx: context for cancellation and deadlines.
//   - tx:  active GORM transaction.
//   - id:  UUID of the role to retrieve.
//
// Returns:
//   - Role entity if found.
//   - An error if not found or DB query fails.
func (r *RoleRepositoryImpl) GetByID(ctx context.Context, tx *gorm.DB, id uuid.UUID) (*entities.Role, error) {
	var role *entities.Role
	if err := tx.WithContext(ctx).First(&role, "id = ?", id).Error; err != nil {
		_, appErr := errors.HandleRepoError(role, err)
		return role, appErr
	}
	return role, nil
}

// GetByName retrieves a role by its unique name.
//
// Parameters:
//   - ctx: context for cancellation and deadlines.
//   - tx:  active GORM transaction.
//   - name: unique role name.
//
// Returns:
//   - Role entity if found.
//   - An error if not found or DB query fails.
func (r *RoleRepositoryImpl) GetByName(ctx context.Context, tx *gorm.DB, name string) (*entities.Role, error) {
	var role *entities.Role
	if err := tx.WithContext(ctx).First(&role, "name = ?", name).Error; err != nil {
		_, appErr := errors.HandleRepoError(role, err)
		return role, appErr
	}
	return role, nil
}

// List returns a paginated list of roles, ordered by creation time (newest first).
//
// Parameters:
//   - ctx: context for cancellation and deadlines.
//   - tx:  active GORM transaction.
//   - limit: number of records to fetch.
//   - offset: number of records to skip.
//
// Returns:
//   - Slice of Role entities.
//   - An error if the query fails.
func (r *RoleRepositoryImpl) List(ctx context.Context, tx *gorm.DB, limit, offset int) ([]*entities.Role, error) {
	var roles []*entities.Role
	if err := tx.WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&roles).Error; err != nil {
		_, appErr := errors.HandleRepoError(roles, err)
		return nil, appErr
	}
	return roles, nil
}

// Count returns the total number of roles (excluding soft-deleted ones).
//
// Parameters:
//   - ctx: context for cancellation and deadlines.
//   - tx:  active GORM transaction.
//
// Returns:
//   - Total count of roles.
//   - An error if the query fails.
func (r *RoleRepositoryImpl) Count(ctx context.Context, tx *gorm.DB) (int64, error) {
	var count int64
	if err := tx.WithContext(ctx).
		Model(&entities.Role{}).
		Where("deleted_at IS NULL").
		Count(&count).Error; err != nil {
		_, appErr := errors.HandleRepoError(count, err)
		return 0, appErr
	}
	return count, nil
}

// Update modifies an existing role's fields (only name and updated_at in this case).
//
// Parameters:
//   - ctx: context for cancellation and deadlines.
//   - tx:  active GORM transaction.
//   - role: pointer to the Role entity with updated values.
//
// Returns:
//   - Error if the update fails or the role does not exist.
func (r *RoleRepositoryImpl) Update(ctx context.Context, tx *gorm.DB, role *entities.Role) error {
	updates := map[string]interface{}{
		"name":       role.Name,
		"updated_at": role.UpdatedAt,
	}

	if err := tx.WithContext(ctx).
		Model(&entities.Role{}).
		Where("id = ?", role.ID).
		Updates(updates).Error; err != nil {
		_, appErr := errors.HandleRepoError(&entities.Role{}, err)
		return appErr
	}
	return nil
}

// Delete performs a soft delete on a role by its UUID.
// The record is not removed permanently, but marked as deleted (deleted_at set).
//
// Parameters:
//   - ctx: context for cancellation and deadlines.
//   - tx:  active GORM transaction.
//   - id:  UUID of the role to delete.
//
// Returns:
//   - Error if the delete fails or the role does not exist.
func (r *RoleRepositoryImpl) Delete(ctx context.Context, tx *gorm.DB, id uuid.UUID) error {
	if err := tx.WithContext(ctx).Delete(&entities.Role{}, "id = ?", id).Error; err != nil {
		_, appErr := errors.HandleRepoError(&entities.Role{}, err)
		return appErr
	}
	return nil
}
