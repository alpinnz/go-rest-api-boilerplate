.PHONY: help run dev build build-seeder test test-coverage clean all docker-up docker-down docker-rebuild docker-logs docker-status migrate-up migrate-down migrate-create migrate-status migrate-force seed lint fmt vet swagger-generate install-tools

help:
	@echo "Available commands:"
	@echo ""
	@echo "Development:"
	@echo "  make dev              - Run with hot reload (recommended)"
	@echo "  make run              - Run application"
	@echo "  make build            - Build API binary"
	@echo "  make build-seeder     - Build seeder binary"
	@echo "  make test             - Run tests"
	@echo "  make test-coverage    - Run tests with coverage report"
	@echo "  make clean            - Clean build artifacts"
	@echo "  make all              - Run fmt, vet, lint, test, and build"
	@echo ""
	@echo "Docker:"
	@echo "  make docker-up        - Start API Docker container"
	@echo "  make docker-down      - Stop API Docker container"
	@echo "  make docker-rebuild   - Rebuild and restart API container"
	@echo "  make docker-logs      - Show API container logs"
	@echo "  make docker-status    - Show API Docker container status"
	@echo ""
	@echo "Database:"
	@echo "  make migrate-up       - Run database migrations"
	@echo "  make migrate-down     - Rollback database migrations"
	@echo "  make migrate-create NAME=name - Create new migration"
	@echo "  make migrate-status   - Show migration files"
	@echo "  make seed             - Run database seeder"
	@echo ""
	@echo "Code Quality:"
	@echo "  make fmt              - Format code"
	@echo "  make vet              - Run go vet"
	@echo "  make lint             - Run linter"
	@echo "  make swagger-generate - Generate OpenAPI specification"
	@echo ""
	@echo "Tools:"
	@echo "  make install-tools    - Install all development tools"

dev:
	@echo "Starting application with hot reload..."
	@which air > /dev/null || (echo "Air not installed. Run 'make install-tools' first." && exit 1)
	@air

run:
	@echo "Starting application..."
	@go run cmd/api/main.go

# Default target - run full pipeline
all: fmt vet lint test build
	@echo ""
	@echo "✓ All checks passed!"

install-tools:
	@echo "Installing development tools..."
	@echo ""
	@echo "1. Installing Air (hot reload)..."
	@go install github.com/air-verse/air@latest
	@echo "2. Installing Swag (OpenAPI generator)..."
	@go install github.com/swaggo/swag/cmd/swag@latest
	@echo "3. Installing golangci-lint..."
	@which golangci-lint > /dev/null || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin
	@echo "4. Installing golang-migrate..."
	@if command -v brew > /dev/null 2>&1; then \
		brew list golang-migrate > /dev/null 2>&1 || brew install golang-migrate; \
	else \
		go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest; \
	fi
	@echo ""
	@echo "✓ All tools installed successfully!"
	@echo "Run 'make dev' to start development with hot reload"

build:
	@echo "Building API..."
	@go build -o bin/api cmd/api/main.go
	@echo "✓ Binary created: bin/api"

build-seeder:
	@echo "Building seeder..."
	@go build -o bin/seeder cmd/seeder/main.go
	@echo "✓ Binary created: bin/seeder"

test:
	@echo "Running tests..."
	@go test -v -cover ./...

test-coverage:
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✓ Coverage report: coverage.html"
	@echo ""
	@go tool cover -func=coverage.out | grep total

clean:
	@echo "Cleaning up..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@go clean
	@echo "✓ Clean completed"

docker-up:
	@echo "Starting API Docker container..."
	@docker compose up -d api
	@echo "✓ Container started"
	@echo "Run 'make docker-logs' to view logs"

docker-down:
	@echo "Stopping API Docker container..."
	@docker compose down api
	@echo "✓ Container stopped"

docker-rebuild:
	@echo "Rebuilding API Docker container..."
	@docker compose down api
	@docker compose build api
	@docker compose up -d api
	@echo "✓ Container rebuilt and started"

docker-logs:
	@echo "Showing API container logs (Ctrl+C to exit)..."
	@docker compose logs -f api

docker-status:
	@echo "API Docker container status:"
	@docker compose ps api
	@echo ""
	@CONTAINER=$$(docker ps --format '{{.Names}}' | grep -i go-rest-api-boilerplate | head -n 1); \
	if [ -z "$$CONTAINER" ]; then \
		echo "  Status: NOT RUNNING"; \
	else \
		echo "  Name: $$CONTAINER"; \
		echo "  Status: RUNNING"; \
		docker inspect --format='  Ports: {{range .NetworkSettings.Ports}}{{.}}{{end}}' $$CONTAINER; \
	fi

migrate-up:
	@echo "Running migrations..."
	@if [ ! -f .env ]; then \
		echo "✗ Error: .env file not found"; \
		echo "  Run: cp .env.example .env"; \
		exit 1; \
	fi
	@if ! command -v migrate > /dev/null 2>&1; then \
		echo "✗ Error: golang-migrate not installed"; \
		echo "  Run: make install-tools"; \
		exit 1; \
	fi

	@set -a; . ./.env; set +a; \
	migrate -path migrations \
		-database "postgresql://$$DB_USER:$$DB_PASS@$$DB_HOST:$$DB_PORT/$$DB_NAME?sslmode=disable" \
		up && echo "✓ Up completed!" || (echo "✗ Up failed"; exit 1)

migrate-down:
	@echo "Rolling back migrations..."
	@if [ ! -f .env ]; then \
		echo "✗ Error: .env file not found"; \
		echo "  Run: cp .env.example .env"; \
		exit 1; \
	fi
	@if ! command -v migrate > /dev/null 2>&1; then \
		echo "✗ Error: golang-migrate not installed"; \
		echo "  Run: make install-tools"; \
		exit 1; \
	fi
	@set -a; . ./.env; set +a; \
	migrate -path migrations \
		-database "postgresql://$$DB_USER:$$DB_PASS@$$DB_HOST:$$DB_PORT/$$DB_NAME?sslmode=disable" \
		down && echo "✓ Drop completed!" || (echo "✗ Drop failed"; exit 1)

migrate-create:
	@if [ -z "$(NAME)" ]; then \
		echo "✗ Error: NAME is required"; \
		echo "  Usage: make migrate-create NAME=migration_name"; \
		exit 1; \
	fi
	@TIMESTAMP=$$(date +%Y%m%d%H%M%S); \
	touch migrations/$${TIMESTAMP}_$(NAME).up.sql; \
	touch migrations/$${TIMESTAMP}_$(NAME).down.sql; \
	echo "✓ Migration files created:"; \
	echo "  - migrations/$${TIMESTAMP}_$(NAME).up.sql"; \
	echo "  - migrations/$${TIMESTAMP}_$(NAME).down.sql"

migrate-status:
	@echo "Migration files:"
	@echo ""
	@ls -1 migrations/*.up.sql 2>/dev/null | sed 's|migrations/||' | sed 's|.up.sql||' || echo "No migration files found"

seed: migrate-up
	@echo "Running database seeder..."
	@go run cmd/seeder/main.go

lint:
	@echo "Running linter..."
	@golangci-lint run
	@echo "✓ Linting completed"

fmt:
	@echo "Formatting code..."
	@go fmt ./...
	@echo "✓ Code formatted"

vet:
	@echo "Running go vet..."
	@go vet ./...
	@echo "✓ Vet completed"

swagger-generate:
	@echo "Generating OpenAPI specification from code annotations..."
	@go run github.com/swaggo/swag/cmd/swag@latest init \
		-g cmd/api/main.go \
		--output docs \
		--outputTypes json \
		--parseDependency \
		--parseInternal 2>&1 | grep -v "warning: failed to get package name" | grep -v "warning: failed to evaluate const" || true
	@rm -f docs/docs.go docs/swagger.yaml
	@echo "✓ Generated: docs/swagger.json"
	@echo ""
	@echo "To view documentation:"
	@echo "  https://editor.swagger.io (import docs/swagger.json)"
	@echo ""
	@echo "To validate:"
	@echo "  npx @apidevtools/swagger-cli validate docs/swagger.json"

migrate-force:
	@if [ -z "$(VERSION)" ]; then \
		echo "✗ Error: VERSION is required"; \
		echo "  Usage: make migrate-force VERSION=20251231062946"; \
		exit 1; \
	fi
	@if [ ! -f .env ]; then \
		echo "✗ Error: .env file not found"; \
		echo "  Run: cp .env.example .env"; \
		exit 1; \
	fi
	@if ! command -v migrate > /dev/null 2>&1; then \
		echo "✗ Error: golang-migrate not installed"; \
		echo "  Run: make install-tools"; \
		exit 1; \
	fi
	@set -a; . ./.env; set +a; \
	echo "Forcing migration version to $(VERSION)..."; \
	migrate -path migrations \
		-database "postgresql://$$DB_USER:$$DB_PASS@$$DB_HOST:$$DB_PORT/$$DB_NAME?sslmode=disable" \
		force $(VERSION) && echo "✓ Migration version forced to $(VERSION)"


