package commands

import (
	"fmt"

	"github.com/alpinnz/go-rest-api-boilerplate/cmd/cli/commands/generator"
	"github.com/alpinnz/go-rest-api-boilerplate/templates"
	"github.com/spf13/cobra"
)

var genHandlerCmd = &cobra.Command{
	Use:   "handler [name]",
	Short: "Generate HTTP handler",
	Long:  `Generate HTTP handler with CRUD operations.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		config := generator.NewFeatureConfig(name)

		fmt.Printf("Generating handler for %s...\n", config.Name)

		dst := fmt.Sprintf("internal/delivery/http/handler/%s_handler.go", config.NameLower)

		if generator.FileExists(dst) {
			fmt.Printf("Warning: %s already exists\n", dst)
			return
		}

		if err := generator.WriteFromTemplate(templates.HandlerTemplate, dst); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		if err := generator.ReplaceInFile(dst, config); err != nil {
			fmt.Printf("Error replacing placeholders: %v\n", err)
			return
		}

		fmt.Printf("✓ Created: %s\n", dst)
		fmt.Printf("\nNext steps:\n")
		fmt.Printf("1. Implement handler methods in %s\n", dst)
		fmt.Printf("2. Register routes in internal/delivery/http/router/router.go\n")
		fmt.Printf("3. Run: app test\n")
	},
}

var genRepositoryCmd = &cobra.Command{
	Use:   "repository [name]",
	Short: "Generate repository interface and implementation",
	Long:  `Generate repository interface and implementation for data access.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		config := generator.NewFeatureConfig(name)

		fmt.Printf("Generating repository for %s...\n", config.Name)

		dstInterface := fmt.Sprintf("internal/domain/repository/%s_repository.go", config.NameLower)

		if !generator.FileExists(dstInterface) {
			if err := generator.WriteFromTemplate(templates.RepositoryInterfaceTemplate, dstInterface); err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
			if err := generator.ReplaceInFile(dstInterface, config); err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
			fmt.Printf("✓ Created: %s\n", dstInterface)
		} else {
			fmt.Printf("⊘ Skipped (exists): %s\n", dstInterface)
		}

		dstImpl := fmt.Sprintf("internal/repository/%s_repository.go", config.NameLower)

		if !generator.FileExists(dstImpl) {
			if err := generator.WriteFromTemplate(templates.RepositoryImplementationTemplate, dstImpl); err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
			if err := generator.ReplaceInFile(dstImpl, config); err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
			fmt.Printf("✓ Created: %s\n", dstImpl)
		} else {
			fmt.Printf("⊘ Skipped (exists): %s\n", dstImpl)
		}

		fmt.Printf("\nNext steps:\n")
		fmt.Printf("1. Implement repository methods in %s\n", dstImpl)
		fmt.Printf("2. Add custom query methods if needed\n")
		fmt.Printf("3. Run: app test\n")
	},
}

var genServiceCmd = &cobra.Command{
	Use:     "service [name]",
	Aliases: []string{"usecase"},
	Short:   "Generate service/use case",
	Long:    `Generate service (use case) with business logic.`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		config := generator.NewFeatureConfig(name)

		fmt.Printf("Generating service for %s...\n", config.Name)

		dst := fmt.Sprintf("internal/usecase/%s_usecase.go", config.NameLower)

		if generator.FileExists(dst) {
			fmt.Printf("Warning: %s already exists\n", dst)
			return
		}

		if err := generator.WriteFromTemplate(templates.UsecaseTemplate, dst); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		if err := generator.ReplaceInFile(dst, config); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		fmt.Printf("✓ Created: %s\n", dst)
		fmt.Printf("\nNext steps:\n")
		fmt.Printf("1. Implement business logic in %s\n", dst)
		fmt.Printf("2. Add validation rules\n")
		fmt.Printf("3. Run: app test\n")
	},
}

var genModuleCmd = &cobra.Command{
	Use:   "module [name]",
	Short: "Generate complete module (entity, repository, service, handler, DTO)",
	Long:  `Generate a complete module with all layers following Clean Architecture.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		config := generator.NewFeatureConfig(name)

		fmt.Printf("Generating complete module for %s...\n\n", config.Name)

		files := []struct {
			template string
			output   string
			desc     string
		}{
			{
				templates.EntityTemplate,
				fmt.Sprintf("internal/domain/entity/%s.go", config.NameLower),
				"Entity",
			},
			{
				templates.RepositoryInterfaceTemplate,
				fmt.Sprintf("internal/domain/repository/%s_repository.go", config.NameLower),
				"Repository Interface",
			},
			{
				templates.RepositoryImplementationTemplate,
				fmt.Sprintf("internal/repository/%s_repository.go", config.NameLower),
				"Repository Implementation",
			},
			{
				templates.UsecaseTemplate,
				fmt.Sprintf("internal/usecase/%s_usecase.go", config.NameLower),
				"Use Case",
			},
			{
				templates.DTOTemplate,
				fmt.Sprintf("internal/delivery/http/dto/%s_dto.go", config.NameLower),
				"DTO",
			},
			{
				templates.HandlerTemplate,
				fmt.Sprintf("internal/delivery/http/handler/%s_handler.go", config.NameLower),
				"Handler",
			},
		}

		for _, f := range files {
			if generator.FileExists(f.output) {
				fmt.Printf("⊘ Skipped (%s): %s\n", f.desc, f.output)
				continue
			}

			if err := generator.WriteFromTemplate(f.template, f.output); err != nil {
				fmt.Printf("✗ Error creating %s: %v\n", f.desc, err)
				continue
			}

			if err := generator.ReplaceInFile(f.output, config); err != nil {
				fmt.Printf("✗ Error processing %s: %v\n", f.desc, err)
				continue
			}

			fmt.Printf("✓ Created %s: %s\n", f.desc, f.output)
		}

		fmt.Printf("\n✓ Module %s generated successfully!\n\n", config.Name)
		fmt.Printf("Next steps:\n")
		fmt.Printf("1. Generate migration: app gen migration create_table_%s\n", config.NamePlural)
		fmt.Printf("2. Customize entity fields\n")
		fmt.Printf("3. Implement business logic\n")
		fmt.Printf("4. Register routes in router.go\n")
		fmt.Printf("5. Initialize in cmd/api/main.go\n")
		fmt.Printf("6. Run: app migrate up\n")
		fmt.Printf("7. Test: app test\n")
	},
}

var genMigrationCmd = &cobra.Command{
	Use:   "migration [name]",
	Short: "Generate database migration files",
	Long:  `Generate up and down migration files with timestamp.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		config := generator.NewFeatureConfig(name)

		fmt.Printf("Generating migration for %s...\n", name)

		upFile := fmt.Sprintf("migrations/%s_%s.up.sql", config.Timestamp, name)
		downFile := fmt.Sprintf("migrations/%s_%s.down.sql", config.Timestamp, name)

		if err := generator.WriteFromTemplate(templates.MigrationUpTemplate, upFile); err != nil {
			fmt.Printf("Error creating up migration: %v\n", err)
			return
		}
		if err := generator.ReplaceInFile(upFile, config); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		if err := generator.WriteFromTemplate(templates.MigrationDownTemplate, downFile); err != nil {
			fmt.Printf("Error creating down migration: %v\n", err)
			return
		}
		if err := generator.ReplaceInFile(downFile, config); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		fmt.Printf("✓ Created: %s\n", upFile)
		fmt.Printf("✓ Created: %s\n", downFile)
		fmt.Printf("\nNext steps:\n")
		fmt.Printf("1. Edit %s to add your schema\n", upFile)
		fmt.Printf("2. Edit %s to add rollback logic\n", downFile)
		fmt.Printf("3. Run: app migrate up\n")
	},
}

var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate code from templates",
	Long:  `Generate boilerplate code for handlers, repositories, services, and complete modules.`,
}

func init() {
	rootCmd.AddCommand(genCmd)

	genCmd.AddCommand(genHandlerCmd)
	genCmd.AddCommand(genRepositoryCmd)
	genCmd.AddCommand(genServiceCmd)
	genCmd.AddCommand(genModuleCmd)
	genCmd.AddCommand(genMigrationCmd)
}
