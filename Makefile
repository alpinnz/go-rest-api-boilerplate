.PHONY: help run build test clean docker-up docker-down docker-status migrate-up migrate-down migrate-create migrate-status seed lint fmt vet swagger-generate

help:
	@echo "Available commands:"
	@echo "  make run              - Run the application"
	@echo "  make build            - Build the application"
	@echo "  make test             - Run tests"
	@echo "  make clean            - Clean build artifacts"
	@echo "  make docker-up        - Start Docker containers"
	@echo "  make docker-down      - Stop Docker containers"
	@echo "  make docker-status    - Show Docker containers status"
	@echo "  make migrate-up       - Run database migrations"
	@echo "  make migrate-down     - Rollback database migrations"
	@echo "  make migrate-create NAME=migration_name - Create new migration"
	@echo "  make migrate-status   - Show migration files"
	@echo "  make seed             - Run database seeder"
	@echo "  make swagger-generate - Generate OpenAPI specification"
	@echo "  make lint             - Run linter"
	@echo "  make fmt              - Format code"
	@echo "  make vet              - Run go vet"

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

docker-status:
	@echo "Docker containers status:"
	@docker compose ps
	@echo ""
	@echo "PostgreSQL container:"
	@CONTAINER=$$(docker ps --format '{{.Names}}' | grep -i postgres | head -n 1); \
	if [ -z "$$CONTAINER" ]; then \
		echo "  Status: NOT RUNNING"; \
	else \
		echo "  Name: $$CONTAINER"; \
		echo "  Status: RUNNING"; \
	fi

migrate-up:
	@echo "Running migrations..."
	@CONTAINER=$$(docker ps --format '{{.Names}}' | grep -i postgres | head -n 1); \
	if [ -z "$$CONTAINER" ]; then \
		echo "Error: PostgreSQL container not running. Run 'make docker-up' first."; \
		exit 1; \
	fi; \
	echo "Using container: $$CONTAINER"; \
	for file in migrations/*.up.sql; do \
		echo "Applying: $$file"; \
		docker exec -i $$CONTAINER psql -U postgres -d go-rest-api-boilerplate < $$file || exit 1; \
	done; \
	echo "Migrations completed!"

migrate-down:
	@echo "Rolling back migrations..."
	@CONTAINER=$$(docker ps --format '{{.Names}}' | grep -i postgres | head -n 1); \
	if [ -z "$$CONTAINER" ]; then \
		echo "Error: PostgreSQL container not running. Run 'make docker-up' first."; \
		exit 1; \
	fi; \
	echo "Using container: $$CONTAINER"; \
	for file in $$(ls -r migrations/*.down.sql); do \
		echo "Rolling back: $$file"; \
		docker exec -i $$CONTAINER psql -U postgres -d go-rest-api-boilerplate < $$file || exit 1; \
	done; \
	echo "Rollback completed!"

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

migrate-status:
	@echo "Migration files:"
	@echo ""
	@echo "UP migrations:"
	@ls -1 migrations/*.up.sql 2>/dev/null || echo "  No migration files found"
	@echo ""
	@echo "DOWN migrations:"
	@ls -1 migrations/*.down.sql 2>/dev/null || echo "  No migration files found"

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
