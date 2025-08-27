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
	"github.com/alpinnz/go-rest-api-boilerplate/pkg/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AuthUsecase defines the contract for authentication business logic.
type AuthUsecase interface {
	// Login authenticates a user using input login,
	// generates JWT tokens, stores the session, and returns token info.
	Login(ctx context.Context, req dto.Login) (*dto.AuthSession, error)
	// Register authenticates a user using input register,
	// generates JWT tokens, stores the session, and returns token info.
	Register(ctx context.Context, req dto.Register) (*dto.AuthSession, error)
}

// AuthUsecaseImpl is the concrete implementation of AuthUsecase.
type AuthUsecaseImpl struct {
	Env                   *config.Env
	DB                    *gorm.DB
	UserRepository        repositories.UserRepository
	AuthSessionRepository repositories.AuthSessionRepository
	RoleRepository        repositories.RoleRepository     // Access role if needed
	UserRoleRepo          repositories.UserRoleRepository // Assign role to user
}

// NewAuthUsecase creates a new instance of AuthUsecaseImpl.
func NewAuthUsecase(
	env *config.Env,
	db *gorm.DB,
	userRepository repositories.UserRepository,
	authSessionRepository repositories.AuthSessionRepository,
	roleRepository repositories.RoleRepository,
	userRoleRepo repositories.UserRoleRepository,
) AuthUsecase {
	return &AuthUsecaseImpl{
		DB:                    db,
		Env:                   env,
		UserRepository:        userRepository,
		AuthSessionRepository: authSessionRepository,
		RoleRepository:        roleRepository,
		UserRoleRepo:          userRoleRepo,
	}
}

// Login handles the login flow:
//
// 1. Validate credentials (email & password).
// 2. Generate JWT access and refresh tokens.
// 3. Save authentication session.
// 4. Return the session response.
func (u *AuthUsecaseImpl) Login(ctx context.Context, req dto.Login) (*dto.AuthSession, error) {
	res, err := helper.WithGormTransaction(ctx, u.DB, func(tx *gorm.DB) (interface{}, error) {
		// Step 1: Find user by email.
		user, err := u.UserRepository.GetByEmail(ctx, tx, utils.EmailNormalizer(req.Email))
		if err != nil {
			// Do not expose whether email exists.
			return nil, errors.NewBadRequest("invalid email or password", err.Error())
		}

		// Step 2: Validate password.
		if _, err = encrypt.ComparePassword(user.Password, req.Password, u.Env.Auth.PasswordSecret); err != nil {
			return nil, errors.NewBadRequest("invalid email or password", err.Error())
		}

		// Step 3: Generate JWT tokens.
		now := time.Now().UTC()
		accessExpiresAt := now.Add(8 * time.Hour)
		refreshExpiresAt := now.Add(7 * 24 * time.Hour)

		accessClaims := encrypt.BuildClaims(user.ID, []entities.Role{}, now, accessExpiresAt)
		refreshClaims := encrypt.BuildClaims(user.ID, []entities.Role{}, now, refreshExpiresAt)

		accessToken, err := encrypt.GenerateHash(accessClaims, u.Env.Auth.AccessTokenSecret)
		if err != nil {
			return nil, errors.NewInternalError("failed to generate access token", err.Error())
		}
		refreshToken, err := encrypt.GenerateHash(refreshClaims, u.Env.Auth.RefreshTokenSecret)
		if err != nil {
			return nil, errors.NewInternalError("failed to generate refresh token", err.Error())
		}

		// Step 4: Save session in DB.
		authSession := &entities.AuthSession{
			ID:               utils.GenerateUUID(),
			UserID:           user.ID,
			AccessToken:      accessToken,
			AccessExpiresAt:  accessExpiresAt,
			RefreshToken:     refreshToken,
			RefreshExpiresAt: refreshExpiresAt,
		}
		authSessionCreated, err := u.AuthSessionRepository.Create(ctx, tx, authSession)
		if err != nil {
			return nil, err
		}

		authSessionCreated.User = user
		result := mapper.ToAuthSessionResponse(authSessionCreated)

		// Step 5: Return response DTO.
		return result, nil
	})
	if err != nil {
		return nil, err
	}

	return res.(*dto.AuthSession), nil
}

func (u *AuthUsecaseImpl) Register(ctx context.Context, req dto.Register) (*dto.AuthSession, error) {
	res, err := helper.WithGormTransaction(ctx, u.DB, func(tx *gorm.DB) (interface{}, error) {

		// Check if email already exists
		existing, err := u.UserRepository.GetByEmail(ctx, tx, req.Email)
		if err == nil && existing.ID != uuid.Nil {
			return nil, errors.NewErrorUserEmailExist()
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
			return nil, errors.NewBadRequest("failed to hash password", err.Error())
		}
		now := time.Now().UTC()

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

		// Step 3: Generate JWT tokens.
		accessExpiresAt := now.Add(8 * time.Hour)
		refreshExpiresAt := now.Add(7 * 24 * time.Hour)

		accessClaims := encrypt.BuildClaims(user.ID, []entities.Role{}, now, accessExpiresAt)
		refreshClaims := encrypt.BuildClaims(user.ID, []entities.Role{}, now, refreshExpiresAt)

		accessToken, err := encrypt.GenerateHash(accessClaims, u.Env.Auth.AccessTokenSecret)
		if err != nil {
			return nil, errors.NewInternalError("failed to generate access token", err.Error())
		}
		refreshToken, err := encrypt.GenerateHash(refreshClaims, u.Env.Auth.RefreshTokenSecret)
		if err != nil {
			return nil, errors.NewInternalError("failed to generate refresh token", err.Error())
		}

		// Step 4: Save session in DB.
		authSession := &entities.AuthSession{
			ID:               utils.GenerateUUID(),
			UserID:           user.ID,
			AccessToken:      accessToken,
			AccessExpiresAt:  accessExpiresAt,
			RefreshToken:     refreshToken,
			RefreshExpiresAt: refreshExpiresAt,
		}
		authSessionCreated, err := u.AuthSessionRepository.Create(ctx, tx, authSession)
		if err != nil {
			return nil, err
		}
		user.Roles = roles
		authSessionCreated.User = &user
		result := mapper.ToAuthSessionResponse(authSessionCreated)

		// Step 5: Return response DTO.
		return result, nil
	})
	if err != nil {
		return nil, err
	}

	return res.(*dto.AuthSession), nil
}
