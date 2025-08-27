package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRole struct {
	ID        uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	UserID    uuid.UUID      `gorm:"type:char(36);index" json:"user_id"`
	User      *User          `gorm:"foreignKey:UserID;references:ID" json:"user"`
	RoleID    uuid.UUID      `gorm:"type:char(36);index" json:"role_id"`
	Role      *Role          `gorm:"foreignKey:RoleID;references:ID" json:"role"`
	CreatedAt time.Time      `gorm:"type:timestamptz;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"type:timestamptz;autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"type:timestamptz;index" json:"deleted_at,omitempty"`
}

//func (e *UserRole) ToModel() *models.UserRole {
//	if e == nil {
//		return nil
//	}
//	return &models.UserRole{
//		ID:        e.ID.String(),
//		UserID:    e.UserID.String(),
//		User:      e.User.ToModel(),
//		RoleID:    e.RoleID.String(),
//		Role:      e.Role.ToModel(),
//		CreatedAt: e.CreatedAt,
//		UpdatedAt: e.UpdatedAt,
//		DeletedAt: utils.NullableTime(e.DeletedAt),
//	}
//}
