package entity

import (
	"time"

	"github.com/google/uuid"
)

// User represents an authenticated user entity in the system.
// This is a pure domain entity - no JSON tags (serialization is DTO concern).
type User struct {
	ID        uuid.UUID
	Email     string
	Password  string // bcrypt hashed password
	FirstName string
	LastName  string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time // soft delete timestamp
	Roles     []*Role    // User's assigned roles (lazy loaded)
}

// FullName returns the full name of the user.
func (u *User) FullName() string {
	return u.FirstName + " " + u.LastName
}

// IsDeleted checks if the user is soft deleted.
func (u *User) IsDeleted() bool {
	return u.DeletedAt != nil
}

// SoftDelete marks the user as deleted.
func (u *User) SoftDelete() {
	now := time.Now()
	u.DeletedAt = &now
}
