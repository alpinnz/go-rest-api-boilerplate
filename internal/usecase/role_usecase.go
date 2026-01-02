package usecase

import (
	"context"
	"time"

	"github.com/alpinnz/go-rest-api-boilerplate/internal/domain"
	"github.com/alpinnz/go-rest-api-boilerplate/internal/domain/entity"
	"github.com/alpinnz/go-rest-api-boilerplate/internal/domain/repository"
	"github.com/google/uuid"
)

type RoleUseCase struct {
	roleRepo       repository.RoleRepository
	contextTimeout time.Duration
}

func NewRoleUseCase(
	roleRepo repository.RoleRepository,
	timeout time.Duration,
) *RoleUseCase {
	return &RoleUseCase{
		roleRepo:       roleRepo,
		contextTimeout: timeout,
	}
}

// CreateRoleInput represents input for creating a role
type CreateRoleInput struct {
	Name        string
	Description string
}

// UpdateRoleInput represents input for updating a role
type UpdateRoleInput struct {
	Name        string
	Description string
}

// ListRolesInput represents input for listing roles with pagination
type ListRolesInput struct {
	Limit  int
	Offset int
}

// ListRolesResult represents paginated roles result
type ListRolesResult struct {
	Roles   []*entity.Role
	Total   int64
	Limit   int
	Offset  int
	HasMore bool
}

func (uc *RoleUseCase) Create(ctx context.Context, input CreateRoleInput) (*entity.Role, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()

	if input.Name == "" {
		return nil, domain.ErrInvalidInput
	}

	existing, _ := uc.roleRepo.FindByName(ctx, input.Name)
	if existing != nil {
		return nil, domain.ErrConflict
	}

	role := &entity.Role{
		Name:        input.Name,
		Description: input.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := uc.roleRepo.Create(ctx, role); err != nil {
		return nil, err
	}

	return role, nil
}

func (uc *RoleUseCase) GetByID(ctx context.Context, id uuid.UUID) (*entity.Role, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()

	role, err := uc.roleRepo.FindByID(ctx, id)
	if err != nil {
		return nil, domain.ErrNotFound
	}

	return role, nil
}

func (uc *RoleUseCase) List(ctx context.Context, input ListRolesInput) (*ListRolesResult, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()

	if input.Limit <= 0 {
		input.Limit = 10
	}
	if input.Offset < 0 {
		input.Offset = 0
	}

	roles, err := uc.roleRepo.FindAll(ctx, input.Limit, input.Offset)
	if err != nil {
		return nil, err
	}

	total, err := uc.roleRepo.Count(ctx)
	if err != nil {
		return nil, err
	}

	return &ListRolesResult{
		Roles:   roles,
		Total:   total,
		Limit:   input.Limit,
		Offset:  input.Offset,
		HasMore: int64(input.Offset+input.Limit) < total,
	}, nil
}

func (uc *RoleUseCase) Update(ctx context.Context, id uuid.UUID, input UpdateRoleInput) (*entity.Role, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()

	role, err := uc.roleRepo.FindByID(ctx, id)
	if err != nil {
		return nil, domain.ErrNotFound
	}

	// Check name uniqueness if name is being changed
	if input.Name != "" && input.Name != role.Name {
		existing, _ := uc.roleRepo.FindByName(ctx, input.Name)
		if existing != nil {
			return nil, domain.ErrConflict
		}
		role.Name = input.Name
	}

	if input.Description != "" {
		role.Description = input.Description
	}

	role.UpdatedAt = time.Now()

	if err := uc.roleRepo.Update(ctx, role); err != nil {
		return nil, err
	}

	return role, nil
}

func (uc *RoleUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()

	_, err := uc.roleRepo.FindByID(ctx, id)
	if err != nil {
		return domain.ErrNotFound
	}

	return uc.roleRepo.Delete(ctx, id)
}
