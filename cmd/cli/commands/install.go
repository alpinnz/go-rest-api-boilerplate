package commands

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install development tools",
	Long:  `Install all required development tools.`,
}

var installToolsCmd = &cobra.Command{
	Use:   "tools",
	Short: "Install all development tools",
	Long:  `Install Air, Swag, golangci-lint, golang-migrate, and mockery.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Installing development tools...")
		fmt.Println()

		tools := []struct {
			name    string
			command []string
		}{
			{"Air (hot reload)", []string{"go", "install", "github.com/air-verse/air@latest"}},
			{"Swag (API docs)", []string{"go", "install", "github.com/swaggo/swag/cmd/swag@latest"}},
			{"Mockery (mocks)", []string{"go", "install", "github.com/vektra/mockery/v2@latest"}},
		}

		for _, tool := range tools {
			fmt.Printf("Installing %s...\n", tool.name)
			installCmd := exec.Command(tool.command[0], tool.command[1:]...)
			installCmd.Stdout = os.Stdout
			installCmd.Stderr = os.Stderr
			if err := installCmd.Run(); err != nil {
				fmt.Printf("  Warning: Failed to install %s\n", tool.name)
			}
		}

		// Install golangci-lint
		fmt.Println("Installing golangci-lint...")
		lintCmd := exec.Command("sh", "-c", "curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin")
		lintCmd.Stdout = os.Stdout
		lintCmd.Stderr = os.Stderr
		lintCmd.Run()

		// Install golang-migrate
		fmt.Println("Installing golang-migrate...")
		migrateCmd := exec.Command("brew", "list", "golang-migrate")
		if err := migrateCmd.Run(); err != nil {
			installMigrate := exec.Command("brew", "install", "golang-migrate")
			installMigrate.Stdout = os.Stdout
			installMigrate.Stderr = os.Stderr
			installMigrate.Run()
		}

		fmt.Println()
		fmt.Println("✓ All tools installed successfully!")
		fmt.Println()
		fmt.Println("Run 'app dev' to start development server")
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
	installCmd.AddCommand(installToolsCmd)
}
