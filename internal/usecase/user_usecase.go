package usecase

import (
	"context"
	"time"

	"github.com/alpinnz/go-rest-api-boilerplate/internal/config"
	"github.com/alpinnz/go-rest-api-boilerplate/internal/domain/entities"
	"github.com/alpinnz/go-rest-api-boilerplate/internal/domain/repositories"
	"github.com/alpinnz/go-rest-api-boilerplate/internal/interfaces/http/dto"
	"github.com/alpinnz/go-rest-api-boilerplate/internal/interfaces/http/mapper"
	"github.com/alpinnz/go-rest-api-boilerplate/pkg/constants"
	"github.com/alpinnz/go-rest-api-boilerplate/pkg/encrypt"
	"github.com/alpinnz/go-rest-api-boilerplate/pkg/errors"
	"github.com/alpinnz/go-rest-api-boilerplate/pkg/helper"
	"github.com/alpinnz/go-rest-api-boilerplate/pkg/translations"
	"github.com/alpinnz/go-rest-api-boilerplate/pkg/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserUsecase defines the contract for user-related business logic
type UserUsecase interface {
	FindAll(ctx context.Context) ([]dto.User, error)
	FindAllWithPagination(ctx context.Context, limit int, offset int) (dto.PaginationData, error)
	FindByID(ctx context.Context, id uuid.UUID) (dto.User, error)
	Create(ctx context.Context, req *dto.Register) (*dto.User, error)
	Update(ctx context.Context, id uuid.UUID, req *dto.UserUpdate) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// UserUsecaseImpl is the concrete implementation of UserUsecase
type UserUsecaseImpl struct {
	Env            *config.Env
	DB             *gorm.DB
	Tr             *translations.Store
	UserRepository repositories.UserRepository
	RoleRepository repositories.RoleRepository     // Access role if needed
	UserRoleRepo   repositories.UserRoleRepository // Assign role to user
}

// NewUserUsecase initializes the usecase with dependencies
func NewUserUsecase(env *config.Env, db *gorm.DB, tr *translations.Store, userRepository repositories.UserRepository, roleRepository repositories.RoleRepository, userRoleRepo repositories.UserRoleRepository) UserUsecase {
	return &UserUsecaseImpl{
		Env:            env,
		DB:             db,
		Tr:             tr,
		UserRepository: userRepository,
		RoleRepository: roleRepository,
		UserRoleRepo:   userRoleRepo,
	}
}

// FindAll returns all users (default max 1000)
func (u *UserUsecaseImpl) FindAll(ctx context.Context) ([]dto.User, error) {
	res, err := helper.WithGormTransaction(ctx, u.DB, func(tx *gorm.DB) (interface{}, error) {
		users, err := u.UserRepository.List(ctx, tx, 10, 0, false)
		if err != nil {
			return nil, err
		}
		return mapper.ToUserResponses(users), nil
	})
	if err != nil {
		return nil, err
	}
	return res.([]dto.User), nil
}

// FindAllWithPagination returns paginated users
func (u *UserUsecaseImpl) FindAllWithPagination(ctx context.Context, limit int, offset int) (dto.PaginationData, error) {
	res, err := helper.WithGormTransaction(ctx, u.DB, func(tx *gorm.DB) (interface{}, error) {
		users, err := u.UserRepository.List(ctx, tx, limit, offset, true)
		if err != nil {
			return nil, err
		}
		total, err := u.UserRepository.Count(ctx, tx)
		if err != nil {
			return nil, err
		}
		return dto.PaginationData{
			TotalData: int(total),
			Data:      mapper.ToUserResponses(users),
		}, nil
	})
	if err != nil {
		return dto.PaginationData{}, err
	}
	return res.(dto.PaginationData), nil
}

// FindByID returns user details by UUID
func (u *UserUsecaseImpl) FindByID(ctx context.Context, id uuid.UUID) (dto.User, error) {
	res, err := helper.WithGormTransaction(ctx, u.DB, func(tx *gorm.DB) (interface{}, error) {
		user, err := u.UserRepository.GetByID(ctx, tx, id, true)
		if err != nil {
			return nil, err
		}
		return mapper.ToUserResponse(user), nil
	})
	if err != nil {
		return dto.User{}, err
	}
	return res.(dto.User), nil
}

// Create a new user with hashed password and assign default role (student)
func (u *UserUsecaseImpl) Create(ctx context.Context, req *dto.Register) (*dto.User, error) {
	res, err := helper.WithGormTransaction(ctx, u.DB, func(tx *gorm.DB) (interface{}, error) {
		// Check if email already exists
		existing, err := u.UserRepository.GetByEmail(ctx, tx, req.Email)
		if err == nil && existing.ID != uuid.Nil {
			msg := u.Tr.TContext(ctx, translations.APP_EMAIL_ALREADY_EXISTS, nil)
			return nil, errors.NewBadRequest(msg, err.Error())

		}

		// kalau error selain "not found" → langsung return
		if err != nil {
			if !errors.IsNotFound(err) {
				return nil, err
			}
		}
		// Hash password
		hashedPassword, err := encrypt.HashPassword(req.Password, u.Env.Auth.PasswordSecret)
		if err != nil {
			msg := u.Tr.TContext(ctx, translations.APP_INTERNAL_SERVER_ERROR, nil)
			return nil, errors.NewBadRequest(msg, err.Error())
		}
		now := time.Now()

		// Create new user entity
		user := entities.User{
			ID:          utils.GenerateUUID(),
			FirstName:   req.FirstName,
			LastName:    req.LastName,
			Email:       utils.EmailNormalizer(req.Email),
			Password:    hashedPassword,
			ActivatedAt: &now,
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		// Save user in DB
		createdUser, err := u.UserRepository.Create(ctx, tx, &user)
		if err != nil {
			return nil, err
		}

		// Assign default role "student"
		userRole := entities.UserRole{
			ID:        utils.GenerateUUID(),
			UserID:    createdUser.ID,
			RoleID:    constants.RoleIDUser,
			CreatedAt: now,
			UpdatedAt: now,
		}
		if err := u.UserRoleRepo.AssignRole(ctx, tx, &userRole); err != nil {
			return nil, err
		}

		roles, err := u.UserRoleRepo.GetRolesByUser(ctx, tx, createdUser.ID)
		if err != nil {
			if !errors.IsNotFound(err) {
				return nil, err
			}
		}
		createdUser.Roles = roles

		return mapper.ToUserResponse(createdUser), nil
	})

	if err != nil {
		return nil, err
	}

	return res.(*dto.User), nil
}

// Update modifies an existing user
func (u *UserUsecaseImpl) Update(ctx context.Context, id uuid.UUID, req *dto.UserUpdate) error {
	_, err := helper.WithGormTransaction(ctx, u.DB, func(tx *gorm.DB) (interface{}, error) {
		// Fetch user
		user, err := u.UserRepository.GetByID(ctx, tx, id, false)
		if err != nil {
			return nil, err
		}

		// Update fields
		user.FirstName = req.FirstName
		user.LastName = req.LastName

		// Save update
		if err := u.UserRepository.Update(ctx, tx, user); err != nil {
			return nil, err
		}
		return nil, nil
	})
	return err
}

// Delete removes a user by UUID
func (u *UserUsecaseImpl) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := helper.WithGormTransaction(ctx, u.DB, func(tx *gorm.DB) (interface{}, error) {
		// Ensure user exists
		_, err := u.UserRepository.GetByID(ctx, tx, id, false)
		if err != nil {
			return nil, err
		}

		// Delete user
		if err := u.UserRepository.Delete(ctx, tx, id); err != nil {
			return nil, err
		}
		return nil, nil
	})
	return err
}
