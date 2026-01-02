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

type UserUseCase struct {
	userRepo       repository.UserRepository
	roleRepo       repository.RoleRepository
	contextTimeout time.Duration
}

func NewUserUseCase(
	userRepo repository.UserRepository,
	roleRepo repository.RoleRepository,
	timeout time.Duration,
) *UserUseCase {
	return &UserUseCase{
		userRepo:       userRepo,
		roleRepo:       roleRepo,
		contextTimeout: timeout,
	}
}

// UpdateUserInput represents input for updating user
type UpdateUserInput struct {
	Email     string
	FirstName string
	LastName  string
}

// ListUsersInput represents input for listing users with pagination
type ListUsersInput struct {
	Limit  int
	Offset int
}

// ListUsersResult represents paginated users result
type ListUsersResult struct {
	Users   []*entity.User
	Total   int64
	Limit   int
	Offset  int
	HasMore bool
}

func (uc *UserUseCase) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()

	user, err := uc.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, domain.ErrNotFound
	}

	// Load user roles
	roles, _ := uc.roleRepo.FindByUserID(ctx, user.ID)
	user.Roles = roles

	return user, nil
}

func (uc *UserUseCase) List(ctx context.Context, input ListUsersInput) (*ListUsersResult, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()

	if input.Limit <= 0 {
		input.Limit = 10
	}
	if input.Offset < 0 {
		input.Offset = 0
	}

	users, err := uc.userRepo.FindAll(ctx, input.Limit, input.Offset)
	if err != nil {
		return nil, err
	}

	total, err := uc.userRepo.Count(ctx)
	if err != nil {
		return nil, err
	}

	// Load roles for each user
	for _, user := range users {
		roles, _ := uc.roleRepo.FindByUserID(ctx, user.ID)
		user.Roles = roles
	}

	return &ListUsersResult{
		Users:   users,
		Total:   total,
		Limit:   input.Limit,
		Offset:  input.Offset,
		HasMore: int64(input.Offset+input.Limit) < total,
	}, nil
}

func (uc *UserUseCase) Update(ctx context.Context, id uuid.UUID, input UpdateUserInput) (*entity.User, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()

	user, err := uc.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, domain.ErrNotFound
	}

	// Check email uniqueness if email is being changed
	if input.Email != "" && input.Email != user.Email {
		existing, _ := uc.userRepo.FindByEmail(ctx, input.Email)
		if existing != nil {
			return nil, domain.ErrConflict
		}
		user.Email = input.Email
	}

	if input.FirstName != "" {
		user.FirstName = input.FirstName
	}

	if input.LastName != "" {
		user.LastName = input.LastName
	}

	user.UpdatedAt = time.Now()

	if err := uc.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	// Load user roles
	roles, _ := uc.roleRepo.FindByUserID(ctx, user.ID)
	user.Roles = roles

	return user, nil
}

func (uc *UserUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()

	_, err := uc.userRepo.FindByID(ctx, id)
	if err != nil {
		return domain.ErrNotFound
	}

	return uc.userRepo.Delete(ctx, id)
}

func (uc *UserUseCase) ChangePassword(ctx context.Context, id uuid.UUID, oldPassword, newPassword string) error {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()

	if oldPassword == "" || newPassword == "" {
		return domain.ErrInvalidInput
	}

	user, err := uc.userRepo.FindByID(ctx, id)
	if err != nil {
		return domain.ErrNotFound
	}

	if !auth.CheckPasswordHash(oldPassword, user.Password) {
		return domain.ErrInvalidCredentials
	}

	hashedPassword, err := auth.HashPassword(newPassword)
	if err != nil {
		return err
	}

	user.Password = hashedPassword
	user.UpdatedAt = time.Now()

	return uc.userRepo.Update(ctx, user)
}

func (uc *UserUseCase) AssignRole(ctx context.Context, userID uuid.UUID, roleID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()

	// Verify user exists
	_, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return domain.ErrNotFound
	}

	// Verify role exists
	_, err = uc.roleRepo.FindByID(ctx, roleID)
	if err != nil {
		return domain.ErrNotFound
	}

	return uc.roleRepo.AssignToUser(ctx, userID, roleID)
}

func (uc *UserUseCase) RemoveRole(ctx context.Context, userID, roleID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()

	return uc.roleRepo.RemoveFromUser(ctx, userID, roleID)
}
