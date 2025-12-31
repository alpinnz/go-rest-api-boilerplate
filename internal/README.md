# internal/

Application-specific code following Clean Architecture principles.

## Directory Structure

```
internal/
├── delivery/http/       HTTP handlers and routing
│   ├── dto/            Data Transfer Objects
│   ├── handler/        HTTP handlers
│   └── router/         Route definitions
├── domain/              Business entities and interfaces
├── infrastructure/      External service integrations
├── localization/        Multi-language support
├── middleware/          HTTP middleware
├── repository/          Data access implementations
└── usecase/             Business logic orchestration
```

## Architecture Overview

This boilerplate follows Clean Architecture with clear separation of concerns:

```
┌─────────────────────────────────────────────┐
│           Delivery Layer (HTTP)             │
│  - Handlers                                 │
│  - Routing                                  │
│  - Request/Response mapping                 │
└──────────────────┬──────────────────────────┘
                   │
┌──────────────────▼──────────────────────────┐
│          Middleware Layer                   │
│  - Authentication                           │
│  - CORS, Logging, Recovery                  │
└──────────────────┬──────────────────────────┘
                   │
┌──────────────────▼──────────────────────────┐
│           Usecase Layer                     │
│  - Business logic orchestration             │
│  - Application services                     │
└──────────────────┬──────────────────────────┘
                   │
┌──────────────────▼──────────────────────────┐
│            Domain Layer                     │
│  - Entities                                 │
│  - Repository interfaces                    │
│  - Domain errors                            │
└──────────────────┬──────────────────────────┘
                   │
┌──────────────────▼──────────────────────────┐
│         Repository Layer                    │
│  - Database implementations                 │
│  - Cache implementations                    │
└──────────────────┬──────────────────────────┘
                   │
┌──────────────────▼──────────────────────────┐
│       Infrastructure Layer                  │
│  - Database connections                     │
│  - External services                        │
└─────────────────────────────────────────────┘
```

## Packages

### delivery/http/

HTTP layer for handling requests and responses.

**Structure:**
```
delivery/http/
├── dto/        Data Transfer Objects (JSON serialization)
├── handler/    Request handlers
└── router/     Route definitions
```

**Responsibilities:**
- Parse HTTP requests
- Validate request format (via DTOs)
- Map DTOs to usecase inputs
- Call usecase layer
- Map usecase outputs to DTOs
- Format HTTP responses

**DTO Pattern:**
- DTOs contain JSON and validation tags
- Handlers map between DTOs and domain entities
- Separation between HTTP format and business logic

### domain/

Core business entities and domain logic.

**Files:**
- `user.go` - User entity and repository interface
- `errors.go` - Domain-specific errors

**Characteristics:**
- No external dependencies
- Pure business logic
- Technology-agnostic

[View detailed documentation](./domain/README.md)

### infrastructure/

External service integrations and connections.

**Structure:**
```
infrastructure/
└── database/
    ├── postgres.go    PostgreSQL connection
    └── redis.go       Redis connection
```

**Responsibilities:**
- Database connections
- External service clients
- Infrastructure configuration

### middleware/

HTTP middleware for cross-cutting concerns.

**Files:**
- `auth.go` - JWT authentication
- `cors.go` - CORS headers
- `logger.go` - Request logging
- `recovery.go` - Panic recovery

[View detailed documentation](./middleware/README.md)

### repository/

Data access layer implementations.

**Files:**
- `user_repository.go` - User database operations
- `session_repository.go` - Session cache operations

**Responsibilities:**
- Implement domain repository interfaces
- Execute database queries
- Handle data persistence

### usecase/

Business logic orchestration.

**Files:**
- `user_usecase.go` - User business operations

**Responsibilities:**
- Orchestrate business workflows
- Coordinate repositories
- Apply business rules
- Transaction management

## Dependency Flow

Dependencies flow inward following Clean Architecture:

```
delivery → middleware → usecase → domain ← repository ← infrastructure
```

Rules:
1. Inner layers cannot depend on outer layers
2. Domain is the innermost layer (no dependencies)
3. Infrastructure is the outermost layer
4. Dependencies point inward through interfaces

## Example: User Registration Flow

```
1. HTTP Request
   ↓
2. Handler (delivery/http/handler/user_handler.go)
   - Parse JSON
   - Validate format
   ↓
3. Middleware (middleware/auth.go)
   - Optional authentication check
   ↓
4. Usecase (usecase/user_usecase.go)
   - Hash password
   - Check duplicate email
   - Create user
   - Generate token
   ↓
5. Repository (repository/user_repository.go)
   - Insert into database
   ↓
6. Infrastructure (infrastructure/database/postgres.go)
   - Execute SQL query
   ↓
7. HTTP Response
   - Return user and token
```

## Package Guidelines

### domain/

```go
// Define entities
type User struct {
    ID    int64
    Email string
}

// Define interfaces
type UserRepository interface {
    Create(ctx context.Context, user *User) error
}

// Define errors
var ErrNotFound = errors.New("resource not found")
```

### usecase/

```go
// Depend on domain interfaces
type UserUseCase struct {
    userRepo domain.UserRepository
}

// Implement business logic
func (uc *UserUseCase) Register(ctx context.Context, input RegisterInput) error {
    // Business logic here
}
```

### repository/

```go
// Implement domain interface
type userRepository struct {
    db *sql.DB
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
    // Database operation
}
```

### delivery/

```go
// Handle HTTP requests
func (h *UserHandler) Register(c *gin.Context) {
    var input RegisterInput
    if err := c.ShouldBindJSON(&input); err != nil {
        response.BadRequest(c, "INVALID_INPUT", "Invalid request", nil)
        return
    }
    
    result, err := h.userUseCase.Register(c.Request.Context(), input)
    if err != nil {
        // Handle error
        return
    }
    
    response.Created(c, result)
}
```

## Testing Strategy

### Unit Tests

Test each layer independently using mocks:

```go
// Test usecase with mock repositories
func TestUserUseCase_Register(t *testing.T) {
    mockRepo := new(MockUserRepository)
    usecase := NewUserUseCase(mockRepo)
    
    mockRepo.On("Create", mock.Anything).Return(nil)
    
    err := usecase.Register(ctx, input)
    assert.NoError(t, err)
}
```

### Integration Tests

Test repository implementations with real databases:

```go
// Test repository with test database
func TestUserRepository_Create(t *testing.T) {
    db := setupTestDB(t)
    defer db.Close()
    
    repo := NewUserRepository(db)
    user := &domain.User{Email: "test@example.com"}
    
    err := repo.Create(context.Background(), user)
    assert.NoError(t, err)
    assert.NotZero(t, user.ID)
}
```

### Handler Tests

Test HTTP endpoints with mock usecases:

```go
// Test handler with mock usecase
func TestUserHandler_Register(t *testing.T) {
    gin.SetMode(gin.TestMode)
    
    mockUseCase := new(MockUserUseCase)
    handler := NewUserHandler(mockUseCase)
    
    router := gin.New()
    router.POST("/register", handler.Register)
    
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("POST", "/register", body)
    router.ServeHTTP(w, req)
    
    assert.Equal(t, 201, w.Code)
}
```

## Data Flow Examples

### User Registration Flow

```
1. HTTP POST /api/v1/register
   {
     "email": "user@example.com",
     "password": "password123",
     "name": "John Doe"
   }
   ↓
2. Handler (delivery/http/handler/user_handler.go)
   - Parse JSON body
   - Validate request format
   ↓
3. Usecase (usecase/user_usecase.go)
   - Check email uniqueness
   - Hash password with bcrypt
   - Create user entity
   ↓
4. Repository (repository/user_repository.go)
   - Insert user into PostgreSQL
   - Return created user with ID
   ↓
5. Response
   201 Created
   {
     "code": null,
     "message": "Resource created successfully",
     "data": {
       "id": 1,
       "email": "user@example.com",
       "name": "John Doe",
       "created_at": "2025-12-31T10:00:00Z"
     }
   }
```

### User Login Flow

```
1. HTTP POST /api/v1/login
   {
     "email": "user@example.com",
     "password": "password123"
   }
   ↓
2. Handler (delivery/http/handler/user_handler.go)
   - Parse JSON body
   - Validate request format
   ↓
3. Usecase (usecase/user_usecase.go)
   - Get user by email from repository
   - Verify password with bcrypt
   - Generate JWT token
   - Store session in Redis
   ↓
4. Repository
   - User Repository: Get user from PostgreSQL
   - Session Repository: Store token in Redis
   ↓
5. Response
   200 OK
   {
     "code": null,
     "message": "Success",
     "data": {
       "user": {...},
       "token": "eyJhbGciOiJIUzI1NiIs..."
     }
   }
```

### Protected Request Flow

```
1. HTTP GET /api/v1/profile
   Headers:
     Authorization: Bearer <token>
   ↓
2. Auth Middleware (middleware/auth.go)
   - Extract token from header
   - Validate Bearer format
   ↓
3. Usecase (usecase/user_usecase.go)
   - Validate JWT signature
   - Check session in Redis
   - Extract user ID from token
   ↓
4. Repository (repository/session_repository.go)
   - Get session from Redis
   - Verify token exists and not expired
   ↓
5. Handler (delivery/http/handler/user_handler.go)
   - Get userID from context
   - Fetch user profile
   ↓
6. Response
   200 OK
   {
     "code": null,
     "message": "Success",
     "data": {
       "id": 1,
       "email": "user@example.com",
       "name": "John Doe"
     }
   }
```

## Design Decisions

### Repository Organization by Domain

Repositories are organized by business entity, not by technology:

```
repository/
├── user_repository.go      # User domain (PostgreSQL)
└── session_repository.go   # Session domain (Redis)

NOT:
repository/
├── postgres/
│   └── user_repository.go
└── redis/
    └── session_repository.go
```

**Benefits:**
- Single repository can use multiple data sources
- Clean Architecture compliance
- Domain-driven design alignment
- Easier to test with mocks
- Infrastructure layer is reusable

### Infrastructure Layer Separation

Database connections are separated into infrastructure layer:

```go
// Infrastructure provides connections
db := infrastructure.NewPostgresDB(config)
redis := infrastructure.NewRedisClient(config)

// Repositories use those connections
userRepo := repository.NewUserRepository(db)
sessionRepo := repository.NewSessionRepository(redis)
```

**Benefits:**
- Single source of truth for connections
- Reusable across multiple repositories
- Easy to switch database technology
- Clear separation of concerns
- Simplified testing

### Redis-Only Sessions

Sessions stored only in Redis, not in PostgreSQL:

**Rationale:**
- Fast in-memory access
- Built-in TTL (time-to-live)
- Automatic expiration
- Simpler architecture
- Horizontal scaling
- No cleanup jobs needed

**When to add SQL sessions:**
- Compliance requires audit trail
- Need session history analytics
- Multi-device session management
- Admin session revocation

### Error Handling

Domain errors for consistent error handling:

```go
// Domain layer defines errors
var ErrNotFound = errors.New("resource not found")
var ErrConflict = errors.New("resource already exists")

// Repository returns domain errors
if err == sql.ErrNoRows {
    return nil, domain.ErrNotFound
}

// Handler maps to HTTP responses
if err == domain.ErrNotFound {
    response.NotFound(c, "user")
}
```

**Benefits:**
- Consistent error handling
- Technology-agnostic
- Easy to test
- Clear error mapping

## References

- [Clean Architecture by Robert C. Martin](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Domain-Driven Design](https://martinfowler.com/bliki/DomainDrivenDesign.html)
- [Repository Pattern](https://martinfowler.com/eaaCatalog/repository.html)

### Example

```go
type userRepository struct {
    db *sql.DB
}

func NewUserRepository(db *sql.DB) domain.UserRepository {
    return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
    // Execute SQL INSERT
}
```

## delivery/http/

HTTP layer (handlers, routes, middleware).

### Structure

```
delivery/http/
├── handler/          Request handlers
├── router/           Route definitions
└── middleware/       HTTP middlewares
```

### handler/

HTTP request handlers.

**Files:**
- `user_handler.go` - User endpoints
- `health_handler.go` - Health check

**Responsibilities:**
- Parse HTTP request
- Validate input
- Call use case
- Format response

### router/

Route configuration.

**Files:**
- `router.go` - Route setup

**Responsibilities:**
- Register routes
- Apply middlewares
- Group endpoints

### middleware/

HTTP middlewares.

**Files:**
- `auth.go` - JWT authentication
- `cors.go` - CORS configuration
- `logger.go` - Request logging
- `recovery.go` - Panic recovery

## infrastructure/

External service connections.

### Structure

```
infrastructure/
└── database/         Database connections
```

### Characteristics

- Manages connections
- Reusable across layers
- Configuration-driven
- Error handling

### Files

- `postgres.go` - PostgreSQL connection
- `redis.go` - Redis client

### Example

```go
// NewPostgresDB creates PostgreSQL connection
func NewPostgresDB(dsn string) (*sql.DB, error) {
    db, err := sql.Open("postgres", dsn)
    if err != nil {
        return nil, err
    }
    
    // Configure connection pool
    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(25)
    db.SetConnMaxLifetime(5 * time.Minute)
    
    return db, nil
}
```

## References

- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Go Project Layout](https://github.com/golang-standards/project-layout)

