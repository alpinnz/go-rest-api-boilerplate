package dto

import (
	"time"
)

type User struct {
	ID          string     `json:"id"`
	FirstName   string     `json:"first_name"`
	LastName    string     `json:"last_name"`
	Email       string     `json:"email"`
	Roles       []*Role    `json:"roles,omitempty"`
	ActivatedAt *time.Time `json:"activated_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

type UserCreate struct {
	FirstName      string `json:"first_name" field:"first_name"  validate:"required,min=4,max=200,regex_name" `
	LastName       string `json:"last_name" field:"last_name"  validate:"required,min=4,max=200,regex_name" `
	Email          string `json:"email" field:"email" validate:"required,min=4,max=200,email" `
	Password       string `json:"password" field:"password" validate:"required,strong_password" `
	PasswordRepeat string `json:"password_repeat" field:"password_repeat" validate:"required,eqfield=Password" `
}

type UserUpdate struct {
	FirstName string `json:"first_name" field:"first_name"  validate:"required,min=4,max=200,regex_name" `
	LastName  string `json:"last_name" field:"last_name"  validate:"required,min=4,max=200,regex_name" `
}

type UserChangePassword struct {
	Id             int    `json:"id" validate:"required"`
	Password       string `json:"password" field:"password" validate:"required,strong_password" `
	PasswordRepeat string `json:"password_repeat" field:"password_repeat" validate:"required,eqfield=Password" `
}

type UserSession struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Role     int    `json:"role"`
}
