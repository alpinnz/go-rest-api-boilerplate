package seeder

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/alpinnz/go-rest-api-boilerplate/internal/domain/entity"
	"github.com/alpinnz/go-rest-api-boilerplate/internal/domain/repository"
)

// RoleSeeder handles role data seeding
type RoleSeeder struct {
	roleRepo repository.RoleRepository
}

// NewRoleSeeder creates a new role seeder instance
func NewRoleSeeder(roleRepo repository.RoleRepository) *RoleSeeder {
	return &RoleSeeder{
		roleRepo: roleRepo,
	}
}

// Seed creates default roles
func (s *RoleSeeder) Seed(ctx context.Context) error {
	log.Println("Seeding roles...")

	roles := s.getDefaultRoles()

	for _, roleData := range roles {
		if err := s.seedRole(ctx, roleData); err != nil {
			return err
		}
	}

	log.Println("Roles seeding completed!")
	return nil
}

// seedRole creates a single role if it doesn't exist
func (s *RoleSeeder) seedRole(ctx context.Context, roleData RoleData) error {
	// Check if role already exists
	existingRole, _ := s.roleRepo.FindByName(ctx, roleData.Name)
	if existingRole != nil {
		log.Printf("Role %s already exists, skipping...", roleData.Name)
		return nil
	}

	// Create role entity
	role := &entity.Role{
		Name:        roleData.Name,
		Description: roleData.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Save to database
	if err := s.roleRepo.Create(ctx, role); err != nil {
		return fmt.Errorf("failed to create role %s: %w", roleData.Name, err)
	}

	log.Printf("Created role: %s - %s (ID: %d)",
		role.Name, role.Description, role.ID)

	return nil
}

// getDefaultRoles returns the default roles to seed
func (s *RoleSeeder) getDefaultRoles() []RoleData {
	return []RoleData{
		{
			Name:        "admin",
			Description: "Administrator with full access",
		},
		{
			Name:        "user",
			Description: "Regular user with standard access",
		},
		{
			Name:        "moderator",
			Description: "Moderator with elevated access",
		},
	}
}
