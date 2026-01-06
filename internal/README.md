# Internal

Internal application code following Clean Architecture principles.

## Structure

```
internal/
├── domain/           # Domain layer - Core business entities and interfaces
├── usecase/          # Use case layer - Business logic orchestration
├── repository/       # Repository layer - Data access implementations
├── delivery/         # Delivery layer - HTTP handlers, DTOs, routes
├── middleware/       # HTTP middleware components
├── infrastructure/   # External dependencies (database, cache)
├── localization/     # Multi-language support
└── seeder/          # Database seeding utilities
```

## Layer Responsibilities

### Domain Layer (domain/)

Core business logic with no external dependencies.

- entity/ - Domain entities (User, Role) with business methods
- repository/ - Repository interfaces (contracts)
- errors.go - Domain-specific errors

Key principle: This layer should not depend on any other layer.

### Use Case Layer (usecase/)

Orchestrates business logic by coordinating repositories and entities.

Files:
- auth_usecase.go - Authentication (login, register, logout)
- user_usecase.go - User management
- role_usecase.go - Role management

Responsibilities:
- Business rule validation
- Transaction coordination
- Error handling with context

### Repository Layer (repository/)

Implements data access interfaces defined in domain layer.

Files:
- user_repository.go - User data access
- role_repository.go - Role data access
- session_repository.go - Redis session management

Responsibilities:
- Database queries
- Data mapping
- Connection management

### Delivery Layer (delivery/)

Handles HTTP request/response and API presentation.

Structure:
- http/dto/ - Data Transfer Objects for API
- http/handler/ - HTTP request handlers
- http/router/ - Route definitions and middleware setup

Responsibilities:
- Request validation
- Response formatting
- Route configuration

### Middleware Layer (middleware/)

Cross-cutting concerns for HTTP requests.

Available middleware:
- auth.go - JWT authentication
- cors.go - CORS configuration
- logger.go - Request logging
- recovery.go - Panic recovery
- rate_limiter.go - Rate limiting (applied to auth endpoints)
- timeout.go - Request timeout (30 seconds global)
- sanitize.go - Input sanitization (XSS protection)
- language.go - Language detection and i18n

All middleware actively used in router for production security and reliability.

### Infrastructure Layer (infrastructure/)

Manages external service connections.

Files:
- database/postgres.go - PostgreSQL connection with pooling
- database/redis.go - Redis connection for sessions

### Localization Layer (localization/)

Multi-language support system.

Files:
- loader.go - Language file loader
- types.go - Translation types
- lang/en.json - English translations
- lang/id.json - Indonesian translations

### Seeder (seeder/)

Database seeding utilities for development and testing.

## Design Principles

Dependency Rule (dependencies flow inward):
```
Delivery -> Use Case -> Repository -> Domain
```

Interface Segregation:
- Domain defines repository interfaces
- Use cases depend on repository interfaces
- Repositories implement domain interfaces

Dependency Injection through constructors:
```go
func NewUserUseCase(userRepo domain.UserRepository, sessionRepo domain.SessionRepository) *UserUseCase {
    return &UserUseCase{userRepo: userRepo, sessionRepo: sessionRepo}
}
```

## Common Patterns

Context Usage:
```go
func (uc *UserUseCase) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
    ctx, cancel := context.WithTimeout(ctx, uc.timeout)
    defer cancel()
    return uc.userRepo.FindByID(ctx, id)
}
```

Error Handling:
```go
if err != nil {
    return errors.Wrap(err, "USER_CREATE_FAILED", "Failed to create user").
        WithContext("email", user.Email)
}
```

Transaction Management:
```go
err := database.WithTransaction(ctx, db, func(ctx context.Context, tx *sql.Tx) error {
    // Multiple operations
    return nil
})
```

## Adding New Features

**Manual approach:**
1. Define entity in `domain/entity/`
2. Define repository interface in `domain/repository/`
3. Implement repository in `repository/`
4. Create use case in `usecase/`
5. Create DTOs in `delivery/http/dto/`
6. Create handler in `delivery/http/handler/`
7. Register routes in `delivery/http/router/`

**Automated approach:**
```bash
make gen-module name=product  # Generates all layers at once
```

This creates:
- Entity, repository interface, repository implementation
- Use case with business logic
- Handler with CRUD operations
- DTOs for requests and responses

## Testing

Each layer has corresponding tests following the project structure:

**Test locations:**
- `usecase/` - Use case tests with mocked repositories
- `repository/` - Repository integration tests with test database
- `delivery/http/handler/` - Handler tests with mocked use cases

**Generate mocks:**
```bash
make mocks  # Mocks available in domain/repository/mocks/
```

**Run tests:**
```bash
make test           # All tests
make test-coverage  # Coverage report
```
