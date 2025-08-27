package databases

import (
	"fmt"

	"github.com/alpinnz/go-rest-api-boilerplate/internal/domain/entities"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&entities.User{},
		&entities.Role{},
		&entities.UserRole{},
		&entities.AuthSession{},
	)
	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}
	return nil
}
