package commands

import (
	"fmt"
	"log"

	"github.com/alpinnz/go-rest-api-boilerplate/config"
	"github.com/alpinnz/go-rest-api-boilerplate/internal/infrastructure/database"
	"github.com/alpinnz/go-rest-api-boilerplate/pkg/migration"
	"github.com/spf13/cobra"
)

var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "Database management commands",
	Long:  `Manage database migrations and seeding.`,
}

var migrateUpCmd = &cobra.Command{
	Use:   "migrate-up",
	Short: "Run all pending migrations",
	Long:  `Apply all pending database migrations.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running migrations...")

		cfg, err := config.Load()
		if err != nil {
			log.Fatalf("Failed to load config: %v", err)
		}

		poolConfig := database.PoolConfig{
			MaxOpenConns:    cfg.Database.MaxOpenConns,
			MaxIdleConns:    cfg.Database.MaxIdleConns,
			ConnMaxLifetime: cfg.Database.ConnMaxLifetime,
			ConnMaxIdleTime: cfg.Database.ConnMaxIdleTime,
		}

		db, err := database.NewPostgresDB(cfg.Database.DSN(), poolConfig)
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}
		defer db.Close()

		runner := migration.NewRunner(db, "migrations")
		if err := runner.Up(); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
	},
}

var migrateDownCmd = &cobra.Command{
	Use:   "migrate-down",
	Short: "Rollback last migration",
	Long:  `Rollback the last applied database migration.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Rolling back last migration...")

		cfg, err := config.Load()
		if err != nil {
			log.Fatalf("Failed to load config: %v", err)
		}

		poolConfig := database.PoolConfig{
			MaxOpenConns:    cfg.Database.MaxOpenConns,
			MaxIdleConns:    cfg.Database.MaxIdleConns,
			ConnMaxLifetime: cfg.Database.ConnMaxLifetime,
			ConnMaxIdleTime: cfg.Database.ConnMaxIdleTime,
		}

		db, err := database.NewPostgresDB(cfg.Database.DSN(), poolConfig)
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}
		defer db.Close()

		runner := migration.NewRunner(db, "migrations")
		if err := runner.Down(); err != nil {
			log.Fatalf("Rollback failed: %v", err)
		}
	},
}

var migrateStatusCmd = &cobra.Command{
	Use:   "migrate-status",
	Short: "Show migration status",
	Long:  `Show status of all migrations.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			log.Fatalf("Failed to load config: %v", err)
		}

		poolConfig := database.PoolConfig{
			MaxOpenConns:    cfg.Database.MaxOpenConns,
			MaxIdleConns:    cfg.Database.MaxIdleConns,
			ConnMaxLifetime: cfg.Database.ConnMaxLifetime,
			ConnMaxIdleTime: cfg.Database.ConnMaxIdleTime,
		}

		db, err := database.NewPostgresDB(cfg.Database.DSN(), poolConfig)
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}
		defer db.Close()

		runner := migration.NewRunner(db, "migrations")
		if err := runner.Status(); err != nil {
			log.Fatalf("Failed to get status: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(dbCmd)
	dbCmd.AddCommand(migrateUpCmd)
	dbCmd.AddCommand(migrateDownCmd)
	dbCmd.AddCommand(migrateStatusCmd)
}
