package helper

import (
	"github.com/alpinnz/go-rest-api-boilerplate/internal/domain/entities"
	"github.com/google/uuid"
)

// ContainsRole helper to check if role exists
func ContainsRole(roles []entities.Role, roleID uuid.UUID) bool {
	for _, r := range roles {
		if r.ID == roleID {
			return true
		}
	}
	return false
}
