package commands

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Run tests",
	Long:  `Run tests with various options.`,
	Run: func(cmd *cobra.Command, args []string) {
		verbose, _ := cmd.Flags().GetBool("verbose")
		cover, _ := cmd.Flags().GetBool("cover")
		race, _ := cmd.Flags().GetBool("race")

		testArgs := []string{"test"}

		if verbose {
			testArgs = append(testArgs, "-v")
		}
		if cover {
			testArgs = append(testArgs, "-cover")
		}
		if race {
			testArgs = append(testArgs, "-race")
		}

		testArgs = append(testArgs, "./...")

		fmt.Println("Running tests...")
		testCmd := exec.Command("go", testArgs...)
		testCmd.Stdout = os.Stdout
		testCmd.Stderr = os.Stderr

		if err := testCmd.Run(); err != nil {
			os.Exit(1)
		}

		fmt.Println("\n✓ All tests passed")
	},
}

var testCoverageCmd = &cobra.Command{
	Use:   "coverage",
	Short: "Run tests with coverage report",
	Long:  `Run tests and generate HTML coverage report.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running tests with coverage...")

		// Run tests with coverage
		testCmd := exec.Command("go", "test", "-v", "-coverprofile=coverage.out", "./...")
		testCmd.Stdout = os.Stdout
		testCmd.Stderr = os.Stderr

		if err := testCmd.Run(); err != nil {
			fmt.Println("Tests failed")
			os.Exit(1)
		}

		// Generate HTML report
		fmt.Println("\nGenerating coverage report...")
		coverCmd := exec.Command("go", "tool", "cover", "-html=coverage.out", "-o", "coverage.html")
		if err := coverCmd.Run(); err != nil {
			fmt.Printf("Error generating report: %v\n", err)
			os.Exit(1)
		}

		// Show total coverage
		totalCmd := exec.Command("go", "tool", "cover", "-func=coverage.out")
		output, _ := totalCmd.Output()
		fmt.Println(string(output))

		fmt.Println("✓ Coverage report: coverage.html")
	},
}

var lintCmd = &cobra.Command{
	Use:   "lint",
	Short: "Run linter",
	Long:  `Run golangci-lint to check code quality.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running linter...")

		lintCmd := exec.Command("golangci-lint", "run")
		lintCmd.Stdout = os.Stdout
		lintCmd.Stderr = os.Stderr

		if err := lintCmd.Run(); err != nil {
			os.Exit(1)
		}

		fmt.Println("✓ Linting completed")
	},
}

var fmtCmd = &cobra.Command{
	Use:   "fmt",
	Short: "Format code",
	Long:  `Format all Go code using gofmt.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Formatting code...")

		fmtCmd := exec.Command("go", "fmt", "./...")
		fmtCmd.Stdout = os.Stdout
		fmtCmd.Stderr = os.Stderr

		if err := fmtCmd.Run(); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("✓ Code formatted")
	},
}

var vetCmd = &cobra.Command{
	Use:   "vet",
	Short: "Run go vet",
	Long:  `Run go vet to examine Go source code and report suspicious constructs.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running go vet...")

		vetCmd := exec.Command("go", "vet", "./...")
		vetCmd.Stdout = os.Stdout
		vetCmd.Stderr = os.Stderr

		if err := vetCmd.Run(); err != nil {
			os.Exit(1)
		}

		fmt.Println("✓ Vet completed")
	},
}

var mocksCmd = &cobra.Command{
	Use:   "mocks",
	Short: "Generate mocks for testing",
	Long:  `Generate mocks using mockery.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Generating mocks...")

		mocksCmd := exec.Command("mockery")
		mocksCmd.Stdout = os.Stdout
		mocksCmd.Stderr = os.Stderr

		if err := mocksCmd.Run(); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("✓ Mocks generated in internal/domain/repository/mocks/")
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
	rootCmd.AddCommand(lintCmd)
	rootCmd.AddCommand(fmtCmd)
	rootCmd.AddCommand(vetCmd)
	rootCmd.AddCommand(mocksCmd)

	testCmd.AddCommand(testCoverageCmd)

	testCmd.Flags().BoolP("verbose", "v", false, "verbose output")
	testCmd.Flags().BoolP("cover", "c", false, "show coverage")
	testCmd.Flags().BoolP("race", "r", false, "enable race detector")
}
