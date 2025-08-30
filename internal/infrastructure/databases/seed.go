package databases

import (
	"context"
	"fmt"
	"time"

	"github.com/alpinnz/go-rest-api-boilerplate/internal/domain/entities"
	"github.com/alpinnz/go-rest-api-boilerplate/pkg/constants"
	"github.com/alpinnz/go-rest-api-boilerplate/pkg/encrypt"
	"github.com/alpinnz/go-rest-api-boilerplate/pkg/helper"
	"github.com/alpinnz/go-rest-api-boilerplate/pkg/logger"
	"github.com/alpinnz/go-rest-api-boilerplate/pkg/utils"
	"gorm.io/gorm"
)

// SeedRoles inserts default roles with fixed IDs
func SeedRoles(ctx context.Context, db *gorm.DB) error {
	roles := []entities.Role{
		{ID: constants.RoleIDAdmin, Name: constants.RoleNameAdmin},
		{ID: constants.RoleIDUser, Name: constants.RoleNameUser},
	}

	for _, role := range roles {
		var count int64
		if err := db.WithContext(ctx).
			Model(&entities.Role{}).
			Where("id = ?", role.ID).
			Or("name = ?", role.Name).
			Count(&count).Error; err != nil {
			return fmt.Errorf("failed to check role %s: %w", role.Name, err)
		}

		if count == 0 {
			_, err := helper.WithGormTransaction(ctx, db, func(tx *gorm.DB) (interface{}, error) {
				if err := tx.WithContext(ctx).Create(&role).Error; err != nil {
					return nil, fmt.Errorf("failed to seed role %s: %w", role.Name, err)
				}
				logger.Log.Info("Seeded role", "id", role.ID, "name", role.Name)
				return role, nil
			})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// SeedUsers inserts default system users with roles
func SeedUsers(ctx context.Context, db *gorm.DB, passwordSecret string) error {
	now := time.Now()

	users := []entities.User{
		{
			ID:        utils.GenerateUUID(),
			Email:     "admin@system.com",
			FirstName: "Admin",
			LastName:  "System",
			Roles: []*entities.Role{
				{ID: constants.RoleIDAdmin, Name: constants.RoleNameAdmin},
				{ID: constants.RoleIDUser, Name: constants.RoleNameUser}, // admin wajib punya role user juga
			},
			Password: "!Password123",
		},
		{
			ID:        utils.GenerateUUID(),
			Email:     "user@system.com",
			FirstName: "User",
			LastName:  "System",
			Roles:     []*entities.Role{},
			Password:  "!Password123",
		},
	}

	for _, user := range users {
		var count int64
		if err := db.WithContext(ctx).
			Model(&entities.User{}).
			Where("id = ?", user.ID).
			Or("email = ?", user.Email).
			Count(&count).Error; err != nil {
			return fmt.Errorf("failed to check user %s: %w", user.Email, err)
		}

		if count == 0 {
			_, err := helper.WithGormTransaction(ctx, db, func(tx *gorm.DB) (interface{}, error) {
				// buat user
				password, err := encrypt.HashPassword(user.Password, passwordSecret)
				if err != nil {
					return nil, err
				}
				user.Password = password
				if err := tx.WithContext(ctx).Create(
					&entities.User{
						ID:          user.ID,
						Email:       user.Email,
						FirstName:   user.FirstName,
						LastName:    user.LastName,
						Password:    user.Password,
						ActivatedAt: &now,
						CreatedAt:   now,
						UpdatedAt:   now,
					}).Error; err != nil {
					return nil, fmt.Errorf("failed to seed user %s: %w", user.Email, err)
				}

				// jika user.Roles kosong → kasih default role user
				if len(user.Roles) == 0 {
					user.Roles = []*entities.Role{
						{ID: constants.RoleIDUser, Name: constants.RoleNameUser},
					}
				}

				// assign roles ke join table user_roles (manual)
				userRoles := make([]entities.UserRole, 0)
				for _, role := range user.Roles {
					userRoles = append(userRoles, entities.UserRole{
						ID:        utils.GenerateUUID(),
						UserID:    user.ID,
						RoleID:    role.ID,
						CreatedAt: now,
						UpdatedAt: now,
					})
				}

				// logging pakai slog, bukan fmt.Println
				logger.Log.Debug("Assigning roles to user",
					"user", user.Email,
					"roles", userRoles,
				)

				if err := tx.WithContext(ctx).Create(&userRoles).Error; err != nil {
					return nil, fmt.Errorf("failed to seed user roles for %s: %w", user.Email, err)
				}

				logger.Log.Info("Seeded user", "id", user.ID, "email", user.Email)
				return user, nil
			})
			if err != nil {
				return err
			}
		}
	}

	return nil
}
