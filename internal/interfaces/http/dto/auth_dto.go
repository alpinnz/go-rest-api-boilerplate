package dto

import (
	"time"
)

type AuthSession struct {
	UserID           *string   `json:"user_id,omitempty"`
	User             *User     `json:"user,omitempty"`
	AccessToken      string    `json:"access_token"`
	AccessExpiresAt  time.Time `json:"access_expires_at"`
	RefreshToken     string    `json:"refresh_token"`
	RefreshExpiresAt time.Time `json:"refresh_expires_at"`
}

type Login struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type Register struct {
	FirstName      string `json:"first_name" field:"first_name"  validate:"required,min=4,max=200,regex_name" `
	LastName       string `json:"last_name" field:"last_name"  validate:"required,min=4,max=200,regex_name" `
	Email          string `json:"email" field:"email" validate:"required,min=4,max=200,email" `
	Password       string `json:"password" field:"password" validate:"required,strong_password" `
	PasswordRepeat string `json:"password_repeat" field:"password_repeat" validate:"required,eqfield=Password" `
}
