package seeder

import "context"

// Seeder interface defines the contract for all seeders
type Seeder interface {
	Seed(ctx context.Context) error
}

// UserData represents user seed data
type UserData struct {
	Email     string
	Password  string
	FirstName string
	LastName  string
	Roles     []string // Role names to assign to this user
}

// RoleData represents role seed data
type RoleData struct {
	Name        string
	Description string
}

// UserRoleAssignment represents user-role assignment data
type UserRoleAssignment struct {
	UserEmail string
	RoleName  string
}
