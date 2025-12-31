.PHONY: help run build test clean docker-up docker-down migrate-up migrate-down migrate-create seed lint fmt vet

help:
	@echo "Available commands:"
	@echo "  make run           - Run the application"
	@echo "  make build         - Build the application"
	@echo "  make test          - Run tests"
	@echo "  make clean         - Clean build artifacts"
	@echo "  make docker-up     - Start Docker containers"
	@echo "  make docker-down   - Stop Docker containers"
	@echo "  make migrate-up    - Run database migrations"
	@echo "  make migrate-down  - Rollback database migrations"
	@echo "  make migrate-create NAME=migration_name - Create new migration"
	@echo "  make seed          - Run database seeder"
	@echo "  make lint          - Run linter"
	@echo "  make fmt           - Format code"
	@echo "  make vet           - Run go vet"

run:
	@echo "Starting application..."
	@go run cmd/api/main.go

build:
	@echo "Building application..."
	@go build -o bin/api cmd/api/main.go

test:
	@echo "Running tests..."
	@go test -v -cover ./...

clean:
	@echo "Cleaning up..."
	@rm -rf bin/
	@go clean

docker-up:
	@echo "Starting Docker containers..."
	@docker compose up -d

docker-down:
	@echo "Stopping Docker containers..."
	@docker compose down

migrate-up:
	@echo "Running migrations..."
	@docker exec -i go-rest-api-postgres psql -U postgres -d go-rest-api-boilerplate < migrations/20251231062946_create_table_users.up.sql
	@echo "Migrations completed!"

migrate-down:
	@echo "Rolling back migrations..."
	@docker exec -i go-rest-api-postgres psql -U postgres -d go-rest-api-boilerplate < migrations/20251231062946_create_table_users.down.sql
	@echo "Rollback completed!"

migrate-create:
	@if [ -z "$(NAME)" ]; then \
		echo "Error: NAME is required. Usage: make migrate-create NAME=migration_name"; \
		exit 1; \
	fi; \
	TIMESTAMP=$$(date +%s); \
	touch migrations/$${TIMESTAMP}_$(NAME).up.sql; \
	touch migrations/$${TIMESTAMP}_$(NAME).down.sql; \
	echo "Created: migrations/$${TIMESTAMP}_$(NAME).up.sql"; \
	echo "Created: migrations/$${TIMESTAMP}_$(NAME).down.sql"

seed:
	@echo "Running database seeder..."
	@go run cmd/seeder/main.go

lint:
	@echo "Running linter..."
	@golangci-lint run

fmt:
	@echo "Formatting code..."
	@go fmt ./...

vet:
	@echo "Running go vet..."
	@go vet ./...

