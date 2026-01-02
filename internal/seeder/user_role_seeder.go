package seeder

import (
	"context"
	"fmt"
	"log"

	"github.com/alpinnz/go-rest-api-boilerplate/internal/domain/repository"
)

// UserRoleSeeder handles user-role assignment seeding
type UserRoleSeeder struct {
	userRepo repository.UserRepository
	roleRepo repository.RoleRepository
}

// NewUserRoleSeeder creates a new user-role seeder instance
func NewUserRoleSeeder(userRepo repository.UserRepository, roleRepo repository.RoleRepository) *UserRoleSeeder {
	return &UserRoleSeeder{
		userRepo: userRepo,
		roleRepo: roleRepo,
	}
}

// Seed assigns roles to users
func (s *UserRoleSeeder) Seed(ctx context.Context) error {
	log.Println("Seeding user-role assignments...")

	assignments := s.getDefaultAssignments()

	for _, assignment := range assignments {
		if err := s.assignRole(ctx, assignment); err != nil {
			return err
		}
	}

	log.Println("User-role assignments completed!")
	return nil
}

// assignRole assigns a role to a user
func (s *UserRoleSeeder) assignRole(ctx context.Context, assignment UserRoleAssignment) error {
	// Find user by email
	user, err := s.userRepo.FindByEmail(ctx, assignment.UserEmail)
	if err != nil {
		return fmt.Errorf("failed to find user %s: %w", assignment.UserEmail, err)
	}

	// Find role by name
	role, err := s.roleRepo.FindByName(ctx, assignment.RoleName)
	if err != nil {
		return fmt.Errorf("failed to find role %s: %w", assignment.RoleName, err)
	}

	// Assign role to user
	if err := s.roleRepo.AssignToUser(ctx, user.ID, role.ID); err != nil {
		// Check if it's a duplicate assignment (conflict error)
		// Just log and continue, don't fail the seeding
		log.Printf("Role assignment exists or failed for %s -> %s: %v",
			assignment.UserEmail, assignment.RoleName, err)
		return nil
	}

	log.Printf("Assigned role '%s' to user '%s' (User ID: %d, Role ID: %d)",
		role.Name, user.Email, user.ID, role.ID)

	return nil
}

// getDefaultAssignments returns the default user-role assignments
func (s *UserRoleSeeder) getDefaultAssignments() []UserRoleAssignment {
	return []UserRoleAssignment{
		{
			UserEmail: "admin@example.com",
			RoleName:  "admin",
		},
		{
			UserEmail: "admin@example.com",
			RoleName:  "user",
		},
		{
			UserEmail: "user@example.com",
			RoleName:  "user",
		},
		{
			UserEmail: "test@example.com",
			RoleName:  "user",
		},
		{
			UserEmail: "test@example.com",
			RoleName:  "moderator",
		},
	}
}
