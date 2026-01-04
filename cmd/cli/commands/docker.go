package commands

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var dockerCmd = &cobra.Command{
	Use:   "docker",
	Short: "Docker commands",
	Long:  `Manage Docker containers for the application.`,
}

var dockerUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Start Docker containers",
	Long:  `Start all services (PostgreSQL, Redis, API) using Docker Compose.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting Docker containers...")

		dockerCmd := exec.Command("docker", "compose", "up", "-d")
		dockerCmd.Stdout = os.Stdout
		dockerCmd.Stderr = os.Stderr

		if err := dockerCmd.Run(); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("\n✓ All services started")
		fmt.Println("\nServices running:")

		// Show status
		psCmd := exec.Command("docker", "compose", "ps")
		psCmd.Stdout = os.Stdout
		psCmd.Run()
	},
}

var dockerDownCmd = &cobra.Command{
	Use:   "down",
	Short: "Stop Docker containers",
	Long:  `Stop all Docker containers.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Stopping Docker containers...")

		dockerCmd := exec.Command("docker", "compose", "down")
		dockerCmd.Stdout = os.Stdout
		dockerCmd.Stderr = os.Stderr

		if err := dockerCmd.Run(); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("✓ All services stopped")
	},
}

var dockerLogsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Show Docker logs",
	Long:  `Show logs from Docker containers.`,
	Run: func(cmd *cobra.Command, args []string) {
		follow, _ := cmd.Flags().GetBool("follow")

		logsArgs := []string{"compose", "logs"}
		if follow {
			logsArgs = append(logsArgs, "-f")
		}

		fmt.Println("Showing Docker logs (Ctrl+C to exit)...")
		logsCmd := exec.Command("docker", logsArgs...)
		logsCmd.Stdout = os.Stdout
		logsCmd.Stderr = os.Stderr
		logsCmd.Stdin = os.Stdin

		logsCmd.Run()
	},
}

var dockerRebuildCmd = &cobra.Command{
	Use:   "rebuild",
	Short: "Rebuild and restart containers",
	Long:  `Rebuild Docker images and restart containers.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Rebuilding Docker containers...")

		// Stop containers
		downCmd := exec.Command("docker", "compose", "down")
		downCmd.Run()

		// Build
		buildCmd := exec.Command("docker", "compose", "build")
		buildCmd.Stdout = os.Stdout
		buildCmd.Stderr = os.Stderr
		if err := buildCmd.Run(); err != nil {
			fmt.Printf("Build failed: %v\n", err)
			os.Exit(1)
		}

		// Start
		upCmd := exec.Command("docker", "compose", "up", "-d")
		upCmd.Stdout = os.Stdout
		upCmd.Stderr = os.Stderr
		if err := upCmd.Run(); err != nil {
			fmt.Printf("Start failed: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("✓ Containers rebuilt and started")
	},
}

func init() {
	rootCmd.AddCommand(dockerCmd)

	dockerCmd.AddCommand(dockerUpCmd)
	dockerCmd.AddCommand(dockerDownCmd)
	dockerCmd.AddCommand(dockerLogsCmd)
	dockerCmd.AddCommand(dockerRebuildCmd)

	dockerLogsCmd.Flags().BoolP("follow", "f", false, "follow log output")
}
