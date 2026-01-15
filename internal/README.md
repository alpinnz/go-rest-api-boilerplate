# Internal

Internal application code following Clean Architecture principles.

This folder contains the application-specific implementation. For setup, environment variables, Make commands, and API endpoints, keep the single source of truth in the root `README.md`.

## Structure

```
internal/
├── container/        # Dependency injection container (wiring + lifecycle)
├── domain/           # Domain layer: entities, errors, and repository interfaces (contracts)
├── usecase/          # Use case layer: business logic orchestration
├── repository/       # Infrastructure-facing repository implementations (DB/Redis)
├── delivery/         # Delivery layer: HTTP handlers, DTOs, routes
├── middleware/       # HTTP middleware components
├── infrastructure/   # External connections (PostgreSQL, Redis)
├── localization/     # Multi-language support
├── websocket/        # WebSocket hub and client management
└── seeder/           # Database seeding utilities
```

## Layer Responsibilities

### Container Layer (`container/`)

Dependency injection container that wires dependencies and manages resource lifecycle.

Responsibilities:
- Initialize resources (DB, Redis, localization, etc.)
- Construct repositories, use cases, handlers, and middleware
- Own cleanup (close DB/Redis) for graceful shutdown

Why it exists:
- Keeps `cmd/api/main.go` small
- Avoids scattered/manual wiring
- Makes testing easier by injecting dependencies

### Domain Layer (`domain/`)

Core business rules with no dependencies on other internal layers.

- `entity/` - Domain entities (e.g. User, Role)
- `repository/` - Repository interfaces (contracts)
- `errors.go` - Domain-specific errors

Rule of thumb: The domain layer should not import `internal/...` implementations.

### Use Case Layer (`usecase/`)

Orchestrates business flows by coordinating domain entities and repository interfaces.

Responsibilities:
- Business rule validation
- Transaction coordination (when needed)
- Error mapping/wrapping (with context)

### Repository Layer (`repository/`)

Implements the interfaces defined in `internal/domain/repository`.

Responsibilities:
- Data access and mapping
- DB/Redis interactions

Naming note:
- `internal/domain/repository/*` = interfaces (contracts)
- `internal/repository/*` = implementations (Postgres/Redis)

### Delivery Layer (`delivery/`)

HTTP request/response concerns.

Structure:
- `http/dto/` - request/response DTOs
- `http/handler/` - HTTP handlers
- `http/router/` - routes + middleware setup

### Middleware Layer (`middleware/`)

Cross-cutting concerns for HTTP requests.

Available middleware:
- `auth.go` - JWT authentication + session validation
- `cors.go` - CORS (configurable via env)
- `logger.go` - structured request logging
- `recovery.go` - panic recovery + logging
- `rate_limiter.go` - configurable per-IP rate limiting
- `timeout.go` - request timeout
- `sanitize.go` - input sanitization
- `language.go` - language detection and i18n

### Infrastructure Layer (`infrastructure/`)

External service connections.

- `database/postgres.go` - PostgreSQL connection + pooling
- `database/redis.go` - Redis client for sessions/caching

### Localization Layer (`localization/`)

Multi-language support system.

### WebSocket (`websocket/`)

Real-time communication using Hub pattern.

### Seeder (`seeder/`)

Database seeding utilities.

## Design Principles

Dependency rule (source code dependencies point inward):

```
Delivery  -> Usecase -> Domain
Repository->/
Infrastructure ->/
Middleware -> Delivery (HTTP pipeline)
```

Interface segregation:
- Domain defines repository interfaces
- Use cases depend on those interfaces
- Repositories implement those interfaces

## Testing Notes

General guidance:
- Use case tests: mock repository interfaces
- Repository tests: integration tests (DB/Redis)
- Handler tests: mock use cases or use a test router

Mocks (if generated) should implement interfaces from `internal/domain/repository`.
