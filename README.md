# Go REST API Boilerplate

Production-ready REST API boilerplate built with Go, following Clean Architecture and industry best practices.

## Features

### Core
- Clean Architecture with dependency inversion
- JWT-based authentication with Redis sessions
- User registration & login
- Password hashing with bcrypt
- RESTful API design
- Health check endpoint

### Database
- PostgreSQL for persistent data
- Redis for session storage & caching
- Migration system with versioning
- Database seeder for development

### Security
- JWT token authentication
- Session management via Redis
- Password hashing (bcrypt)
- Session invalidation on logout
- CORS middleware
- Input validation
- SQL injection prevention

### Development
- Hot reload ready
- Docker Compose for local development
- Make commands for all operations
- Linter configuration (golangci-lint)
- Environment-based configuration

### Production Ready
- Graceful shutdown
- Connection pooling
- Structured logging middleware
- Recovery middleware (panic handling)
- Docker support

## Tech Stack

- **Go 1.21+** - Programming language
- **Gin** - HTTP web framework
- **PostgreSQL 15** - Primary database
- **Redis 7** - Session storage & caching
- **Docker Compose** - Container orchestration
- **JWT** - Authentication tokens

## Quick Start

```bash
# 1. Setup environment
cp .env.example .env

# 2. Start services
make docker-up

# 3. Run migrations
make migrate-up

# 4. Seed database (optional)
make seed

# 5. Run application
make run

# 6. Test health endpoint
curl http://localhost:8080/api/v1/health
```

## Using This Boilerplate

### For New Project

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

### Project Structure

```
cmd/
  api/              Application entry point
  seeder/           Database seeder
config/             Configuration management
internal/
  domain/           Business entities & repository interfaces
  usecase/          Business logic implementation
  delivery/http/    HTTP handlers, routes & middleware
  repository/       Data access layer (by domain)
  infrastructure/   Database connections (PostgreSQL, Redis)
  middleware/       HTTP middleware (auth, cors, logger, recovery)
pkg/                Public reusable packages
migrations/         SQL migration files
```


## Commands

```bash
# Application
make run           # Run application
make build         # Build binary
make test          # Run tests
make clean         # Clean artifacts

# Docker
make docker-up     # Start Docker services
make docker-down   # Stop Docker services

# Database
make migrate-up    # Run migrations
make migrate-down  # Rollback migrations
make seed          # Run database seeder

# Code Quality
make fmt           # Format code
make vet           # Vet code
make lint          # Run linter
```

## API Endpoints

### Public
- `GET  /api/v1/health` - Health check
- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/login` - User login

### Protected (Requires JWT)
- `POST /api/v1/auth/logout` - User logout
- `GET  /api/v1/users/me` - Get current user profile
- `GET  /api/v1/users/:id` - Get user by ID

## Documentation

- [cmd/](cmd/README.md) - Application entry points
- [internal/](internal/README.md) - Application code
- [internal/domain/](internal/domain/README.md) - Domain entities and errors
- [internal/middleware/](internal/middleware/README.md) - HTTP middleware
- [pkg/](pkg/README.md) - Public packages
- [pkg/auth/](pkg/auth/README.md) - Authentication utilities
- [pkg/response/](pkg/response/README.md) - Response formatting
- [pkg/validator/](pkg/validator/README.md) - Input validation
- [migrations/](migrations/README.md) - Database migrations

## Environment Variables

```bash
# Application
APP_PORT=8080

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
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

**Flow:**
1. **Login**: Generate JWT token, store in Redis with user ID
2. **Request**: Validate token from Redis, extract user ID
3. **Logout**: Delete token from Redis (immediate invalidation)
4. **Expire**: Auto-expire based on JWT_EXPIRATION

**Why Redis Only:**
- Fast in-memory access
- Built-in TTL (time-to-live)
- Automatic expiration
- Easier horizontal scaling
- No need for cleanup jobs

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for:
- Development workflow
- Branch naming conventions
- Commit message format
- Code style guidelines
- Testing requirements

## License

MIT License - see [LICENSE](LICENSE) file

## Support

- **Issues**: [GitHub Issues](https://github.com/username/repo/issues)
- **Discussions**: [GitHub Discussions](https://github.com/username/repo/discussions)

