package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "app",
	Short: "Go REST API Boilerplate CLI",
	Long: `A powerful CLI tool for managing your Go REST API Boilerplate project.

Examples:
  app gen handler user          # Generate user handler
  app gen repository product    # Generate product repository  
  app gen service auth          # Generate auth service
  app gen module order          # Generate complete order module
  app dev                       # Start development server
  app test                      # Run tests
  app migrate up                # Run migrations`,
	Version: "1.0.0",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
}
