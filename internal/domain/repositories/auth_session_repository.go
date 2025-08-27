package repositories

import (
	"context"

	"github.com/alpinnz/go-rest-api-boilerplate/internal/domain/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AuthSessionRepository handles persistence of authentication sessions
type AuthSessionRepository interface {
	// Create inserts a new auth session into the database
	Create(ctx context.Context, tx *gorm.DB, session *entities.AuthSession) (*entities.AuthSession, error)

	// GetByAccessToken retrieves a session by access token
	GetByAccessToken(ctx context.Context, tx *gorm.DB, token string) (*entities.AuthSession, error)

	// GetByRefreshToken retrieves a session by refresh token
	GetByRefreshToken(ctx context.Context, tx *gorm.DB, token string) (*entities.AuthSession, error)

	// Delete removes a session by ID (used for logout/cleanup)
	Delete(ctx context.Context, tx *gorm.DB, id uuid.UUID) error
}
