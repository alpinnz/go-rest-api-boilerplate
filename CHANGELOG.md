# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- WebSocket server for real-time bidirectional communication
  - Hub pattern for connection management
  - JWT authentication required for WebSocket connections
  - Message broadcasting to all clients
  - Targeted message sending to specific users
  - Structured message types (notification, system, chat)
  - Ping/Pong health checks (54 second interval)
  - Graceful connection handling with automatic cleanup
  - Thread-safe operations with sync.RWMutex
  - Detailed connection statistics endpoint with 6 metrics
  - Per-user connection tracking
  - Queue health monitoring
  - Top 10 users by connection count
  - WebSocket use case layer for business logic
  - Complete documentation in internal/websocket/README.md
  - Unit tests for WebSocket hub functionality
- Separate JWT configuration for access and refresh tokens
  - ACCESS_TOKEN_SECRET and ACCESS_TOKEN_EXPIRATION (default: 15m)
  - REFRESH_TOKEN_SECRET and REFRESH_TOKEN_EXPIRATION (default: 168h/7 days)
- GetAccessTokenExpiration() and GetRefreshTokenExpiration() helper functions
- SetJWTConfig() function to configure both token types
- Enhanced JWT validation with automatic secret selection based on token type
- Makefile automation with essential commands
  - Code generation commands (gen-module, gen-handler, gen-repository, gen-usecase, gen-migration)
  - Documentation generation (swagger)
  - Development commands (dev, run, build)
  - Testing commands (test, test-coverage, fmt, lint, vet, mocks, check)
  - Database commands (migrate-up, migrate-down, migrate-status, seed)
  - Docker commands (docker-up, docker-down, docker-logs, docker-restart)
  - Quick actions (start, stop)
  - Tool installation command (install)
- CLI tool using Cobra framework (used for code generation)
- Colored output in Makefile with comprehensive help system
- Self-documenting commands via `make help`
- Production middleware implementation
  - Sanitize middleware for XSS protection (applied globally)
  - Timeout middleware with 30-second timeout (applied globally)
  - Rate limiter middleware for brute force protection (10 req/min on auth endpoints)
 - Code generator templates moved to cmd/cli/commands/generator/templates/
- Documentation for generator package and templates
- Swagger API documentation in internal/delivery/http/docs/
- Swagger schemes annotation (http/https) for proper baseUrl configuration
- Swagger UI auto token management with Bearer Token authorization
  - Auto-save tokens after login/refresh
  - Auto-authorize Swagger's "Authorize" button with Bearer token
  - Auto-inject access token in Authorization header (Bearer {access_token})
  - Auto-fill refresh_token in request body
  - Auto-clear tokens on logout from localStorage and Swagger authorization
  - Visual token status indicator
  - Helper buttons (Show Tokens, Clear Tokens)
  - Info banner with usage instructions
  - Integration with Swagger's built-in security scheme (BearerAuth)

### Changed
- JWT configuration refactored from single JWT_SECRET to separate access/refresh token secrets
- JWT expiration now configurable separately for access and refresh tokens
- Updated all documentation to reflect new JWT configuration
- Primary interface now uses Makefile task runner (no direct CLI binary usage)
- All commands simplified: `make dev`, `make test`, `make gen-module product`
- Updated GitHub Actions CI/CD workflow to use Makefile
- Simplified Quick Start guide with Makefile commands
- Consolidated documentation structure (single source of truth for each topic)
- Documentation cleanup: removed redundant files, consolidated WebSocket docs
- All markdown files now consistent without emojis
- Response package now uses pagination.Meta to avoid redundancy
- Templates folder moved from root to cmd/cli/commands/generator/templates/
- Swagger docs moved from root docs/ to internal/delivery/http/docs/
- All CLI output messages use `make` commands instead of `app` binary
- Import path for templates updated to cmd/cli/commands/generator/templates
- WebSocket documentation consolidated into internal/websocket/README.md

### Removed
- Single JWT_SECRET and JWT_EXPIRATION environment variables (replaced with separate access/refresh configs)
- Hardcoded token expiration constants (now configurable via environment)
- SetJWTSecret() function (replaced with SetJWTConfig())
- Redundant documentation files (cmd/cli/README.md, redundant WebSocket docs)
- Duplicate workflow guides
- RBAC middleware (redundant - role checking done at handler level)
- Metrics handler and middleware (unused, not implemented)
- Redundant Pagination struct from response package
- App binary from root directory (now uses go run via Makefile)
- Templates folder from root level (moved to generator package)
- Docs folder from root level (moved to internal/delivery/http/docs/)
- All emojis from markdown documentation
- internal/websocket/examples.go (examples moved to README.md)
- Seven redundant WebSocket documentation files

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

### Changed
- **BREAKING**: Updated `health_handler.go` to require database and Redis client
  - Migration: Update handler initialization in `main.go` to pass `db` and `redisClient`
  
- **BREAKING**: Updated `postgres.go` to accept pool configuration
  - Migration: Update database initialization to pass `PoolConfig` struct

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

Run the updated install command to get Mockery:

```bash
make install
```

#### 5. Optional: Generate Mocks

If you want to use mocks for testing:

```bash
make mocks
```

#### 6. Optional: Use Templates

Check out the new templates in the `cmd/cli/commands/generator/templates/` directory for rapid feature development.

---

## Upcoming Features

### Planned for Next Release

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

- **Issues**: https://github.com/alpinnz/go-rest-api-boilerplate/issues
- **Discussions**: https://github.com/alpinnz/go-rest-api-boilerplate/discussions

---

[Unreleased]: https://github.com/alpinnz/go-rest-api-boilerplate/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/alpinnz/go-rest-api-boilerplate/releases/tag/v1.0.0

