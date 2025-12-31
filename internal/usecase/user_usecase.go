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
	Token string       `json:"token"`
	User  *domain.User `json:"user"`
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

	token, err := auth.GenerateToken(user.ID)
	if err != nil {
		return nil, err
	}

	if err := uc.sessionRepo.Set(ctx, token, user.ID, uc.jwtExpiration); err != nil {
		return nil, err
	}

	return &LoginResponse{
		Token: token,
		User:  user,
	}, nil
}

func (uc *UserUseCase) Logout(ctx context.Context, token string) error {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()

	return uc.sessionRepo.Delete(ctx, token)
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
