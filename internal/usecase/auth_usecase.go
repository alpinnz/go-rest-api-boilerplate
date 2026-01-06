package usecase

import (
	"context"
	"time"

	"github.com/alpinnz/go-rest-api-boilerplate/internal/domain"
	"github.com/alpinnz/go-rest-api-boilerplate/internal/domain/entity"
	"github.com/alpinnz/go-rest-api-boilerplate/internal/domain/repository"
	"github.com/alpinnz/go-rest-api-boilerplate/pkg/auth"
	"github.com/google/uuid"
)

type AuthUseCase struct {
	userRepo       repository.UserRepository
	sessionRepo    repository.SessionRepository
	contextTimeout time.Duration
}

func NewAuthUseCase(
	userRepo repository.UserRepository,
	sessionRepo repository.SessionRepository,
	timeout time.Duration,
) *AuthUseCase {
	return &AuthUseCase{
		userRepo:       userRepo,
		sessionRepo:    sessionRepo,
		contextTimeout: timeout,
	}
}

// RegisterInput represents plain input for registration (no HTTP concerns)
type RegisterInput struct {
	Email     string
	Password  string
	FirstName string
	LastName  string
}

// LoginInput represents plain input for login (no HTTP concerns)
type LoginInput struct {
	Email    string
	Password string
}

// AuthResult represents authentication result (no HTTP concerns)
type AuthResult struct {
	AccessToken  string
	RefreshToken string
	User         *entity.User
}

func (uc *AuthUseCase) Register(ctx context.Context, input RegisterInput) (*entity.User, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()

	if input.Email == "" || input.Password == "" || input.FirstName == "" || input.LastName == "" {
		return nil, domain.ErrInvalidInput
	}

	existing, _ := uc.userRepo.FindByEmail(ctx, input.Email)
	if existing != nil {
		return nil, domain.ErrConflict
	}

	hashedPassword, err := auth.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	user := &entity.User{
		Email:     input.Email,
		Password:  hashedPassword,
		FirstName: input.FirstName,
		LastName:  input.LastName,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (uc *AuthUseCase) Login(ctx context.Context, input LoginInput) (*AuthResult, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()

	if input.Email == "" || input.Password == "" {
		return nil, domain.ErrInvalidInput
	}

	user, err := uc.userRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	if !auth.CheckPasswordHash(input.Password, user.Password) {
		return nil, domain.ErrInvalidCredentials
	}

	// Generate access token (short-lived)
	accessToken, err := auth.GenerateAccessToken(user.ID)
	if err != nil {
		return nil, err
	}

	// Generate refresh token (long-lived)
	refreshToken, err := auth.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	// Store access token session
	if err := uc.sessionRepo.SetAccessToken(ctx, accessToken, user.ID, auth.GetAccessTokenExpiration()); err != nil {
		return nil, err
	}

	// Store refresh token mapping
	if err := uc.sessionRepo.SetRefreshToken(ctx, refreshToken, accessToken, auth.GetRefreshTokenExpiration()); err != nil {
		return nil, err
	}

	return &AuthResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}

func (uc *AuthUseCase) Logout(ctx context.Context, token string) error {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()

	return uc.sessionRepo.DeleteAccessToken(ctx, token)
}

// RefreshToken generates new access token from valid refresh token.
// Returns new access token and updated refresh token.
func (uc *AuthUseCase) RefreshToken(ctx context.Context, refreshToken string) (*AuthResult, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()

	// Validate refresh token
	claims, err := auth.ValidateToken(refreshToken)
	if err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	// Check token type
	if claims.TokenType != auth.TokenTypeRefresh {
		return nil, domain.ErrInvalidCredentials
	}

	// Check if refresh token exists in Redis
	oldAccessToken, err := uc.sessionRepo.GetByRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, domain.ErrSessionExpired
	}

	// Delete old access token session
	_ = uc.sessionRepo.DeleteAccessToken(ctx, oldAccessToken)

	// Delete used refresh token
	_ = uc.sessionRepo.DeleteRefreshToken(ctx, refreshToken)

	// Get user to generate new tokens
	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	// Generate new access token
	newAccessToken, err := auth.GenerateAccessToken(user.ID)
	if err != nil {
		return nil, err
	}

	// Generate new refresh token
	newRefreshToken, err := auth.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	// Store new access token session
	if err := uc.sessionRepo.SetAccessToken(ctx, newAccessToken, user.ID, auth.GetAccessTokenExpiration()); err != nil {
		return nil, err
	}

	// Store new refresh token mapping
	if err := uc.sessionRepo.SetRefreshToken(ctx, newRefreshToken, newAccessToken, auth.GetRefreshTokenExpiration()); err != nil {
		return nil, err
	}

	return &AuthResult{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		User:         user,
	}, nil
}

// ValidateSession validates access token and returns user ID.
// This is used by the authentication middleware.
func (uc *AuthUseCase) ValidateSession(ctx context.Context, token string) (uuid.UUID, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()

	userID, err := uc.sessionRepo.GetAccessToken(ctx, token)
	if err != nil {
		return uuid.Nil, err
	}

	return userID, nil
}
