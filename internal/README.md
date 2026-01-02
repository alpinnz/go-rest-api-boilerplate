# Internal Package

Application-specific code following Clean Architecture, SOLID principles, and Domain-Driven Design (DDD).

## Structure

```
internal/
├── domain/              Core business domain
│   ├── entity/         Domain entities (User, Role)
│   ├── repository/     Repository interfaces (contracts)
│   └── errors.go       Domain errors
├── repository/          Repository implementations
├── usecase/            Business logic orchestration
├── delivery/http/      HTTP layer (handlers, DTOs, router)
├── middleware/         HTTP middleware
├── localization/       Multi-language support
└── infrastructure/     External services (PostgreSQL, Redis)
```

## Design Principles

### Domain-Driven Design (DDD)

Entity: Objects with identity (ID), mutable, has lifecycle
- Examples: User, Role
- Location: `domain/entity/`

Repository Pattern:
- Interface in domain layer (contract): `domain/repository/`
- Implementation in infrastructure layer: `repository/`

### SOLID Principles

Single Responsibility:
- Entity: Domain data and behavior
- Repository Interface: Data access contract
- Repository Implementation: Data access logic
- UseCase: Business logic orchestration
- Handler: HTTP request/response handling

Dependency Inversion:
- UseCase depends on Repository interface (abstraction)
- Not on concrete implementation
- Easy to mock for testing

### Clean Architecture

Dependency Rule:
```
HTTP Handler → UseCase → Domain Entity
     ↓            ↓            ↑
    DTO      Repository   Repository
                Interface   Implementation
```

Dependencies flow inward. Domain layer has no dependencies.

## Layers

### Domain Layer (`domain/`)

Core business logic. No dependencies on external frameworks.

entity/:
```go
// entity/user.go - Domain entity with behavior
type User struct {
    ID        uuid.UUID  // UUIDv7 for time-ordered IDs
    FirstName string
    LastName  string
    Email     string
    Password  string
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt *time.Time
    Roles     []*Role
}

func (u *User) FullName() string {
    return u.FirstName + " " + u.LastName
}

func (u *User) IsDeleted() bool {
    return u.DeletedAt != nil
}
```

repository/:
```go
// repository/user_repository.go - Interface contract
type UserRepository interface {
    Create(ctx context.Context, user *entity.User) error
    FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
    FindByEmail(ctx context.Context, email string) (*entity.User, error)
    Update(ctx context.Context, user *entity.User) error
    Delete(ctx context.Context, id uuid.UUID) error
}
```

### Repository Layer (`repository/`)

Data access implementations.

```go
type userRepository struct {
    db *sql.DB
}

func (r *userRepository) FindByID(ctx context.Context, id int64) (*entity.User, error) {
    // SQL implementation
    // Auto exclude soft-deleted: WHERE deleted_at IS NULL
}
```

### UseCase Layer (`usecase/`)

Business logic orchestration. Coordinates between repositories and external services.

```go
type UserUseCase struct {
    userRepo domain.UserRepository
    roleRepo domain.RoleRepository
}

func (uc *UserUseCase) GetUser(ctx context.Context, id int64) (*entity.User, error) {
    user, err := uc.userRepo.FindByID(ctx, id)
    if err != nil {
        return nil, err
    }
    
    // Load user roles
    roles, _ := uc.roleRepo.FindByUserID(ctx, user.ID)
    user.Roles = roles
    
    return user, nil
}
```

### Delivery Layer (`delivery/http/`)

HTTP request/response handling.

Structure:
```
delivery/http/
├── dto/        Data Transfer Objects
├── handler/    HTTP handlers
└── router/     Route definitions
```

Flow:
1. Handler receives HTTP request
2. Parse and validate via DTO
3. Map DTO to UseCase input
4. Call UseCase
5. Map UseCase output to DTO
6. Return HTTP response

### Middleware Layer (`middleware/`)

HTTP middleware for cross-cutting concerns:
- auth.go - JWT authentication
- cors.go - CORS configuration
- language.go - Language detection
- logger.go - Request logging
- recovery.go - Panic recovery

### Infrastructure Layer (`infrastructure/`)

External service connections:
- database/postgres.go - PostgreSQL connection
- database/redis.go - Redis connection

### Localization Layer (`localization/`)

Multi-language support:
- Language detection from Accept-Language header
- Translation files in `lang/` directory
- Parameter interpolation

## Adding New Feature

1. Define entity in `domain/entity/`
2. Define repository interface in `domain/repository/`
3. Implement repository in `repository/`
4. Create use case in `usecase/`
5. Add DTO in `delivery/http/dto/`
6. Create handler in `delivery/http/handler/`
7. Register route in `delivery/http/router/`

## Example: Adding New Feature

Let's add a "Post" feature:

Step 1: Define entity (`domain/entity/post.go`):
```go
type Post struct {
    ID        int64
    UserID    int64
    Title     string
    Content   string
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt *time.Time
}

func (p *Post) IsDeleted() bool {
    return p.DeletedAt != nil
}
```

Step 2: Define repository interface (`domain/repository/post_repository.go`):
```go
type PostRepository interface {
    Create(ctx context.Context, post *entity.Post) error
    FindByID(ctx context.Context, id int64) (*entity.Post, error)
    FindByUserID(ctx context.Context, userID int64) ([]*entity.Post, error)
    Update(ctx context.Context, post *entity.Post) error
    Delete(ctx context.Context, id int64) error
}
```

Step 3: Implement repository (`repository/post_repository.go`):
```go
type postRepository struct {
    db *sql.DB
}

func NewPostRepository(db *sql.DB) domain.PostRepository {
    return &postRepository{db: db}
}

func (r *postRepository) FindByID(ctx context.Context, id int64) (*entity.Post, error) {
    // Implementation with WHERE deleted_at IS NULL
}
```

Step 4: Create use case (`usecase/post_usecase.go`):
```go
type PostUseCase struct {
    postRepo domain.PostRepository
    userRepo domain.UserRepository
}

func NewPostUseCase(postRepo domain.PostRepository, userRepo domain.UserRepository) *PostUseCase {
    return &PostUseCase{
        postRepo: postRepo,
        userRepo: userRepo,
    }
}

func (uc *PostUseCase) GetPost(ctx context.Context, id int64) (*entity.Post, error) {
    return uc.postRepo.FindByID(ctx, id)
}
```

Step 5: Add DTO (`delivery/http/dto/post_dto.go`):
```go
type CreatePostRequest struct {
    Title   string `json:"title" binding:"required,min=3,max=200"`
    Content string `json:"content" binding:"required,min=10"`
}

type PostResponse struct {
    ID        string `json:"id"`        // UUID as string
    Title     string `json:"title"`
    Content   string `json:"content"`
    CreatedAt string `json:"created_at"`
}
```

Step 6: Create handler (`delivery/http/handler/post_handler.go`):
```go
type PostHandler struct {
    postUseCase *usecase.PostUseCase
}

func NewPostHandler(postUseCase *usecase.PostUseCase) *PostHandler {
    return &PostHandler{postUseCase: postUseCase}
}

func (h *PostHandler) GetPost(c *gin.Context) {
    // Implementation
}
```

Step 7: Register route (`delivery/http/router/router.go`):
```go
posts := v1.Group("/posts")
posts.Use(middleware.Auth())
{
    posts.GET("/:id", postHandler.GetPost)
    posts.POST("", postHandler.CreatePost)
    posts.PUT("/:id", postHandler.UpdatePost)
    posts.DELETE("/:id", postHandler.DeletePost)
}
```

## Best Practices

1. Keep domain layer independent of frameworks
2. Use interfaces for repository contracts
3. Implement repository in separate layer
4. Business logic goes in use cases
5. Handlers only handle HTTP concerns
6. Use DTOs for API input/output
7. Always handle errors explicitly
8. Use context for cancellation and timeout
9. Implement soft delete for data preservation
10. Document all exported functions

## References

- Clean Architecture by Robert C. Martin
- Domain-Driven Design by Eric Evans
- SOLID Principles
- Go Project Layout: https://github.com/golang-standards/project-layout

