package usecase

import (
	"context"
	"time"

	"github.com/alpinnz/go-rest-api-boilerplate/internal/domain"
	"github.com/alpinnz/go-rest-api-boilerplate/pkg/auth"
)

type UserUseCase struct {
	userRepo       domain.UserRepository
	sessionRepo    domain.SessionRepository
	jwtExpiration  time.Duration
	contextTimeout time.Duration
}

func NewUserUseCase(
	userRepo domain.UserRepository,
	sessionRepo domain.SessionRepository,
	jwtExpiration time.Duration,
	timeout time.Duration,
) *UserUseCase {
	return &UserUseCase{
		userRepo:       userRepo,
		sessionRepo:    sessionRepo,
		jwtExpiration:  jwtExpiration,
		contextTimeout: timeout,
	}
}

type RegisterInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Name     string `json:"name" validate:"required,min=2"`
}

type LoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	User         *domain.User `json:"user"`
}

func (uc *UserUseCase) Register(ctx context.Context, input RegisterInput) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()

	if input.Email == "" || input.Password == "" || input.Name == "" {
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

	user := &domain.User{
		Email:     input.Email,
		Password:  hashedPassword,
		Name:      input.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (uc *UserUseCase) Login(ctx context.Context, input LoginInput) (*LoginResponse, error) {
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
	if err := uc.sessionRepo.Set(ctx, accessToken, user.ID, auth.AccessTokenExpiration); err != nil {
		return nil, err
	}

	// Store refresh token mapping
	if err := uc.sessionRepo.SetRefreshToken(ctx, refreshToken, accessToken, auth.RefreshTokenExpiration); err != nil {
		return nil, err
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}

func (uc *UserUseCase) Logout(ctx context.Context, token string) error {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()

	return uc.sessionRepo.Delete(ctx, token)
}

// RefreshToken generates new access token from valid refresh token.
// Returns new access token and updated refresh token.
func (uc *UserUseCase) RefreshToken(ctx context.Context, refreshToken string) (*LoginResponse, error) {
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
	_ = uc.sessionRepo.Delete(ctx, oldAccessToken)

	// Delete used refresh token
	_ = uc.sessionRepo.DeleteRefreshToken(ctx, refreshToken)

	// Get user data
	user, err := uc.userRepo.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, domain.ErrNotFound
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
	if err := uc.sessionRepo.Set(ctx, newAccessToken, user.ID, auth.AccessTokenExpiration); err != nil {
		return nil, err
	}

	// Store new refresh token mapping
	if err := uc.sessionRepo.SetRefreshToken(ctx, newRefreshToken, newAccessToken, auth.RefreshTokenExpiration); err != nil {
		return nil, err
	}

	return &LoginResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		User:         user,
	}, nil
}

func (uc *UserUseCase) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()

	user, err := uc.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, domain.ErrNotFound
	}

	return user, nil
}

func (uc *UserUseCase) ValidateSession(ctx context.Context, token string) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()

	userID, err := uc.sessionRepo.Get(ctx, token)
	if err != nil {
		return 0, err
	}

	return userID, nil
}
