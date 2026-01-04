# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Makefile task runner with 50+ commands
  - Code generation commands (gen-module, gen-handler, gen-repository, gen-service, gen-migration)
  - Development commands (dev, run, build, build-all)
  - Testing commands (test, test-verbose, test-coverage, fmt, lint, vet, mocks, check)
  - Database commands (migrate-up, migrate-down, migrate-status, seed)
  - Docker commands (docker-up, docker-down, docker-logs, docker-rebuild)
  - Quick actions (start, stop, restart, status)
  - Tool installation command
  - CI/CD helpers (ci-lint, ci-test, ci-build, ci)
- CLI tool using Cobra framework (accessed via Makefile)
- Colored output in Makefile with help system
- Self-documenting commands via `make help`
- Production middleware implementation
  - Sanitize middleware for XSS protection (applied globally)
  - Timeout middleware with 30-second timeout (applied globally)
  - Rate limiter middleware for brute force protection (10 req/min on auth endpoints)

### Changed
- Primary interface now uses Makefile task runner
- All commands simplified: `make dev`, `make test`, `make gen-module product`
- Updated GitHub Actions CI/CD workflow to use Makefile
- Simplified Quick Start guide with Makefile commands
- Consolidated documentation structure
- Response package now uses pagination.Meta to avoid redundancy

### Removed
- Redundant documentation files
- Duplicate workflow guides
- RBAC middleware (redundant - role checking done at handler level)
- Metrics handler and middleware (unused, not implemented)
- Redundant Pagination struct from response package

## [1.0.0] - 2026-01-04

### Added
- CI/CD pipeline with GitHub Actions
  - Automated testing with PostgreSQL and Redis services
  - Code linting with golangci-lint
  - Multi-platform binary builds (Linux, macOS, Windows - AMD64/ARM64)
  - Docker image building and publishing to GHCR
  - Automated releases with semantic versioning
  - Code coverage reporting to Codecov

- Enhanced health check system
  - Comprehensive health check endpoint with DB and Redis status
  - Liveness probe endpoint (`/health/live`)
  - Readiness probe endpoint (`/health/ready`)
  - Uptime tracking
  - Detailed dependency status reporting
  - Service unavailable response (503) when unhealthy

- Database connection pooling configuration
  - Configurable max open connections
  - Configurable max idle connections
  - Configurable connection max lifetime
  - Configurable connection max idle time
  - Environment variable support for all pool settings

- Developer productivity tools
  - Feature templates for rapid development (entity, repository, usecase, handler, DTO, migration)
  - Mock generation setup with Mockery
  - Context helpers for request tracing
  - Template usage guide with examples
  - Makefile command for mock generation

- Interactive API documentation
  - Swagger UI at `/docs` endpoint
  - Beautiful interface for API exploration
  - Try-it-out functionality
  - CDN-hosted Swagger UI 5.x


### Changed
- **BREAKING**: Updated `health_handler.go` to require database and Redis client
  - Migration: Update handler initialization in `main.go` to pass `db` and `redisClient`
  
- **BREAKING**: Updated `postgres.go` to accept pool configuration
  - Migration: Update database initialization to pass `PoolConfig` struct

- Enhanced configuration validation
  - Added environment validation (development/staging/production)
  - Added database pool configuration validation
  - Added JWT secret strength checking
  - Added production-specific security checks
  - Better error messages for misconfigurations

- Updated router with new health check endpoints
  - Added `/health/live` endpoint
  - Added `/health/ready` endpoint
  - Added Swagger UI routes at `/docs` and `/docs/`

- Improved `.env.example` with database pool settings
  - Added `DB_MAX_OPEN_CONNS` (default: 25)
  - Added `DB_MAX_IDLE_CONNS` (default: 5)
  - Added `DB_CONN_MAX_LIFETIME` (default: 5m)
  - Added `DB_CONN_MAX_IDLE_TIME` (default: 10m)

- Enhanced Makefile
  - Added `mocks` target for mock generation
  - Updated `install-tools` to include Mockery
  - Updated help text with new commands

- Updated README with new features and documentation links

### Fixed
- None (this is the initial enhancement release)

### Deprecated
- None

### Removed
- None

### Security
- Enhanced JWT secret validation in production mode
- Added comprehensive configuration validation
- Database connection security improvements with configurable timeouts

---

## [1.0.0] - 2026-01-04

### Initial Release

#### Core Features
- Clean Architecture implementation
- Domain-Driven Design (DDD)
- JWT authentication with access and refresh tokens
- Role-Based Access Control (RBAC)
- User management (CRUD)
- Role management (CRUD)
- Soft delete support
- Multi-language support (English, Indonesian)

#### Infrastructure
- PostgreSQL 15 database
- Redis 7 for session management
- Docker and Docker Compose support
- Database migrations
- Database seeder

#### Security
- JWT token authentication
- Password hashing with bcrypt
- Session management via Redis
- CORS middleware
- Input validation
- Rate limiting
- Request timeout
- Input sanitization
- RBAC middleware

#### Development Tools
- Hot reload with Air
- OpenAPI/Swagger JSON generation
- golangci-lint configuration
- Makefile with comprehensive commands
- Structured logging with Zerolog
- Prometheus metrics

#### Testing
- Unit tests for pkg utilities
- Test coverage reporting
- Pagination tests
- Validator tests
- Auth tests
- Error handling tests

---

## Migration Guide

### From Initial Release to Unreleased

#### 1. Update Environment Variables

Add new database pool configuration to your `.env` file:

```bash
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=5
DB_CONN_MAX_LIFETIME=5m
DB_CONN_MAX_IDLE_TIME=10m
```

#### 2. Update Health Handler Initialization

If you've customized `main.go`, update the health handler initialization:

**Before:**
```go
healthHandler := handler.NewHealthHandler()
```

**After:**
```go
healthHandler := handler.NewHealthHandler(db, redisClient)
```

#### 3. Update Database Initialization

**Before:**
```go
db, err := database.NewPostgresDB(cfg.Database.DSN())
```

**After:**
```go
poolConfig := database.PoolConfig{
    MaxOpenConns:    cfg.Database.MaxOpenConns,
    MaxIdleConns:    cfg.Database.MaxIdleConns,
    ConnMaxLifetime: cfg.Database.ConnMaxLifetime,
    ConnMaxIdleTime: cfg.Database.ConnMaxIdleTime,
}
db, err := database.NewPostgresDB(cfg.Database.DSN(), poolConfig)
```

#### 4. Install New Tools

Run the updated install-tools command to get Mockery:

```bash
./app install tools
```

#### 5. Optional: Generate Mocks

If you want to use mocks for testing:

```bash
./app mocks
```

#### 6. Optional: Use Templates

Check out the new templates in the `templates/` directory for rapid feature development.

---

## Upcoming Features

### Planned for Next Release

- [ ] WebSocket support
- [ ] GraphQL API support
- [ ] gRPC API support
- [ ] Event sourcing
- [ ] CQRS pattern implementation
- [ ] Elasticsearch integration
- [ ] Background job processing
- [ ] Email service integration
- [ ] SMS service integration
- [ ] File upload handling
- [ ] Image processing
- [ ] PDF generation
- [ ] CSV import/export
- [ ] Audit logging
- [ ] Advanced RBAC with permissions
- [ ] Multi-tenancy support
- [ ] API versioning
- [ ] Request validation middleware
- [ ] Response caching
- [ ] Database sharding support

### Under Consideration

- [ ] OAuth2 provider support
- [ ] OpenTelemetry integration
- [ ] Service mesh support
- [ ] Message queue integration (RabbitMQ, Kafka)
- [ ] Blockchain integration
- [ ] Machine learning integration
- [ ] Serverless deployment support

---

## Contributors

Thank you to all contributors who helped make this boilerplate better!

---

## Support

- **Issues**: https://github.com/yourusername/go-rest-api-boilerplate/issues
- **Discussions**: https://github.com/yourusername/go-rest-api-boilerplate/discussions
- **Email**: support@yourdomain.com

---

[Unreleased]: https://github.com/yourusername/go-rest-api-boilerplate/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/yourusername/go-rest-api-boilerplate/releases/tag/v1.0.0

