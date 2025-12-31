# Domain Package

Core business entities and domain errors.

## Overview

Domain package contains business entities, interfaces, and domain-specific errors. This is the core of the application that defines business rules and models.

## Structure

```
domain/
├── errors.go       Domain-specific errors
├── user.go         User entity and repository interface
└── session.go      Session repository interface
```

## Domain Errors

Predefined errors for consistent error handling across layers.

### Available Errors

| Error | HTTP Status | Description |
|-------|-------------|-------------|
| `ErrNotFound` | 404 | Resource does not exist |
| `ErrConflict` | 409 | Resource already exists (duplicate) |
| `ErrInvalidInput` | 400 | Input validation failed |
| `ErrUnauthorized` | 401 | Missing or invalid authentication |
| `ErrForbidden` | 403 | User lacks permission |
| `ErrInternalServer` | 500 | Unexpected server error |
| `ErrBadRequest` | 400 | Malformed request |
| `ErrInvalidCredentials` | 401 | Authentication failed |
| `ErrSessionExpired` | 401 | Session token expired |

### Usage

Return domain errors from repository or usecase layer:

```go
func (r *UserRepository) FindByID(ctx context.Context, id int64) (*domain.User, error) {
    user, err := r.db.FindUser(id)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, domain.ErrNotFound
        }
        return nil, err
    }
    return user, nil
}
```

Map domain errors to HTTP responses in handler:

```go
func (h *UserHandler) GetByID(c *gin.Context) {
    user, err := h.userUseCase.GetByID(c.Request.Context(), id)
    if err != nil {
        if err == domain.ErrNotFound {
            response.NotFound(c, "user")
            return
        }
        response.InternalServerError(c, "Failed to fetch user")
        return
    }
    
    response.Success(c, user)
}
```

## Domain Entities

### User

Pure domain entity representing a user in the system.

**Fields:**
- `ID` - Unique user identifier
- `Email` - User email address (unique)
- `Password` - Hashed password (bcrypt)
- `Name` - User full name
- `CreatedAt` - Account creation timestamp
- `UpdatedAt` - Last update timestamp
- `DeletedAt` - Soft delete timestamp (nullable)

**Note**: User entity has no JSON tags. Serialization is handled by DTOs in the delivery layer.

**Repository Interface:**

```go
type UserRepository interface {
    Create(ctx context.Context, user *User) error
    FindByID(ctx context.Context, id int64) (*User, error)
    FindByEmail(ctx context.Context, email string) (*User, error)
    Update(ctx context.Context, user *User) error
    Delete(ctx context.Context, id int64) error
}
```

## Design Principles

### 1. Domain Independence

Domain layer is independent from infrastructure:

```go
// Good: Domain interface
type UserRepository interface {
    FindByID(ctx context.Context, id int64) (*User, error)
}

// Bad: Infrastructure coupling
type UserRepository interface {
    FindByID(ctx context.Context, id int64) (*gorm.Model, error)
}
```

### 2. Consistent Error Handling

Use predefined domain errors:

```go
// Good: Domain error
if user == nil {
    return domain.ErrNotFound
}

// Bad: Generic error
if user == nil {
    return errors.New("user not found")
}
```

### 3. Business Logic in Domain

Keep business rules in domain layer:

```go
// Good: Business logic in domain
func (u *User) CanDelete() error {
    if u.IsAdmin() {
        return domain.ErrForbidden
    }
    return nil
}

// Bad: Business logic in handler
func (h *Handler) DeleteUser(c *gin.Context) {
    if user.IsAdmin() {
        response.Forbidden(c, "user.delete")
        return
    }
}
```

### 4. Interface Segregation

Define interfaces in domain, implement in infrastructure:

```go
// domain/user.go
type UserRepository interface {
    FindByID(ctx context.Context, id int64) (*User, error)
}

// repository/user_repository.go
type userRepository struct {
    db *sql.DB
}

func (r *userRepository) FindByID(ctx context.Context, id int64) (*domain.User, error) {
    // Implementation
}
```

## Error Mapping Examples

### Repository Layer

```go
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
    var user domain.User
    query := "SELECT id, email, password, name, created_at, updated_at FROM users WHERE email = $1"
    err := r.db.QueryRowContext(ctx, query, email).
        Scan(&user.ID, &user.Email, &user.Password, &user.Name, &user.CreatedAt, &user.UpdatedAt)
    
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, domain.ErrNotFound
        }
        return nil, err
    }
    
    return &user, nil
}
```

### Usecase Layer

```go
func (uc *UserUseCase) Login(ctx context.Context, input LoginInput) (*LoginResult, error) {
    user, err := uc.userRepo.FindByEmail(ctx, input.Email)
    if err != nil {
        if err == domain.ErrNotFound {
            return nil, domain.ErrInvalidCredentials
        }
        return nil, err
    }
    
    if !auth.CheckPasswordHash(input.Password, user.Password) {
        return nil, domain.ErrInvalidCredentials
    }
    
    return &LoginResult{User: user, Token: token}, nil
}
```

### Handler Layer

```go
func (h *UserHandler) Login(c *gin.Context) {
    result, err := h.userUseCase.Login(c.Request.Context(), input)
    if err != nil {
        switch err {
        case domain.ErrInvalidCredentials:
            response.Unauthorized(c, "Invalid email or password")
        case domain.ErrNotFound:
            response.NotFound(c, "user")
        default:
            response.InternalServerError(c, "Failed to process login")
        }
        return
    }
    
    response.Success(c, result)
}
```

