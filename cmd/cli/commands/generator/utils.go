package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

type FeatureConfig struct {
	Name       string // e.g., "User"
	NameLower  string // e.g., "user"
	NameUpper  string // e.g., "USER"
	NamePlural string // e.g., "users"
	Timestamp  string // for migrations
}

func NewFeatureConfig(name string) *FeatureConfig {
	return &FeatureConfig{
		Name:       capitalize(name),
		NameLower:  strings.ToLower(name),
		NameUpper:  strings.ToUpper(name),
		NamePlural: pluralize(strings.ToLower(name)),
		Timestamp:  time.Now().Format("20060102150405"),
	}
}

func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + strings.ToLower(s[1:])
}

func pluralize(s string) string {
	if strings.HasSuffix(s, "y") {
		return s[:len(s)-1] + "ies"
	}
	if strings.HasSuffix(s, "s") {
		return s + "es"
	}
	return s + "s"
}

func GenerateFromTemplate(templatePath, outputPath string, config *FeatureConfig) error {
	// Read template
	tmplContent, err := os.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("failed to read template: %w", err)
	}

	// Parse template
	tmpl, err := template.New("gen").Parse(string(tmplContent))
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	// Create output directory if not exists
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Create output file
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Execute template
	if err := tmpl.Execute(file, config); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
}

func ReplaceInFile(filePath string, config *FeatureConfig) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	s := string(content)

	// Replace PLACEHOLDER_ style (new, valid Go syntax)
	s = strings.ReplaceAll(s, "PLACEHOLDER_Entity", config.Name)
	s = strings.ReplaceAll(s, "PLACEHOLDER_entity", config.NameLower)
	s = strings.ReplaceAll(s, "PLACEHOLDER_ENTITY", config.NameUpper)
	s = strings.ReplaceAll(s, "PLACEHOLDER_Entities", config.Name+"s") // For plural in struct fields
	s = strings.ReplaceAll(s, "PLACEHOLDER_entities", config.NamePlural)

	// Replace {{}} style placeholders (backwards compatibility)
	s = strings.ReplaceAll(s, "{{Entity}}", config.Name)
	s = strings.ReplaceAll(s, "{{entity}}", config.NameLower)
	s = strings.ReplaceAll(s, "{{ENTITY}}", config.NameUpper)
	s = strings.ReplaceAll(s, "{{entities}}", config.NamePlural)

	// Replace normal placeholders (backwards compatibility)
	s = strings.ReplaceAll(s, "Feature", config.Name)
	s = strings.ReplaceAll(s, "feature", config.NameLower)
	s = strings.ReplaceAll(s, "FEATURE", config.NameUpper)
	s = strings.ReplaceAll(s, "features", config.NamePlural)

	return os.WriteFile(filePath, []byte(s), 0644)
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func CopyFile(src, dst string) error {
	content, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	dir := filepath.Dir(dst)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return os.WriteFile(dst, []byte(content), 0644)
}

// WriteFromTemplate writes content from template string to destination file
func WriteFromTemplate(templateContent, dst string) error {
	dir := filepath.Dir(dst)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return os.WriteFile(dst, []byte(templateContent), 0644)
}
