package seeder

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/alpinnz/go-rest-api-boilerplate/internal/domain/entity"
	"github.com/alpinnz/go-rest-api-boilerplate/internal/domain/repository"
	"github.com/alpinnz/go-rest-api-boilerplate/pkg/auth"
)

// UserSeeder handles user data seeding
type UserSeeder struct {
	userRepo repository.UserRepository
}

// NewUserSeeder creates a new user seeder instance
func NewUserSeeder(userRepo repository.UserRepository) *UserSeeder {
	return &UserSeeder{
		userRepo: userRepo,
	}
}

// Seed creates default users
func (s *UserSeeder) Seed(ctx context.Context) error {
	log.Println("Seeding users...")

	users := s.getDefaultUsers()

	for _, userData := range users {
		if err := s.seedUser(ctx, userData); err != nil {
			return err
		}
	}

	log.Println("Users seeding completed!")
	return nil
}

// seedUser creates a single user if it doesn't exist
func (s *UserSeeder) seedUser(ctx context.Context, userData UserData) error {
	// Check if user already exists
	existingUser, _ := s.userRepo.FindByEmail(ctx, userData.Email)
	if existingUser != nil {
		log.Printf("User %s already exists, skipping...", userData.Email)
		return nil
	}

	// Hash password
	hashedPassword, err := auth.HashPassword(userData.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password for %s: %w", userData.Email, err)
	}

	// Create user entity
	user := &entity.User{
		Email:     userData.Email,
		Password:  hashedPassword,
		FirstName: userData.FirstName,
		LastName:  userData.LastName,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Save to database
	if err := s.userRepo.Create(ctx, user); err != nil {
		return fmt.Errorf("failed to create user %s: %w", userData.Email, err)
	}

	log.Printf("Created user: %s %s - %s (ID: %d)",
		user.FirstName, user.LastName, user.Email, user.ID)

	return nil
}

// getDefaultUsers returns the default users to seed
func (s *UserSeeder) getDefaultUsers() []UserData {
	return []UserData{
		{
			Email:     "admin@example.com",
			Password:  "!Password123",
			FirstName: "Admin",
			LastName:  "User",
		},
		{
			Email:     "user@example.com",
			Password:  "!Password123",
			FirstName: "Regular",
			LastName:  "User",
		},
		{
			Email:     "test@example.com",
			Password:  "!Password123",
			FirstName: "Test",
			LastName:  "User",
		},
	}
}
