package repositories

import (
	"context"

	"github.com/alpinnz/go-rest-api-boilerplate/internal/domain/entities"
	"github.com/alpinnz/go-rest-api-boilerplate/internal/domain/repositories"
	"github.com/alpinnz/go-rest-api-boilerplate/pkg/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserRepositoryImpl provides the concrete implementation for User repository operations.
// It communicates with the database using GORM.
type UserRepositoryImpl struct{}

// NewUserRepository creates and returns a new instance of UserRepositoryImpl.
func NewUserRepository() repositories.UserRepository {
	return &UserRepositoryImpl{}
}

// Create inserts a new user into the database.
//
// Parameters:
//   - ctx: context for cancellation and deadlines.
//   - tx:  active GORM transaction.
//   - user: pointer to the User entity to be created.
//
// Returns:
//   - The created User entity.
//   - Error if the operation fails.
func (r *UserRepositoryImpl) Create(ctx context.Context, tx *gorm.DB, user *entities.User) (*entities.User, error) {
	if err := tx.WithContext(ctx).Create(user).Error; err != nil {
		return errors.HandleRepoError(user, err)
	}
	return user, nil
}

// GetByID retrieves a user by its UUID.
//
// Parameters:
//   - ctx: context for cancellation and deadlines.
//   - tx:  active GORM transaction.
//   - id:  UUID of the user to retrieve.
//   - withRoles:
//
// Returns:
//   - The found User entity.
//   - Error if the user does not exist or a database error occurs.
func (r *UserRepositoryImpl) GetByID(ctx context.Context, tx *gorm.DB, id uuid.UUID, withRoles bool) (*entities.User, error) {
	var user *entities.User
	q := tx.WithContext(ctx)

	if withRoles {
		q = q.Preload("Roles")
	}

	if err := q.First(&user, "id = ?", id).Error; err != nil {
		return errors.HandleRepoError(user, err)
	}
	return user, nil
}

// GetByEmail retrieves a user by their email address.
//
// Parameters:
//   - ctx:   context for cancellation and deadlines.
//   - tx:    active GORM transaction.
//   - email: email of the user to retrieve.
//
// Returns:
//   - The found User entity.
//   - Error if the user does not exist or a database error occurs.
func (r *UserRepositoryImpl) GetByEmail(ctx context.Context, tx *gorm.DB, email string) (*entities.User, error) {
	var user *entities.User
	if err := tx.WithContext(ctx).First(&user, "email = ?", email).Error; err != nil {
		return errors.HandleRepoError(user, err)
	}
	return user, nil
}

// List retrieves users with pagination (limit & offset).
//
// Parameters:
//   - ctx:    context for cancellation and deadlines.
//   - tx:     active GORM transaction.
//   - limit:  maximum number of users to retrieve.
//   - offset: number of users to skip before starting to collect results.
//   - withRoles:
//
// Returns:
//   - A slice of User entities.
//   - Error if the query fails.
func (r *UserRepositoryImpl) List(ctx context.Context, tx *gorm.DB, limit, offset int, withRoles bool) ([]*entities.User, error) {
	var users []*entities.User
	q := tx.WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC")

	if withRoles {
		q = q.Preload("Roles")
	}

	if err := q.Find(&users).Error; err != nil {
		return errors.HandleRepoError(users, err)
	}
	return users, nil
}

// Count returns the total number of users (excluding soft-deleted ones).
//
// Parameters:
//   - ctx: context for cancellation and deadlines.
//   - tx:  active GORM transaction.
//
// Returns:
//   - The total number of users as int64.
//   - Error if the count query fails.
func (r *UserRepositoryImpl) Count(ctx context.Context, tx *gorm.DB) (int64, error) {
	var count int64
	if err := tx.WithContext(ctx).
		Model(&entities.User{}).
		Where("deleted_at IS NULL").
		Count(&count).Error; err != nil {
		return errors.HandleRepoError(count, err)
	}
	return count, nil
}

// Update modifies an existing user’s details.
//
// Parameters:
//   - ctx:  context for cancellation and deadlines.
//   - tx:   active GORM transaction.
//   - user: pointer to the User entity with updated values.
//
// Returns:
//   - Error if the update fails.
func (r *UserRepositoryImpl) Update(ctx context.Context, tx *gorm.DB, user *entities.User) error {
	updates := map[string]interface{}{
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"password":   user.Password,
		"updated_at": user.UpdatedAt,
	}

	if err := tx.WithContext(ctx).
		Model(&entities.User{}).
		Where("id = ?", user.ID).
		Updates(updates).Error; err != nil {
		_, appErr := errors.HandleRepoError(&entities.User{}, err)
		return appErr
	}
	return nil
}

// Delete performs a soft delete on a user by their UUID.
//
// Parameters:
//   - ctx: context for cancellation and deadlines.
//   - tx:  active GORM transaction.
//   - id:  UUID of the user to delete.
//
// Returns:
//   - Error if the delete fails or the user does not exist.
func (r *UserRepositoryImpl) Delete(ctx context.Context, tx *gorm.DB, id uuid.UUID) error {
	if err := tx.WithContext(ctx).Delete(&entities.User{}, "id = ?", id).Error; err != nil {
		_, appErr := errors.HandleRepoError(&entities.User{}, err)
		return appErr
	}
	return nil
}
