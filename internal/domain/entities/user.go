package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User Entity
type User struct {
	ID          uuid.UUID      `gorm:"type:char(36);primaryKey;"                  json:"id"`
	FirstName   string         `gorm:"type:varchar(255);not null"                 json:"first_name"`
	LastName    string         `gorm:"type:varchar(255);not null"                 json:"last_name"`
	Email       string         `gorm:"type:varchar(255);uniqueIndex;not null"     json:"email"`
	Password    string         `gorm:"type:varchar(255);not null"                 json:"password"`
	ActivatedAt *time.Time     `gorm:"type:timestamptz"                           json:"activated_at,omitempty"`
	CreatedAt   time.Time      `gorm:"type:timestamptz;autoCreateTime"            json:"created_at"`
	UpdatedAt   time.Time      `gorm:"type:timestamptz;autoUpdateTime"            json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"type:timestamptz;index"                     json:"deleted_at,omitempty"`
	Roles       []*Role        `gorm:"many2many:user_roles;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"roles,omitempty"`
}

//func (e *User) ToModel() *models.User {
//	if e == nil {
//		return nil
//	}
//	return &models.User{
//		ID:          e.ID.String(),
//		FirstName:   e.FirstName,
//		LastName:    e.LastName,
//		Email:       e.Email,
//		ActivatedAt: e.ActivatedAt,
//		CreatedAt:   e.CreatedAt,
//		UpdatedAt:   e.UpdatedAt,
//		DeletedAt:   utils.NullableTime(e.DeletedAt),
//	}
//}
