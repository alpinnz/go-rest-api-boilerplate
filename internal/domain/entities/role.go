package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Role Entity
type Role struct {
	ID        uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	Name      string         `gorm:"type:varchar(255);unique;not null" json:"name"`
	CreatedAt time.Time      `gorm:"type:timestamptz;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"type:timestamptz;autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"type:timestamptz;index" json:"deleted_at,omitempty"`
	Users     []*User        `gorm:"many2many:user_roles;joinForeignKey:RoleID;joinReferences:UserID" json:"-"`
}

//func (e *Role) ToModel() *models.Role {
//	if e == nil {
//		return nil
//	}
//	return &models.Role{
//		ID:        e.ID.String(),
//		Name:      e.Name,
//		CreatedAt: e.CreatedAt,
//		UpdatedAt: e.UpdatedAt,
//		DeletedAt: utils.NullableTime(e.DeletedAt),
//	}
//}
