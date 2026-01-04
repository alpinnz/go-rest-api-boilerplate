package commands

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Database migration commands",
	Long:  `Manage database migrations.`,
}

var migrateUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Run all pending migrations",
	Long:  `Apply all pending database migrations.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running migrations...")

		migrateCmd := exec.Command("make", "migrate-up")
		migrateCmd.Stdout = os.Stdout
		migrateCmd.Stderr = os.Stderr

		if err := migrateCmd.Run(); err != nil {
			fmt.Println("Migration failed")
			os.Exit(1)
		}

		fmt.Println("✓ Migrations completed")
	},
}

var migrateDownCmd = &cobra.Command{
	Use:   "down",
	Short: "Rollback all migrations",
	Long:  `Rollback all database migrations.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Rolling back migrations...")

		migrateCmd := exec.Command("make", "migrate-down")
		migrateCmd.Stdout = os.Stdout
		migrateCmd.Stderr = os.Stderr

		if err := migrateCmd.Run(); err != nil {
			fmt.Println("Rollback failed")
			os.Exit(1)
		}

		fmt.Println("✓ Rollback completed")
	},
}

var migrateStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show migration status",
	Long:  `Show all migration files.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Migration files:\n")

		statusCmd := exec.Command("make", "migrate-status")
		statusCmd.Stdout = os.Stdout
		statusCmd.Stderr = os.Stderr

		if err := statusCmd.Run(); err != nil {
			os.Exit(1)
		}
	},
}

var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Seed the database",
	Long:  `Run database seeder to populate initial data.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Seeding database...")

		seedCmd := exec.Command("make", "seed")
		seedCmd.Stdout = os.Stdout
		seedCmd.Stderr = os.Stderr

		if err := seedCmd.Run(); err != nil {
			fmt.Println("Seeding failed")
			os.Exit(1)
		}

		fmt.Println("\n✓ Database seeded successfully")
		fmt.Println("\nDefault users created:")
		fmt.Println("  - admin@example.com / !Password123 (admin)")
		fmt.Println("  - user@example.com / !Password123 (user)")
		fmt.Println("  - test@example.com / !Password123 (user + moderator)")
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
	rootCmd.AddCommand(seedCmd)

	migrateCmd.AddCommand(migrateUpCmd)
	migrateCmd.AddCommand(migrateDownCmd)
	migrateCmd.AddCommand(migrateStatusCmd)
}
