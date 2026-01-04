package commands

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "Start development server with hot reload",
	Long:  `Start the development server using Air for hot reload.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting development server with hot reload...")
		fmt.Println("Press Ctrl+C to stop")
		fmt.Println()

		// Check if Air is installed
		if _, err := exec.LookPath("air"); err != nil {
			fmt.Println("Error: Air is not installed")
			fmt.Println("Run: app install tools")
			return
		}

		// Run Air
		airCmd := exec.Command("air")
		airCmd.Stdout = os.Stdout
		airCmd.Stderr = os.Stderr
		airCmd.Stdin = os.Stdin

		if err := airCmd.Run(); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	},
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the application",
	Long:  `Run the application without hot reload.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting application...")

		runCmd := exec.Command("go", "run", "cmd/api/main.go")
		runCmd.Stdout = os.Stdout
		runCmd.Stderr = os.Stderr
		runCmd.Stdin = os.Stdin

		if err := runCmd.Run(); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	},
}

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build the application",
	Long:  `Build the application binary.`,
	Run: func(cmd *cobra.Command, args []string) {
		output, _ := cmd.Flags().GetString("output")
		if output == "" {
			output = "bin/api"
		}

		fmt.Printf("Building application to %s...\n", output)

		buildCmd := exec.Command("go", "build", "-o", output, "cmd/api/main.go")
		buildCmd.Stdout = os.Stdout
		buildCmd.Stderr = os.Stderr

		if err := buildCmd.Run(); err != nil {
			fmt.Printf("Build failed: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✓ Binary created: %s\n", output)
	},
}

func init() {
	rootCmd.AddCommand(devCmd)
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(buildCmd)

	buildCmd.Flags().StringP("output", "o", "bin/api", "output binary path")
}
