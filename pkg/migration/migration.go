package migration

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// Migration represents a database migration
type Migration struct {
	Version   string
	Name      string
	UpSQL     string
	DownSQL   string
	Timestamp time.Time
}

// Runner handles migration execution
type Runner struct {
	db            *sql.DB
	migrationsDir string
}

// NewRunner creates a new migration runner
func NewRunner(db *sql.DB, migrationsDir string) *Runner {
	return &Runner{
		db:            db,
		migrationsDir: migrationsDir,
	}
}

// Initialize creates the migrations tracking table
func (r *Runner) Initialize() error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`
	_, err := r.db.Exec(query)
	return err
}

// Up runs all pending migrations
func (r *Runner) Up() error {
	if err := r.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize migrations table: %w", err)
	}

	migrations, err := r.loadMigrations()
	if err != nil {
		return err
	}

	appliedVersions, err := r.getAppliedVersions()
	if err != nil {
		return err
	}

	for _, migration := range migrations {
		if _, applied := appliedVersions[migration.Version]; !applied {
			if err := r.applyMigration(migration); err != nil {
				return fmt.Errorf("failed to apply migration %s: %w", migration.Version, err)
			}
			fmt.Printf("Applied migration: %s - %s\n", migration.Version, migration.Name)
		}
	}

	fmt.Println("All migrations completed successfully")
	return nil
}

// Down rolls back the last migration
func (r *Runner) Down() error {
	appliedVersions, err := r.getAppliedVersions()
	if err != nil {
		return err
	}

	if len(appliedVersions) == 0 {
		fmt.Println("No migrations to rollback")
		return nil
	}

	// Get the latest applied version
	var latestVersion string
	for version := range appliedVersions {
		if latestVersion == "" || version > latestVersion {
			latestVersion = version
		}
	}

	migrations, err := r.loadMigrations()
	if err != nil {
		return err
	}

	for _, migration := range migrations {
		if migration.Version == latestVersion {
			if err := r.rollbackMigration(migration); err != nil {
				return fmt.Errorf("failed to rollback migration %s: %w", migration.Version, err)
			}
			fmt.Printf("Rolled back migration: %s - %s\n", migration.Version, migration.Name)
			return nil
		}
	}

	return fmt.Errorf("migration %s not found", latestVersion)
}

// Status shows the status of all migrations
func (r *Runner) Status() error {
	migrations, err := r.loadMigrations()
	if err != nil {
		return err
	}

	appliedVersions, err := r.getAppliedVersions()
	if err != nil {
		return err
	}

	fmt.Println("Migration Status:")
	fmt.Println("=================")
	for _, migration := range migrations {
		status := "Pending"
		if _, applied := appliedVersions[migration.Version]; applied {
			status = "Applied"
		}
		fmt.Printf("%s | %s | %s\n", migration.Version, status, migration.Name)
	}

	return nil
}

// loadMigrations loads all migration files from the migrations directory
func (r *Runner) loadMigrations() ([]Migration, error) {
	files, err := os.ReadDir(r.migrationsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read migrations directory: %w", err)
	}

	migrationsMap := make(map[string]*Migration)

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filename := file.Name()
		if !strings.HasSuffix(filename, ".sql") {
			continue
		}

		// Parse migration filename (format: YYYYMMDDHHMMSS_migration_name.up.sql or .down.sql)
		parts := strings.Split(filename, "_")
		if len(parts) < 2 {
			continue
		}

		version := parts[0]
		isUp := strings.HasSuffix(filename, ".up.sql")
		isDown := strings.HasSuffix(filename, ".down.sql")

		if !isUp && !isDown {
			continue
		}

		// Get or create migration
		migration, exists := migrationsMap[version]
		if !exists {
			nameParts := parts[1:]
			name := strings.Join(nameParts, "_")
			name = strings.TrimSuffix(name, ".up.sql")
			name = strings.TrimSuffix(name, ".down.sql")

			migration = &Migration{
				Version: version,
				Name:    name,
			}
			migrationsMap[version] = migration
		}

		// Read SQL content
		content, err := os.ReadFile(filepath.Join(r.migrationsDir, filename))
		if err != nil {
			return nil, fmt.Errorf("failed to read migration file %s: %w", filename, err)
		}

		if isUp {
			migration.UpSQL = string(content)
		} else {
			migration.DownSQL = string(content)
		}
	}

	// Convert map to sorted slice
	migrations := make([]Migration, 0, len(migrationsMap))
	for _, migration := range migrationsMap {
		migrations = append(migrations, *migration)
	}

	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	return migrations, nil
}

// getAppliedVersions returns a set of applied migration versions
func (r *Runner) getAppliedVersions() (map[string]bool, error) {
	rows, err := r.db.Query("SELECT version FROM schema_migrations")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	versions := make(map[string]bool)
	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		versions[version] = true
	}

	return versions, rows.Err()
}

// applyMigration applies a single migration
func (r *Runner) applyMigration(migration Migration) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Execute migration SQL
	if _, err := tx.Exec(migration.UpSQL); err != nil {
		return err
	}

	// Record migration
	if _, err := tx.Exec("INSERT INTO schema_migrations (version) VALUES ($1)", migration.Version); err != nil {
		return err
	}

	return tx.Commit()
}

// rollbackMigration rolls back a single migration
func (r *Runner) rollbackMigration(migration Migration) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Execute rollback SQL
	if _, err := tx.Exec(migration.DownSQL); err != nil {
		return err
	}

	// Remove migration record
	if _, err := tx.Exec("DELETE FROM schema_migrations WHERE version = $1", migration.Version); err != nil {
		return err
	}

	return tx.Commit()
}
