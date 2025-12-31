package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/alpinnz/go-rest-api-boilerplate/config"
	"github.com/alpinnz/go-rest-api-boilerplate/internal/domain"
	"github.com/alpinnz/go-rest-api-boilerplate/internal/infrastructure/database"
	"github.com/alpinnz/go-rest-api-boilerplate/internal/repository"
	"github.com/alpinnz/go-rest-api-boilerplate/pkg/auth"
)

type Seeder struct {
	db       *sql.DB
	userRepo domain.UserRepository
}

func NewSeeder(db *sql.DB) *Seeder {
	return &Seeder{
		db:       db,
		userRepo: repository.NewUserRepository(db),
	}
}

func (s *Seeder) SeedUsers(ctx context.Context) error {
	log.Println("Seeding users...")

	users := []struct {
		Email    string
		Password string
		Name     string
	}{
		{
			Email:    "admin@example.com",
			Password: "!Password123",
			Name:     "Admin User",
		},
		{
			Email:    "user@example.com",
			Password: "!Password123",
			Name:     "Regular User",
		},
		{
			Email:    "test@example.com",
			Password: "!Password123",
			Name:     "Test User",
		},
	}

	for _, userData := range users {
		// Check if user already exists
		existingUser, _ := s.userRepo.FindByEmail(ctx, userData.Email)
		if existingUser != nil {
			log.Printf("User %s already exists, skipping...", userData.Email)
			continue
		}

		// Hash password
		hashedPassword, err := auth.HashPassword(userData.Password)
		if err != nil {
			return fmt.Errorf("failed to hash password for %s: %w", userData.Email, err)
		}

		// Create user
		user := &domain.User{
			Email:     userData.Email,
			Password:  hashedPassword,
			Name:      userData.Name,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := s.userRepo.Create(ctx, user); err != nil {
			return fmt.Errorf("failed to create user %s: %w", userData.Email, err)
		}

		log.Printf("Created user: %s (ID: %d)", user.Email, user.ID)
	}

	log.Println("Users seeding completed!")
	return nil
}

func (s *Seeder) TruncateUsers(ctx context.Context) error {
	log.Println("Truncating users table...")
	_, err := s.db.ExecContext(ctx, "TRUNCATE TABLE users RESTART IDENTITY CASCADE")
	if err != nil {
		return fmt.Errorf("failed to truncate users: %w", err)
	}
	log.Println("Users table truncated!")
	return nil
}

func (s *Seeder) SeedAll(ctx context.Context, fresh bool) error {
	if fresh {
		if err := s.TruncateUsers(ctx); err != nil {
			return err
		}
	}

	if err := s.SeedUsers(ctx); err != nil {
		return err
	}

	return nil
}

func main() {
	log.Println("Starting seeder...")

	// Load config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to database
	db, err := database.NewPostgresDB(cfg.Database.DSN())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Create seeder
	seeder := NewSeeder(db)
	ctx := context.Background()

	// Run seeder
	// Use fresh = true to truncate tables first
	fresh := true
	if err := seeder.SeedAll(ctx, fresh); err != nil {
		log.Fatalf("Seeder failed: %v", err)
	}

	log.Println("Seeding completed successfully!")
}
