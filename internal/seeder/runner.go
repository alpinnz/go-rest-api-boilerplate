package seeder

import (
	"context"
	"database/sql"
	"fmt"
	"log"
)

// Runner orchestrates the seeding process
type Runner struct {
	db      *sql.DB
	seeders []Seeder
}

// NewRunner creates a new seeder runner
func NewRunner(db *sql.DB) *Runner {
	return &Runner{
		db:      db,
		seeders: make([]Seeder, 0),
	}
}

// Register adds a seeder to the runner
func (r *Runner) Register(seeder Seeder) {
	r.seeders = append(r.seeders, seeder)
}

// Run executes all registered seeders
func (r *Runner) Run(ctx context.Context, options RunOptions) error {
	log.Println("Starting seeder...")

	// Truncate tables if fresh mode is enabled
	if options.Fresh {
		if err := r.truncateAll(ctx); err != nil {
			return fmt.Errorf("failed to truncate tables: %w", err)
		}
	}

	// Run all seeders in order
	for _, seeder := range r.seeders {
		if err := seeder.Seed(ctx); err != nil {
			return fmt.Errorf("seeding failed: %w", err)
		}
	}

	log.Println("Seeding completed successfully!")
	return nil
}

// truncateAll truncates all seeder tables
func (r *Runner) truncateAll(ctx context.Context) error {
	log.Println("Truncating tables...")

	// Order matters: truncate child tables first.
	// Postgres doesn't support TRUNCATE TABLE IF EXISTS, so we probe table existence first.
	tables := []string{
		"user_roles",
		"roles",
		"users",
	}

	for _, table := range tables {
		exists, err := r.tableExists(ctx, table)
		if err != nil {
			return fmt.Errorf("failed to check table %s existence: %w", table, err)
		}
		if !exists {
			log.Printf("Skipping truncate (table does not exist yet): %s", table)
			continue
		}

		query := fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", table)
		if _, err := r.db.ExecContext(ctx, query); err != nil {
			return fmt.Errorf("failed to truncate %s: %w", table, err)
		}
		log.Printf("Truncated table: %s", table)
	}

	log.Println("All tables truncated!")
	return nil
}

func (r *Runner) tableExists(ctx context.Context, table string) (bool, error) {
	// to_regclass returns NULL if the relation doesn't exist
	var regclass sql.NullString
	if err := r.db.QueryRowContext(ctx, "SELECT to_regclass($1)", table).Scan(&regclass); err != nil {
		return false, err
	}
	return regclass.Valid && regclass.String != "", nil
}

// RunOptions configures the seeding behavior
type RunOptions struct {
	// Fresh truncates all tables before seeding
	Fresh bool
}
