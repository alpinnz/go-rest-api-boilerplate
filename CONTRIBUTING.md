# Contributing Guidelines

Thank you for contributing to this project. Follow these guidelines to maintain code quality and consistency.

## Development Workflow

```bash
# 1. Fork repository
# 2. Clone fork
git clone https://github.com/YOUR_USERNAME/go-rest-api-boilerplate.git

# 3. Create feature branch
git checkout -b feat/new-feature

# 4. Make changes and test
make test

# 5. Format and vet code
make fmt && make vet

# 6. Commit with conventional commit format
git commit -m "feat(scope): add new feature"

# 7. Push and create pull request
git push origin feat/new-feature
```

## Branch Naming

Format: `<type>/<description>`

```
feat/add-user-profile       # New feature
fix/jwt-expiration-bug      # Bug fix
docs/update-readme          # Documentation
refactor/user-repository    # Code refactoring
test/auth-middleware        # Test additions
chore/update-dependencies   # Maintenance
```

## Commit Convention

Follow [Conventional Commits](https://www.conventionalcommits.org/) specification:

```
<type>(<scope>): <subject>

[optional body]

[optional footer]
```

### Types

| Type | Description | Example |
|------|-------------|---------|
| `feat` | New feature | `feat(auth): add JWT refresh token` |
| `fix` | Bug fix | `fix(user): handle null email` |
| `docs` | Documentation | `docs: update API endpoints` |
| `style` | Code style/format | `style: format with gofmt` |
| `refactor` | Code refactoring | `refactor(repo): simplify query` |
| `test` | Add/update tests | `test(auth): add login tests` |
| `chore` | Maintenance | `chore: update dependencies` |
| `perf` | Performance | `perf(db): optimize query` |

### Scope

Module or layer being changed: `auth`, `user`, `middleware`, `config`, `db`

### Subject

- Use imperative mood: "add" not "added" or "adds"
- Lowercase, no period at the end
- Maximum 50 characters

### Examples

```bash
# Good
feat(auth): add JWT refresh token
fix(user): prevent duplicate email registration
docs: update API documentation
refactor(repo): use domain-driven structure
test(middleware): add auth middleware tests

# Bad
feat: Added new feature for authentication system
fix: fixed bug
Update README.md
refactored code
```

## Code Documentation

### Go Doc Comments

Follow these principles for all exported code:

```go
// UserRepository provides data access methods for user entities.
// All methods accept context for cancellation and timeout control.
type UserRepository interface {
    // Create inserts a new user into the database.
    // Returns error if email already exists or database operation fails.
    Create(ctx context.Context, user *User) error
    
    // FindByID retrieves a user by their unique identifier.
    // Returns ErrNotFound if user does not exist.
    FindByID(ctx context.Context, id int64) (*User, error)
    
    // FindByEmail retrieves a user by email address.
    // Email comparison is case-insensitive.
    // Returns ErrNotFound if user does not exist.
    FindByEmail(ctx context.Context, email string) (*User, error)
}
```

### Documentation Rules

1. First sentence is a brief summary (one line)
2. Use imperative tone (not "This function will...")
3. Explain:
   - Purpose
   - Important parameters
   - Return values
   - Errors or panics
4. Don't repeat what's obvious from the function name
5. Add examples for non-trivial logic

### Example: Good Documentation

```go
// HashPassword generates a bcrypt hash from plaintext password.
// Uses bcrypt cost factor of 12 for security vs performance balance.
// Returns error if password is empty or hashing fails.
func HashPassword(password string) (string, error) {
    if password == "" {
        return "", errors.New("password cannot be empty")
    }
    
    hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return "", fmt.Errorf("failed to hash password: %w", err)
    }
    
    return string(hash), nil
}

// GenerateToken creates a JWT token for authenticated user.
// Token includes user ID and expires based on configured expiration time.
// Returns token string and error if generation fails.
//
// Example:
//   token, err := GenerateToken(12345)
//   if err != nil {
//       return err
//   }
//   // Use token for Authorization header
func GenerateToken(userID int64) (string, error) {
    // Implementation...
}
```

### Example: Bad Documentation

```go
// This function will hash the password using bcrypt
func HashPassword(password string) (string, error) {
    // Bad: Uses "will", not imperative
    // Bad: Missing error conditions
    // Bad: Missing return value explanation
}

// Hash password
func HashPassword(password string) (string, error) {
    // Bad: Too brief, no useful information
}
```

## Code Style

### Go Standards

```go
// Good: Clear naming, error handling, context
func (r *userRepository) FindByID(ctx context.Context, id int64) (*domain.User, error) {
    user := &domain.User{}
    err := r.db.QueryRowContext(ctx, query, id).Scan(
        &user.ID,
        &user.Email,
        &user.Password,
        &user.Name,
        &user.CreatedAt,
        &user.UpdatedAt,
    )
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, domain.ErrNotFound
        }
        return nil, fmt.Errorf("failed to find user: %w", err)
    }
    return user, nil
}

// Bad: Poor naming, no error handling, no context
func (r *userRepository) Get(id int64) *domain.User {
    u := &domain.User{}
    r.db.QueryRow(query, id).Scan(&u.ID, &u.Email)
    return u
}
```

### Rules

1. **Format**: Run `make fmt` before commit
2. **Vet**: Run `make vet` and fix all issues
3. **Lint**: Run `make lint` (requires golangci-lint)
4. **Naming**: Use clear, descriptive names
5. **Errors**: Always handle errors explicitly with error wrapping
6. **Context**: Pass context.Context for cancellation and timeout
7. **Comments**: Document all exported functions and types
8. **Tests**: Add tests for new features

### REST API Conventions

1. **HTTP Status Codes**: Use appropriate codes
   - `200 OK` - Successful read/update
   - `201 Created` - Successful resource creation
   - `204 No Content` - Successful deletion
   - `400 Bad Request` - Invalid input
   - `401 Unauthorized` - Missing/invalid authentication
   - `404 Not Found` - Resource not found
   - `409 Conflict` - Duplicate resource
   - `500 Internal Server Error` - Unexpected error

2. **Request Validation**: Use struct tags
```go
type CreateUserRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=8"`
    Name     string `json:"name" binding:"required,min=2,max=100"`
}
```

3. **Response Structure**: Keep consistent
```go
// Success
response.Success(c, http.StatusOK, data)

// Error
response.Error(c, http.StatusBadRequest, err)
```

4. **URL Naming**:
   - Lowercase: `/api/v1/users` not `/api/v1/Users`
   - Plural for collections: `/users` not `/user`
   - Kebab-case: `/user-profiles` not `/userProfiles`

## Testing Requirements

### Guidelines

- Write unit tests for all use cases
- Write integration tests for repositories
- Write HTTP tests for handlers
- Use table-driven tests for multiple cases
- Mock external dependencies
- Target 80%+ test coverage

### Test Structure

```go
func TestUserUseCase_Register(t *testing.T) {
    // Arrange
    mockUserRepo := &mockUserRepository{
        users: make(map[int64]*domain.User),
    }
    mockSessionRepo := &mockSessionRepository{
        sessions: make(map[string]int64),
    }
    
    uc := usecase.NewUserUseCase(mockUserRepo, mockSessionRepo, time.Hour)
    
    input := usecase.RegisterInput{
        Email:    "test@example.com",
        Password: "password123",
        Name:     "Test User",
    }
    
    // Act
    user, err := uc.Register(context.Background(), input)
    
    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, user)
    assert.Equal(t, input.Email, user.Email)
    assert.Equal(t, input.Name, user.Name)
    assert.NotEmpty(t, user.Password) // hashed
}
```

### Table-Driven Tests

```go
func TestUserRepository_FindByEmail(t *testing.T) {
    tests := []struct {
        name    string
        email   string
        wantErr bool
    }{
        {
            name:    "existing user",
            email:   "existing@example.com",
            wantErr: false,
        },
        {
            name:    "non-existing user",
            email:   "notfound@example.com",
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            user, err := repo.FindByEmail(context.Background(), tt.email)
            if tt.wantErr {
                assert.Error(t, err)
                assert.Nil(t, user)
            } else {
                assert.NoError(t, err)
                assert.NotNil(t, user)
            }
        })
    }
}
```

## Pull Request Process

1. Update README.md if needed
2. Update documentation if API changed
3. Ensure all tests pass (`make test`)
4. Ensure code is formatted (`make fmt && make vet`)
5. Request review from maintainers
6. Address review comments
7. Squash commits if requested

## Code Review Guidelines

Reviewers should verify:

- Code follows project architecture
- Tests are adequate and passing
- Documentation is updated
- No breaking changes without discussion
- Performance implications considered
- Security best practices followed
- Error handling is comprehensive

## Questions

For questions or discussions, open an issue with the `question` label.

