# Go REST API Boilerplate

Production-ready REST API boilerplate built with Go, following Clean Architecture, SOLID principles, and Domain-Driven Design.

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
- Graceful shutdown
- Connection pooling
- Structured logging middleware
- Recovery middleware (panic handling)
- Docker support
- Comprehensive error handling

## Prerequisites

Before running the application, ensure you have:

- **Go 1.21+** installed
- **PostgreSQL** server running (local or cloud)
- **Redis** server running (local or cloud)
- **PostgreSQL client (psql)** for migrations:
  - macOS: `brew install postgresql`
  - Ubuntu: `sudo apt-get install postgresql-client`
  - Windows: Download from [postgresql.org](https://www.postgresql.org/download/)
  - Or use [golang-migrate](https://github.com/golang-migrate/migrate) as alternative

## Quick Start

```bash
# 0. Install PostgreSQL client if not available (for migrations)
brew install postgresql  # macOS
# sudo apt-get install postgresql-client  # Ubuntu

# 1. Install development tools (Air for hot reload, Swag, etc.)
make install-tools

# 2. Setup environment
cp .env.example .env
# Edit .env with your PostgreSQL and Redis credentials

# 3. Run migrations (requires psql command)
make migrate-up

# 4. Seed database (optional)
make seed

# 5. Run with hot reload (recommended)
make dev

# Or run without hot reload
make run

# 6. Test endpoints
curl http://localhost:8080/api/v1/health
curl http://localhost:8080/docs/swagger.json
```

Note: Ensure your PostgreSQL and Redis servers are running before starting the application.

## Development

### Hot Reload

This boilerplate uses Air for automatic rebuild on file changes.

Setup:
```bash
# Install Air and other development tools
make install-tools

# Ensure GOPATH/bin is in your PATH
export PATH=$PATH:$(go env GOPATH)/bin
```

Usage:
```bash
# Run with hot reload (recommended for development)
make dev

# The application will automatically rebuild when you save changes
```

Configuration:
- Configuration file: `.air.toml`
- Auto rebuild on `.go` file changes
- Excludes test files and vendor directory
- Clear screen on rebuild

Benefits:
- Faster development cycle
- No manual restart needed
- Automatic rebuild on save
- Clean console output

### Using This Boilerplate

For New Project:
```bash
# Clone and setup
git clone <this-repo> my-new-project
cd my-new-project
rm -rf .git
git init

# Update module name in go.mod
# Update container names in docker-compose.yml
# Update import paths in all Go files
```

### Project Structure (SOLID + DDD)

```
cmd/
  api/              Application entry point
  seeder/           Database seeder
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
  middleware/       HTTP middleware (auth, cors, logger, recovery)
  localization/     Multi-language support
pkg/                Public reusable packages
  auth/             JWT and password utilities
  response/         HTTP response formatting
  validator/        Struct validation
migrations/         SQL migration files
docs/               API documentation (OpenAPI/Swagger)
```

Design Decisions:
- **Entity vs Model**: Using DDD terminology - "entity" for objects with identity
- **Separated Concerns**: Entity in `domain/entity/`, interfaces in `domain/repository/`
- **SOLID Principles**: Single responsibility, dependency inversion fully applied
- **Soft Delete**: Both User and Role support soft delete with `deleted_at`
- **Name Structure**: User has `FirstName` + `LastName` for flexibility
- **Domain Methods**: Entities include behavior (e.g., `FullName()`, `IsDeleted()`)

## Commands

```bash
# Development
make dev              # Run with hot reload (recommended)
make run              # Run application
make build            # Build API binary
make build-seeder     # Build seeder binary
make test             # Run tests
make test-coverage    # Run tests with coverage report
make clean            # Clean artifacts
make all              # Run fmt, vet, lint, test, and build

# Tools
make install-tools    # Install Air, Swag, golangci-lint

# Docker
make docker-up        # Start API Docker container
make docker-down      # Stop API Docker container
make docker-rebuild   # Rebuild and restart API container
make docker-logs      # Show API container logs (follow mode)
make docker-status    # Show API Docker container status

# Database
make migrate-up       # Run migrations (uses .env config)
make migrate-down     # Rollback migrations (uses .env config)
make migrate-create NAME=name # Create new migration
make migrate-status   # Show migration files
make seed             # Run database seeder

# Code Quality
make fmt              # Format code
make vet              # Run go vet
make lint             # Run linter
make swagger-generate # Generate OpenAPI specification
```

## API Endpoints

### Public
- `GET  /api/v1/health` - Health check
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

### Supported Languages
- English (en)
- Indonesian (id)

### Usage

```bash
# English response
curl -H "Accept-Language: en" http://localhost:8080/api/v1/auth/login

# Indonesian response
curl -H "Accept-Language: id" http://localhost:8080/api/v1/auth/login
```

### Adding New Language

1. Create new JSON file in `internal/localization/lang/`:
```bash
cp internal/localization/lang/en.json internal/localization/lang/fr.json
```

2. Translate all keys in the new file

3. Load in router (automatic via middleware)

## Environment Variables

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
JWT_SECRET=change-this-secret-key
JWT_EXPIRATION=24h

# Server
READ_TIMEOUT=10s
WRITE_TIMEOUT=10s
SHUTDOWN_TIMEOUT=5s
```

## Session Management

Sessions are stored in Redis for fast access and automatic expiration.

Flow:
1. Login: Generate JWT token, store in Redis with user ID
2. Request: Validate token from Redis, extract user ID
3. Logout: Delete token from Redis (immediate invalidation)
4. Expire: Auto-expire based on JWT_EXPIRATION

Why Redis Only:
- Fast in-memory access
- Built-in TTL (time-to-live)
- Automatic expiration
- Easier horizontal scaling
- No need for cleanup jobs

## Docker

This boilerplate only Dockerizes the API application. PostgreSQL and Redis should run externally (locally or cloud).

### Build and Run API Container

```bash
# Start API container
make docker-up

# Stop API container
make docker-down

# Check status
make docker-status
```

### Docker Configuration

The API container:
- Exposes port 8080
- Uses environment variables from .env file
- Connects to external PostgreSQL and Redis
- Includes all application dependencies

Note: Ensure PostgreSQL and Redis are accessible from Docker container (use host.docker.internal on macOS/Windows or adjust DB_HOST/REDIS_HOST).

## Migrations

Database migrations connect directly to PostgreSQL using .env configuration.

### Commands

```bash
# Run all migrations
make migrate-up

# Rollback all migrations
make migrate-down

# Create new migration
make migrate-create NAME=add_user_avatar

# View migration files
make migrate-status

# Check current migration version
make 
```

### Troubleshooting

**Error: `psql: command not found`**

Simply run:
```bash
make install-tools
```

This will automatically install PostgreSQL client along with other development tools.

### Current Schema

Users Table:
- id, email, password
- first_name, last_name
- created_at, updated_at, deleted_at (soft delete)

Roles Table:
- id, name, description
- created_at, updated_at, deleted_at (soft delete)

User Roles Table (Junction):
- user_id, role_id
- created_at

## Database Seeding

Database seeder populates initial data for development and testing.

### Commands

```bash
# Run seeder (truncates tables and seeds fresh data)
make seed

# Or directly
go run cmd/seeder/main.go
```

### Default Seeded Data

**Roles:**
- admin - Administrator with full access
- user - Regular user with standard access
- moderator - Moderator with elevated access

**Users with Role Assignments:**
- **admin@example.com** / !Password123
  - Name: Admin User
  - Roles: admin, user
  - Full administrative access

- **user@example.com** / !Password123
  - Name: Regular User
  - Roles: user
  - Standard user access

- **test@example.com** / !Password123
  - Name: Test User
  - Roles: user, moderator
  - Standard and moderator access

### Seeder Architecture

The seeder follows SOLID principles with separated concerns:

```
internal/seeder/
├── types.go              - Interfaces and data types
├── runner.go             - Orchestration logic
├── role_seeder.go        - Role seeding
├── user_seeder.go        - User seeding
└── user_role_seeder.go   - User-role assignments
```

**Seeding Order:**
1. Roles (admin, user, moderator)
2. Users (admin, user, test accounts)
3. User-Role Assignments (automatically assigned)

**Features:**
- Fresh mode: Truncates tables before seeding
- Idempotent: Safe to run multiple times
- Extensible: Easy to add new seeders
- Follows SOLID principles

See [internal/seeder/README.md](internal/seeder/README.md) for detailed documentation.

## Code Structure

### Domain Layer (`internal/domain/`)

Core business logic. No dependencies on external frameworks.

entity/:
- User entity with FirstName and LastName
- Role entity with soft delete
- Domain methods (FullName(), IsDeleted())

repository/:
- UserRepository interface
- RoleRepository interface
- SessionRepository interface

### Repository Layer (`internal/repository/`)

Data access implementations.

Features:
- Auto-exclude soft deleted records
- Context support for timeout/cancellation
- Error handling with domain errors
- SQL injection prevention

### UseCase Layer (`internal/usecase/`)

Business logic orchestration.

auth_usecase.go:
- Login, Register, Logout
- Refresh token
- Session management

user_usecase.go:
- User CRUD operations
- Role assignment

role_usecase.go:
- Role CRUD operations

### Delivery Layer (`internal/delivery/http/`)

HTTP request/response handling.

Structure:
- dto/ - Data Transfer Objects
- handler/ - HTTP handlers
- router/ - Route definitions

### Middleware Layer (`internal/middleware/`)

Cross-cutting concerns:
- auth.go - JWT authentication
- cors.go - CORS configuration
- language.go - Language detection
- logger.go - Request logging
- recovery.go - Panic recovery

## Contributing

### Development Workflow

```bash
# 1. Fork and clone repository
git clone https://github.com/YOUR_USERNAME/go-rest-api-boilerplate.git

# 2. Create feature branch
git checkout -b feat/new-feature

# 3. Make changes and test
make test

# 4. Format and vet code
make fmt && make vet

# 5. Commit with conventional commit format
git commit -m "feat(scope): add new feature"

# 6. Push and create pull request
git push origin feat/new-feature
```

### Branch Naming

Format: `<type>/<description>`

Examples:
- `feat/add-user-profile` - New feature
- `fix/jwt-expiration-bug` - Bug fix
- `docs/update-readme` - Documentation
- `refactor/user-repository` - Code refactoring
- `test/auth-middleware` - Test additions

### Commit Convention

Follow Conventional Commits specification:

Format:
```
<type>(<scope>): <subject>
```

Types:
- `feat` - New feature
- `fix` - Bug fix
- `docs` - Documentation
- `style` - Code style/format
- `refactor` - Code refactoring
- `test` - Add/update tests
- `chore` - Maintenance

Examples:
```bash
feat(auth): add JWT refresh token
fix(user): handle null email
docs: update API endpoints
refactor(repo): simplify query
test(auth): add login tests
```

### Code Style

1. Format: Run `make fmt` before commit
2. Vet: Run `make vet` and fix all issues
3. Lint: Run `make lint`
4. Naming: Use clear, descriptive names
5. Errors: Always handle errors explicitly
6. Context: Pass context.Context for cancellation
7. Comments: Document all exported functions
8. Tests: Add tests for new features

## License

MIT License - see LICENSE file

## Support

- Issues: GitHub Issues
- Discussions: GitHub Discussions

