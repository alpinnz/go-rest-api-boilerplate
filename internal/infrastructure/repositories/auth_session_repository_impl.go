package repositories

import (
	"context"

	"github.com/alpinnz/go-rest-api-boilerplate/internal/domain/entities"
	"github.com/alpinnz/go-rest-api-boilerplate/pkg/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AuthSessionRepositoryImpl provides the concrete implementation of the AuthSession repository.
// It encapsulates database operations related to the AuthSession entity using GORM.
type AuthSessionRepositoryImpl struct{}

// NewAuthSessionRepository creates and returns a new instance of AuthSessionRepositoryImpl.
func NewAuthSessionRepository() *AuthSessionRepositoryImpl {
	return &AuthSessionRepositoryImpl{}
}

// Create inserts a new authentication session into the database.
//
// Parameters:
//   - ctx: context for cancellation and deadlines.
//   - tx:  active GORM transaction.
//   - authSession: pointer to the AuthSession entity to be created.
//
// Returns:
//   - The created AuthSession entity with ID populated.
//   - An error if the insert fails (e.g., constraint violation, DB error).
func (r *AuthSessionRepositoryImpl) Create(ctx context.Context, tx *gorm.DB, authSession *entities.AuthSession) (*entities.AuthSession, error) {
	if err := tx.WithContext(ctx).Create(authSession).Error; err != nil {
		_, appErr := errors.HandleRepoError(*authSession, err)
		return authSession, appErr
	}
	return authSession, nil
}

// GetByAccessToken retrieves an authentication session by its access token.
//
// Parameters:
//   - ctx: context for cancellation and deadlines.
//   - tx:  active GORM transaction.
//   - token: access token string.
//
// Returns:
//   - AuthSession entity if found.
//   - An error if not found or DB query fails.
func (r *AuthSessionRepositoryImpl) GetByAccessToken(ctx context.Context, tx *gorm.DB, token string) (*entities.AuthSession, error) {
	var authSession *entities.AuthSession
	if err := tx.WithContext(ctx).First(&authSession, "access_token = ?", token).Error; err != nil {
		_, appErr := errors.HandleRepoError(authSession, err)
		return authSession, appErr
	}
	return authSession, nil
}

// GetByRefreshToken retrieves an authentication session by its refresh token.
//
// Parameters:
//   - ctx: context for cancellation and deadlines.
//   - tx:  active GORM transaction.
//   - token: refresh token string.
//
// Returns:
//   - AuthSession entity if found.
//   - An error if not found or DB query fails.
func (r *AuthSessionRepositoryImpl) GetByRefreshToken(ctx context.Context, tx *gorm.DB, token string) (*entities.AuthSession, error) {
	var authSession *entities.AuthSession
	if err := tx.WithContext(ctx).First(&authSession, "refresh_token = ?", token).Error; err != nil {
		_, appErr := errors.HandleRepoError(authSession, err)
		return authSession, appErr
	}
	return authSession, nil
}

// Delete performs a soft delete of an authentication session by its UUID.
// The record is not removed permanently, but marked as deleted (deleted_at set).
//
// Parameters:
//   - ctx: context for cancellation and deadlines.
//   - tx:  active GORM transaction.
//   - id:  UUID of the session to delete.
//
// Returns:
//   - Error if the delete fails or the session does not exist.
func (r *AuthSessionRepositoryImpl) Delete(ctx context.Context, tx *gorm.DB, id uuid.UUID) error {
	if err := tx.WithContext(ctx).Delete(&entities.AuthSession{}, "id = ?", id).Error; err != nil {
		_, appErr := errors.HandleRepoError(&entities.AuthSession{}, err)
		return appErr
	}
	return nil
}
