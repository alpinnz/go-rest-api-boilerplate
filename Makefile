.PHONY: help dev run build test lint fmt vet mocks swagger migrate seed docker install clean

# Default target
.DEFAULT_GOAL := help

# Load environment variables from .env file
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

# Colors for output
BLUE := \033[0;34m
GREEN := \033[0;32m
YELLOW := \033[0;33m
RED := \033[0;31m
NC := \033[0m # No Color

help: ## Show this help message
	@echo "$(BLUE)Go REST API Boilerplate$(NC)"
	@echo ""
	@echo "$(GREEN)Usage:$(NC)"
	@echo "  make $(YELLOW)<target>$(NC)"
	@echo ""
	@echo "$(GREEN)Available targets:$(NC)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(YELLOW)%-20s$(NC) %s\n", $$1, $$2}'

# ============================================================================
# Code Generation
# ============================================================================

gen-module: ## Generate complete module (usage: make gen-module name=product)
	@go run cmd/cli/main.go gen module $(name)

gen-handler: ## Generate handler (usage: make gen-handler name=user)
	@go run cmd/cli/main.go gen handler $(name)

gen-repository: ## Generate repository (usage: make gen-repository name=order)
	@go run cmd/cli/main.go gen repository $(name)

gen-usecase: ## Generate usecase (usage: make gen-usecase name=auth)
	@go run cmd/cli/main.go gen usecase $(name)

gen-migration: ## Generate migration (usage: make gen-migration name=create_users)
	@go run cmd/cli/main.go gen migration $(name)

swagger: ## Generate Swagger/OpenAPI documentation
	@echo "$(BLUE)Generating Swagger documentation...$(NC)"
	@swag init -g cmd/api/main.go -o internal/delivery/http/docs --parseDependency --parseInternal || (echo "$(RED)Swag not installed. Run: make install$(NC)" && exit 1)
	@echo "$(GREEN)✓ Swagger documentation generated: internal/delivery/http/docs/swagger.json$(NC)"

# ============================================================================
# Development
# ============================================================================

dev: ## Start development server with hot reload
	@echo "$(BLUE)Starting development server with hot reload...$(NC)"
	@air || (echo "$(RED)Air not installed. Run: make install$(NC)" && exit 1)

run: ## Run application without hot reload
	@echo "$(BLUE)Running application...$(NC)"
	@go run cmd/api/main.go

build: ## Build API binary
	@echo "$(BLUE)Building API binary...$(NC)"
	@go build -o bin/api cmd/api/main.go
	@echo "$(GREEN)✓ API binary built: bin/api$(NC)"

# ============================================================================
# Testing & Code Quality
# ============================================================================

test: ## Run all tests
	@echo "$(BLUE)Running tests...$(NC)"
	@go test -v ./...

test-coverage: ## Run tests with coverage report
	@echo "$(BLUE)Running tests with coverage...$(NC)"
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)✓ Coverage report: coverage.html$(NC)"

fmt: ## Format code
	@echo "$(BLUE)Formatting code...$(NC)"
	@gofmt -w $$(find . -name "*.go" -not -path "*/templates/*" -not -path "*/.git/*" -not -path "*/vendor/*")
	@echo "$(GREEN)✓ Code formatted$(NC)"

lint: ## Run linter
	@echo "$(BLUE)Running linter...$(NC)"
	@golangci-lint run || (echo "$(RED)golangci-lint not installed. Run: make install$(NC)" && exit 1)

vet: ## Run go vet
	@echo "$(BLUE)Running go vet...$(NC)"
	@go vet ./...

mocks: ## Generate test mocks
	@echo "$(BLUE)Generating mocks...$(NC)"
	@mockery --all --dir internal/domain/repository --output internal/domain/repository/mocks || (echo "$(RED)Mockery not installed. Run: make install$(NC)" && exit 1)
	@echo "$(GREEN)✓ Mocks generated$(NC)"

check: fmt vet lint test ## Run all code quality checks

# ============================================================================
# Database
# ============================================================================

migrate-up: ## Run all migrations
	@echo "$(BLUE)Running migrations...$(NC)"
	@if [ -z "$${DB_HOST}" ]; then echo "$(RED)Error: DB_HOST not set in .env$(NC)" && exit 1; fi
	@migrate -path migrations -database "postgres://$${DB_USER}:$${DB_PASS}@$${DB_HOST}:$${DB_PORT}/$${DB_NAME}?sslmode=$${DB_SSLMODE}" up || (echo "$(RED)migrate not installed. Run: make install$(NC)" && exit 1)
	@echo "$(GREEN)✓ Migrations completed$(NC)"

migrate-down: ## Rollback last migration
	@echo "$(BLUE)Rolling back migration...$(NC)"
	@migrate -path migrations -database "postgres://$${DB_USER}:$${DB_PASS}@$${DB_HOST}:$${DB_PORT}/$${DB_NAME}?sslmode=$${DB_SSLMODE}" down 1
	@echo "$(GREEN)✓ Migration rolled back$(NC)"

migrate-status: ## Show migration status
	@echo "$(BLUE)Migration status:$(NC)"
	@migrate -path migrations -database "postgres://$${DB_USER}:$${DB_PASS}@$${DB_HOST}:$${DB_PORT}/$${DB_NAME}?sslmode=$${DB_SSLMODE}" version

seed: ## Seed database with initial data
	@echo "$(BLUE)Seeding database...$(NC)"
	@go run cmd/seeder/main.go
	@echo "$(GREEN)✓ Database seeded$(NC)"

# ============================================================================
# Docker
# ============================================================================

docker-up: ## Start Docker containers
	@echo "$(BLUE)Starting Docker containers...$(NC)"
	@docker-compose up -d
	@echo "$(GREEN)✓ Containers started$(NC)"

docker-down: ## Stop Docker containers
	@echo "$(BLUE)Stopping Docker containers...$(NC)"
	@docker-compose down
	@echo "$(GREEN)✓ Containers stopped$(NC)"

docker-logs: ## Show Docker logs
	@docker-compose logs -f

docker-restart: ## Restart Docker containers
	@echo "$(BLUE)Restarting Docker containers...$(NC)"
	@docker-compose restart
	@echo "$(GREEN)✓ Containers restarted$(NC)"

# ============================================================================
# Tools
# ============================================================================

install: ## Install development tools
	@echo "$(BLUE)Installing development tools...$(NC)"
	@echo "Installing Air (hot reload)..."
	@go install github.com/air-verse/air@latest
	@echo "Installing Swag (API docs)..."
	@go install github.com/swaggo/swag/cmd/swag@latest
	@echo "Installing Mockery (mocks)..."
	@go install github.com/vektra/mockery/v2@latest
	@echo "Installing golangci-lint..."
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin
	@echo "Installing golang-migrate..."
	@go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@echo "$(GREEN)✓ All tools installed successfully!$(NC)"

# ============================================================================
# Cleanup
# ============================================================================

clean: ## Clean build artifacts
	@echo "$(BLUE)Cleaning build artifacts...$(NC)"
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@echo "$(GREEN)✓ Cleanup complete$(NC)"

# ============================================================================
# Quick Actions
# ============================================================================

start: ## Quick start: docker + migrate + dev
	@echo "$(BLUE)Starting full stack...$(NC)"
	@make docker-up
	@sleep 3
	@make migrate-up
	@make dev

stop: docker-down ## Stop all services

