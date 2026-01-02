package main

import (
	"context"
	"log"

	"github.com/alpinnz/go-rest-api-boilerplate/config"
	"github.com/alpinnz/go-rest-api-boilerplate/internal/infrastructure/database"
	"github.com/alpinnz/go-rest-api-boilerplate/internal/repository"
	"github.com/alpinnz/go-rest-api-boilerplate/internal/seeder"
)

func main() {
	// Load configuration
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

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	roleRepo := repository.NewRoleRepository(db)

	// Create seeder runner
	runner := seeder.NewRunner(db)

	// Register seeders in order
	// 1. Seed roles first
	runner.Register(seeder.NewRoleSeeder(roleRepo))

	// 2. Seed users
	runner.Register(seeder.NewUserSeeder(userRepo))

	// 3. Assign roles to users
	runner.Register(seeder.NewUserRoleSeeder(userRepo, roleRepo))

	// Run seeding
	ctx := context.Background()
	options := seeder.RunOptions{
		Fresh: true, // Set to false to keep existing data
	}

	if err := runner.Run(ctx, options); err != nil {
		log.Fatalf("Seeder failed: %v", err)
	}
}
