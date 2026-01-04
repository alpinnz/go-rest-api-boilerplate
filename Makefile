.PHONY: help cli gen dev run build test lint fmt vet mocks migrate seed docker install

# Default target
.DEFAULT_GOAL := help

# Colors for output
BLUE := \033[0;34m
GREEN := \033[0;32m
YELLOW := \033[0;33m
RED := \033[0;31m
NC := \033[0m # No Color

help: ## Show this help message
	@echo "$(BLUE)Go REST API Boilerplate - Task Runner$(NC)"
	@echo ""
	@echo "$(GREEN)Usage:$(NC)"
	@echo "  make $(YELLOW)<target>$(NC)"
	@echo ""
	@echo "$(GREEN)Available targets:$(NC)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(YELLOW)%-20s$(NC) %s\n", $$1, $$2}'

# ============================================================================
# CLI Tool Wrapper - Main entry point
# ============================================================================

cli: ## Run CLI tool (usage: make cli gen module product)
	@go run cmd/cli/main.go $(filter-out $@,$(MAKECMDGOALS))

# Catch-all target to pass arguments to cli
%:
	@:

# ============================================================================
# Code Generation (Shortcuts)
# ============================================================================

gen: ## Generate code (usage: make gen module product)
	@go run cmd/cli/main.go gen $(filter-out $@,$(MAKECMDGOALS))

gen-module: ## Generate complete module (usage: make gen-module product)
	@go run cmd/cli/main.go gen module $(filter-out $@,$(MAKECMDGOALS))

gen-handler: ## Generate handler (usage: make gen-handler user)
	@go run cmd/cli/main.go gen handler $(filter-out $@,$(MAKECMDGOALS))

gen-repository: ## Generate repository (usage: make gen-repository order)
	@go run cmd/cli/main.go gen repository $(filter-out $@,$(MAKECMDGOALS))

gen-service: ## Generate service (usage: make gen-service auth)
	@go run cmd/cli/main.go gen service $(filter-out $@,$(MAKECMDGOALS))

gen-migration: ## Generate migration (usage: make gen-migration create_users)
	@go run cmd/cli/main.go gen migration $(filter-out $@,$(MAKECMDGOALS))

# ============================================================================
# Development
# ============================================================================

dev: ## Start development server with hot reload
	@go run cmd/cli/main.go dev

run: ## Run application without hot reload
	@go run cmd/cli/main.go run

build: ## Build API binary
	@go run cmd/cli/main.go build

build-all: ## Build all binaries (API, Seeder, CLI)
	@echo "$(BLUE)Building all binaries...$(NC)"
	@go build -o bin/api cmd/api/main.go
	@go build -o bin/seeder cmd/seeder/main.go
	@go build -o bin/cli cmd/cli/main.go
	@echo "$(GREEN)✓ All binaries built successfully$(NC)"

# ============================================================================
# Testing & Code Quality
# ============================================================================

test: ## Run all tests
	@go run cmd/cli/main.go test

test-verbose: ## Run tests with verbose output
	@go run cmd/cli/main.go test -v

test-coverage: ## Run tests with coverage report
	@go run cmd/cli/main.go test coverage

lint: ## Run linter
	@go run cmd/cli/main.go lint

fmt: ## Format code
	@echo "$(BLUE)Formatting code (excluding templates)...$(NC)"
	@gofmt -w $$(find . -name "*.go" -not -path "./templates/*" -not -path "./.git/*" -not -path "./vendor/*")
	@echo "$(GREEN)✓ Code formatted$(NC)"

vet: ## Run go vet
	@go run cmd/cli/main.go vet

mocks: ## Generate test mocks
	@go run cmd/cli/main.go mocks

check: fmt vet lint test ## Run all code quality checks

# ============================================================================
# Database
# ============================================================================

migrate: ## Run migrations (usage: make migrate up/down/status)
	@go run cmd/cli/main.go migrate $(filter-out $@,$(MAKECMDGOALS))

migrate-up: ## Run all migrations
	@go run cmd/cli/main.go migrate up

migrate-down: ## Rollback all migrations
	@go run cmd/cli/main.go migrate down

migrate-status: ## Show migration status
	@go run cmd/cli/main.go migrate status

seed: ## Seed database with initial data
	@go run cmd/cli/main.go seed

# ============================================================================
# Docker
# ============================================================================

docker: ## Docker commands (usage: make docker up/down/logs)
	@go run cmd/cli/main.go docker $(filter-out $@,$(MAKECMDGOALS))

docker-up: ## Start Docker containers
	@go run cmd/cli/main.go docker up

docker-down: ## Stop Docker containers
	@go run cmd/cli/main.go docker down

docker-logs: ## Show Docker logs
	@go run cmd/cli/main.go docker logs

docker-rebuild: ## Rebuild Docker containers
	@go run cmd/cli/main.go docker rebuild

# ============================================================================
# Tools
# ============================================================================

install: ## Install development tools
	@go run cmd/cli/main.go install tools

deps: ## Download Go dependencies
	@echo "$(BLUE)Downloading dependencies...$(NC)"
	@go mod download
	@echo "$(GREEN)✓ Dependencies downloaded$(NC)"

tidy: ## Tidy Go modules
	@echo "$(BLUE)Tidying modules...$(NC)"
	@go mod tidy
	@echo "$(GREEN)✓ Modules tidied$(NC)"

# ============================================================================
# Cleanup
# ============================================================================

clean: ## Clean build artifacts
	@echo "$(BLUE)Cleaning build artifacts...$(NC)"
	@rm -rf bin/
	@rm -f app coverage.out coverage.html
	@echo "$(GREEN)✓ Cleanup complete$(NC)"

clean-all: clean ## Clean everything including dependencies
	@echo "$(BLUE)Cleaning vendor and cache...$(NC)"
	@go clean -cache -testcache -modcache
	@echo "$(GREEN)✓ Everything cleaned$(NC)"

# ============================================================================
# Quick Actions
# ============================================================================

start: docker-up migrate-up dev ## Quick start: docker + migrate + dev

stop: docker-down ## Stop all services

restart: stop start ## Restart all services

logs: docker-logs ## Show all logs

status: ## Show project status
	@echo "$(BLUE)Project Status:$(NC)"
	@echo ""
	@echo "$(YELLOW)Go Version:$(NC)"
	@go version
	@echo ""
	@echo "$(YELLOW)Docker Status:$(NC)"
	@docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}" 2>/dev/null || echo "Docker not running"
	@echo ""
	@echo "$(YELLOW)Git Status:$(NC)"
	@git status -s || echo "Not a git repository"

# ============================================================================
# CI/CD Helpers
# ============================================================================

ci-test: ## Run tests for CI
	@go test -v -race -coverprofile=coverage.out -covermode=atomic ./...

ci-lint: ## Run linter for CI
	@golangci-lint run --timeout=5m

ci-build: ## Build for CI
	@go build -v -o bin/api cmd/api/main.go
	@go build -v -o bin/seeder cmd/seeder/main.go
	@go build -v -o bin/cli cmd/cli/main.go

ci: ci-lint ci-test ci-build ## Run all CI checks

