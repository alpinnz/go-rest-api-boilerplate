# Go REST API Boilerplate

Production-ready REST API boilerplate built with Go, following Clean Architecture, SOLID principles, and Domain-Driven Design.

## Table of Contents

- [Tech Stack](#tech-stack)
- [Features](#features)
- [Prerequisites](#prerequisites)
- [Quick Start](#quick-start)
- [Project Structure](#project-structure)
- [Development Guide](#development-guide)
- [API Endpoints](#api-endpoints)
- [Authentication](#authentication)
- [Testing](#testing)
- [Configuration](#configuration)
- [Docker Deployment](#docker-deployment)
- [Database](#database)
- [Best Practices](#best-practices)
- [Contributing](#contributing)
- [License](#license)

## Tech Stack

- **Go 1.21+** - Programming language
- **Gin** - HTTP web framework
- **go-playground/validator** - Struct validation
- **PostgreSQL 15** - Primary database
- **Redis 7** - Session storage and caching
- **Docker** - API containerization
- **JWT** - Authentication tokens
- **Air** - Hot reload for development
- **golangci-lint** - Code linting
- **Swag** - OpenAPI/Swagger generation

## Features

### New Features (2026)
- **WebSocket Server** - Real-time bidirectional communication with Hub pattern and user targeting
- **CI/CD Pipeline** - Automated testing, linting, and releases with GitHub Actions
- **Enhanced Health Checks** - Comprehensive checks with DB, Redis status, and liveness/readiness probes
- **Developer Templates** - Rapid feature development with pre-built templates
- **Mock Generation** - Automated mock creation for testing with Mockery
- **Swagger UI** - Interactive API documentation at `/docs`
- **Connection Pool Config** - Configurable database connection pooling
- **Context Helpers** - Request tracing across all layers

### Core Features
- Clean Architecture with SOLID principles
- Domain-Driven Design (DDD) with entity terminology
- JWT-based authentication with Redis sessions
- Role-based access control (RBAC) with soft delete
- User registration and login with full CRUD
- Role management with full CRUD
- Many-to-many user-role relationship
- Soft delete for users and roles
- Password hashing with bcrypt
- RESTful API design
- Health check endpoint

### Database Features
- PostgreSQL for persistent data
- Redis for session storage and caching
- Migration system with versioning
- Database seeder for development
- Soft delete support
- Connection pooling

### Security Features
- JWT token authentication with access and refresh tokens
- Session management via Redis
- Password hashing with bcrypt
- Session invalidation on logout
- CORS middleware
- Input validation with go-playground/validator
- SQL injection prevention
- Token rotation for refresh tokens

### Localization
- Multi-language support (English, Indonesian)
- Accept-Language header detection
- Code-based error translation
- Parameter interpolation
- Extensible language system

### Development Features
- Hot reload with Air (automatic rebuild on file changes)
- Docker support for API containerization
- Make commands for all operations
- Linter configuration (golangci-lint)
- Environment-based configuration
- OpenAPI/Swagger JSON auto-generation

### Production Ready
- Graceful shutdown with context cancellation
- Configurable connection pooling (max open, idle, lifetime)
- Comprehensive health checks (DB, Redis, uptime)
- Liveness and readiness probes support
- Structured logging middleware
- Recovery middleware (panic handling)
- Docker deployment ready
- CI/CD pipeline with automated testing
- Comprehensive error handling
- Rate limiting and timeout middleware
- RBAC and security middleware

## Prerequisites

Before running the application, ensure you have:

- **Go 1.21+** installed
- **PostgreSQL** server running (local or cloud)
- **Redis** server running (local or cloud)
## Quick Start

```bash
# Install development tools
make install

# Setup environment
cp .env.example .env
# Edit .env with your configuration

# Start everything (Docker + Migrations + Dev server)
make start

# Test health endpoint
curl http://localhost:8080/api/v1/health
```

Use `make help` to see all available commands.


## Project Structure

```
cmd/
  api/              Application entry point
  seeder/           Database seeder
  cli/              CLI tool (accessed via Makefile)
    commands/       CLI commands
      generator/    Code generator utilities
        templates/  Code generation templates
config/             Configuration management
internal/
  domain/
    entity/         Domain entities (User, Role) - DDD terminology
    repository/     Repository interfaces (contracts)
    errors.go       Domain errors
  usecase/          Business logic (auth, user, role)
  delivery/http/    HTTP handlers, routes and middleware
  repository/       Repository implementations
  infrastructure/   Database connections (PostgreSQL, Redis)
  middleware/       HTTP middleware (auth, cors, logger, recovery, sanitize, timeout, rate_limiter)
  localization/     Multi-language support
pkg/                Public reusable packages
  auth/             JWT and password utilities
  response/         HTTP response formatting
  validator/        Struct validation
migrations/         SQL migration files
```

Design Decisions:
- **Entity vs Model**: Using DDD terminology - "entity" for objects with identity
- **Separated Concerns**: Entity in `domain/entity/`, interfaces in `domain/repository/`
- **SOLID Principles**: Single responsibility, dependency inversion fully applied
- **Soft Delete**: Both User and Role support soft delete with `deleted_at`
- **Name Structure**: User has `FirstName` + `LastName` for flexibility
- **Domain Methods**: Entities include behavior (e.g., `FullName()`, `IsDeleted()`)

## Development Workflow

### Quick Start Everything
```bash
make start           # Start docker + migrations + dev server
make stop            # Stop all services
```

### Code Generation
```bash
make gen-module name=product      # Generate complete module (all layers)
make gen-handler name=user        # Generate handler only
make gen-repository name=order    # Generate repository
make gen-usecase name=auth        # Generate use case
make gen-migration name=create_products  # Generate migration files
```

### Running Application
```bash
make dev                         # Run with hot reload (recommended)
make run                         # Run without hot reload
make build                       # Build API binary
```

### Testing and Quality
```bash
make test                        # Run all tests
make test-coverage               # Generate HTML coverage report
make fmt                         # Format code
make lint                        # Run linter
make vet                         # Run go vet
make mocks                       # Generate mocks for testing
make check                       # Run all quality checks
```

### Database
```bash
make migrate-up                  # Run migrations
make migrate-down                # Rollback migrations
make migrate-status              # Show migration status
make seed                        # Seed database with initial data
```

### Docker
```bash
make docker-up                   # Start containers
make docker-down                 # Stop containers
make docker-logs                 # Show logs
make docker-restart              # Restart containers
```

### Tools
```bash
make install                     # Install development tools
make clean                       # Clean build artifacts
make help                        # Show all available commands
```

## API Endpoints

### Public
- `GET  /api/v1/health` - Comprehensive health check (DB, Redis, uptime)
- `GET  /api/v1/health/live` - Liveness probe
- `GET  /api/v1/health/ready` - Readiness probe
- `GET  /docs` - Interactive Swagger UI
- `GET  /docs/swagger.json` - OpenAPI specification

### Authentication
- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/login` - User login (returns access_token and refresh_token)
- `POST /api/v1/auth/refresh-token` - Refresh access token using refresh_token
- `POST /api/v1/auth/logout` - User logout (requires authentication)

### User Management (Protected)
- `GET  /api/v1/users/me` - Get current user profile
- `GET  /api/v1/users` - List all users (paginated)
- `GET  /api/v1/users/:id` - Get user by UUID (e.g., 01933e7f-1234-7890-abcd-ef0123456789)
- `PUT  /api/v1/users/:id` - Update user by UUID
- `DELETE /api/v1/users/:id` - Delete user by UUID (soft delete)
- `POST /api/v1/users/:id/roles` - Assign role to user by UUID
- `DELETE /api/v1/users/:id/roles/:roleId` - Remove role from user by UUID

### Role Management (Protected)
- `GET  /api/v1/roles` - List all roles (paginated)
- `POST /api/v1/roles` - Create new role
- `GET  /api/v1/roles/:id` - Get role by UUID
- `PUT  /api/v1/roles/:id` - Update role by UUID
- `DELETE /api/v1/roles/:id` - Delete role by UUID (soft delete)

### WebSocket (Protected)
- `GET  /api/v1/ws` - WebSocket connection upgrade (requires JWT authentication)
- `GET  /api/v1/ws/stats` - Get WebSocket connection statistics with detailed metrics

**WebSocket Documentation**: See [internal/websocket/README.md](internal/websocket/README.md) for detailed implementation, usage examples, and client integration guides.

## Authentication

This API uses JWT-based authentication with access and refresh tokens.

### Token Types

Access Token:
- Lifetime: 15 minutes
- Used for API requests
- Sent in `Authorization: Bearer <token>` header
- Short-lived to minimize risk if compromised

Refresh Token:
- Lifetime: 7 days
- Used to obtain new access tokens
- Sent in request body to `/auth/refresh` endpoint
- Single-use (token rotation for security)

### How It Works

1. Login: User provides credentials, receives both tokens
2. API Request: Client uses access token in Authorization header
3. Token Expired: When 401 received, use refresh token to get new pair
4. Logout: Both tokens invalidated immediately via Redis

### Login Response
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": 1,
    "email": "user@example.com",
    "first_name": "John",
    "last_name": "Doe"
  }
}
```

### Security Features

- Token Rotation: Refresh tokens are single-use, new pair generated each refresh
- Short-Lived Access: 15-minute expiration minimizes exposure
- Server-Side Sessions: All tokens stored in Redis, can be immediately revoked
- Token Type Validation: Access tokens can't be used as refresh tokens
- HTTPS Required: Always use HTTPS in production

### Usage Examples
```bash
# Use access token in Authorization header
curl -H "Authorization: Bearer <access_token>" \
     http://localhost:8080/api/v1/users/me

# Refresh access token when expired
curl -X POST http://localhost:8080/api/v1/auth/refresh \
     -H "Content-Type: application/json" \
     -d '{"refresh_token": "<refresh_token>"}'
```

### Troubleshooting

Access Token Expired (401):
- Use refresh token to get new access token

Refresh Token Expired (401):
- User must login again (7 days passed)

Both Tokens Invalid:
- User logged out or server restarted (Redis cleared)
- User must login again

## Localization

Multi-language support with automatic detection from `Accept-Language` header.

Supported Languages: English (en), Indonesian (id)

Usage:
```bash
# English response
curl -H "Accept-Language: en" http://localhost:8080/api/v1/auth/login

# Indonesian response
curl -H "Accept-Language: id" http://localhost:8080/api/v1/auth/login
```

Adding New Language:
1. Create new JSON file: `internal/localization/lang/fr.json`
2. Translate all keys
3. Middleware will automatically load it

## Testing

**Running Tests:**
```bash
make test                        # Run all tests
make test-coverage               # Generate HTML coverage report
make mocks                       # Generate test mocks
```

**Test Structure:**
- **Unit Tests:** `pkg/*/` (pagination, validator, auth)
- **Repository Tests:** `internal/repository/`
- **Use Case Tests:** `internal/usecase/` (with mocked repositories)
- **Handler Tests:** `internal/delivery/http/handler/`

**Coverage Goals:**
- Unit Tests: > 80%
- Integration Tests: Critical user flows
- Handler Tests: All endpoints

Writing Tests (table-driven pattern):
```go
func TestValidateEmail(t *testing.T) {
    tests := []struct {
        name    string
        email   string
        wantErr bool
    }{
        {"valid email", "test@example.com", false},
        {"invalid format", "invalid", true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateEmail(tt.email)
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```


## Configuration

Environment variables (`.env` file):

```bash
# Application
APP_NAME=go-rest-api-boilerplate
APP_ENV=development
APP_PORT=8080

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASS=postgres
DB_NAME=go-rest-api-boilerplate
DB_SSLMODE=disable

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# JWT
ACCESS_TOKEN_SECRET=change-this-access-token-secret-key
ACCESS_TOKEN_EXPIRATION=15m
REFRESH_TOKEN_SECRET=change-this-refresh-token-secret-key
REFRESH_TOKEN_EXPIRATION=168h

# Server
READ_TIMEOUT=10s
WRITE_TIMEOUT=10s
SHUTDOWN_TIMEOUT=5s
```

Session Management:
- Sessions stored in Redis for fast access
- Access tokens expire based on ACCESS_TOKEN_EXPIRATION (default: 15m)
- Refresh tokens expire based on REFRESH_TOKEN_EXPIRATION (default: 168h/7 days)
- Login: Generate access & refresh tokens, store in Redis
- Request: Validate access token from Redis, extract user ID
- Logout: Delete tokens from Redis (immediate invalidation)
- Refresh: Generate new token pair using valid refresh token

## Docker

The API container:
- Exposes port 8080
- Uses environment variables from .env file
- Connects to external PostgreSQL and Redis
- Includes all application dependencies

**Commands:**
```bash
make docker-up                   # Start API container
make docker-down                 # Stop API container
make docker-logs                 # View container logs
make docker-restart              # Restart container
```

Note: PostgreSQL and Redis run externally (locally or cloud). Use `host.docker.internal` on macOS/Windows for DB_HOST/REDIS_HOST.

## Database Migrations

Uses golang-migrate tool with .env configuration.

**Commands:**
```bash
make migrate-up                     # Run all migrations
make migrate-down                   # Rollback last migration
make migrate-status                 # Show migration status
make gen-migration name=add_user_avatar  # Create new migration
```

Current Schema:
- Users: id, email, password, first_name, last_name, timestamps, deleted_at
- Roles: id, name, description, timestamps, deleted_at
- User Roles: user_id, role_id, created_at (junction table)

## Database Seeding

Populates initial data for development and testing.

Command:
```bash
make seed                        # Run seeder
```

Default Seeded Data:

Roles:
- admin - Administrator with full access
- user - Regular user with standard access
- moderator - Moderator with elevated access

Users:
- admin@example.com / !Password123 (admin, user roles)
- user@example.com / !Password123 (user role)
- test@example.com / !Password123 (user, moderator roles)

Seeder Architecture:
- Fresh mode: Truncates tables before seeding
- Idempotent: Safe to run multiple times
- Follows SOLID principles
- Order: Roles → Users → User-Role Assignments

## Best Practices

Code Organization (Clean Architecture):
1. Domain Layer (`internal/domain/`) - Core business logic, no external dependencies
2. Use Case Layer (`internal/usecase/`) - Business logic orchestration
3. Repository Layer (`internal/repository/`) - Data access implementations
4. Delivery Layer (`internal/delivery/`) - HTTP handlers and DTOs

Repository Pattern:
```go
// Always use context
func (r *UserRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
    // Implementation
}

// Use transactions
err := database.WithTransaction(ctx, db, func(ctx context.Context, tx *sql.Tx) error {
    // Multiple operations
    return nil
})
```

Use Case Pattern:
```go
// Set timeout for long operations
ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
defer cancel()

// Handle errors with context
if err != nil {
    return errors.Wrap(err, "USER_FETCH_FAILED", "Failed to fetch user").
        WithContext("user_id", id)
}
```

Handler Pattern:
```go
// Validate input
validationErrors := validator.Validate(req)
if len(validationErrors) > 0 {
    response.ValidationError(c, validationErrors)
    return
}

// Use pagination for lists
params := pagination.FromContext(c)
users, total, err := h.useCase.List(ctx, params.Offset, params.PerPage)
meta := pagination.NewMeta(params.Page, params.PerPage, total)
response.SuccessWithMeta(c, users, &meta)
```

Security Best Practices:
- Never commit .env file
- Use strong, random JWT secret in production
- Enable rate limiting on public endpoints
- Validate all inputs with validator
- Use parameterized queries (handled by lib/pq)
- Sanitize inputs with sanitize middleware
- Configure CORS allowed origins properly
- Always use TLS in production

Error Handling:
```go
// Create specific errors
err := errors.NotFound("User not found")
err := errors.BadRequest("Invalid email format")
err := errors.Unauthorized("Token expired")

// Add context
err := errors.NotFound("User not found").
    WithContext("user_id", userID).
    WithContext("action", "fetch")
```

Logging:
```go
// Use structured logging
logger.WithFields(map[string]interface{}{
    "user_id": userID,
    "action": "login",
}).Info("User logged in")
```

## Contributing

### Development Workflow

**Setup:**
```bash
# Fork and clone repository
git clone https://github.com/YOUR_USERNAME/go-rest-api-boilerplate.git

# Create feature branch
git checkout -b feat/new-feature

# Make changes and test
make check

# Commit with conventional commit format
git commit -m "feat(scope): add new feature"

# Push and create pull request
git push origin feat/new-feature
```

### Branch Naming

Format: `<type>/<description>`

Examples:
- `feat/add-user-profile` - New feature
- `fix/jwt-expiration-bug` - Bug fix
- `docs/update-readme` - Documentation
- `refactor/cleanup-handlers` - Code refactoring
- `test/add-unit-tests` - Tests

### Commit Convention

Format: `<type>(<scope>): <subject>`

**Types:**
- `feat` - New feature
- `fix` - Bug fix
- `docs` - Documentation changes
- `style` - Code style changes (formatting, etc.)
- `refactor` - Code refactoring
- `test` - Adding or updating tests
- `chore` - Maintenance tasks

**Examples:**
```bash
feat(auth): add JWT refresh token support
fix(user): handle null email in validation
docs: update API endpoints documentation
refactor(handler): simplify error handling
test(repository): add user repository tests
```

### Code Quality

Before committing:
```bash
make check    # Runs fmt, vet, lint, and tests
```

**Requirements:**
- Run `make check` before commit
- Document all exported functions
- Add tests for new features
- Update documentation
- Follow existing code style

## License

MIT License - see LICENSE file


