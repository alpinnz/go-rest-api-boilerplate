package entities

import (
	"time"

	"github.com/google/uuid"
)

type AuthSession struct {
	ID               uuid.UUID `gorm:"type:uuid;primaryKey"              json:"id"`
	UserID           uuid.UUID `gorm:"type:uuid;not null"                json:"user_id"`
	User             *User     `gorm:"foreignKey:UserID"                 json:"user"`
	AccessToken      string    `gorm:"type:varchar(255);not null"        json:"access_token"`
	AccessExpiresAt  time.Time `gorm:"type:timestamptz;not null"         json:"access_expires_at"`
	RefreshToken     string    `gorm:"type:varchar(255);not null"        json:"refresh_token"`
	RefreshExpiresAt time.Time `gorm:"type:timestamptz;not null"         json:"refresh_expires_at"`
	CreatedAt        time.Time `gorm:"type:timestamptz;autoCreateTime"   json:"created_at"`
	UpdatedAt        time.Time `gorm:"type:timestamptz;autoUpdateTime"   json:"updated_at"`
}

func (t *AuthSession) IsAccessExpired() bool {
	return time.Now().UTC().After(t.AccessExpiresAt)
}
func (t *AuthSession) IsRefreshExpired() bool {
	return time.Now().UTC().After(t.RefreshExpiresAt)
}
