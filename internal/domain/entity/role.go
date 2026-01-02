package entity

import (
	"time"

	"github.com/google/uuid"
)

// Role represents a role entity in the system.
// Users can have multiple roles for fine-grained access control.
type Role struct {
	ID          uuid.UUID
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time // soft delete timestamp
}

// IsDeleted checks if the role is soft deleted.
func (r *Role) IsDeleted() bool {
	return r.DeletedAt != nil
}

// SoftDelete marks the role as deleted.
func (r *Role) SoftDelete() {
	now := time.Now()
	r.DeletedAt = &now
}
