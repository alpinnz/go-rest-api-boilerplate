# Seeder Package

Database seeding package following SOLID principles.

## Structure

```
internal/seeder/
├── types.go         - Data types and interfaces
├── runner.go        - Seeder orchestration
├── user_seeder.go   - User data seeding
└── role_seeder.go   - Role data seeding
```

## Design Principles

### Single Responsibility Principle (SRP)
- Each seeder handles only one entity type
- Runner only orchestrates the seeding process
- Main.go only handles initialization

### Open/Closed Principle (OCP)
- Easy to add new seeders without modifying existing code
- Implement Seeder interface for new entities

### Dependency Inversion Principle (DIP)
- Seeders depend on repository interfaces
- Runner depends on Seeder interface

## Usage

### Running Seeder

```bash
# Run seeder with fresh mode (truncates tables first)
make seed

# Or directly
go run cmd/seeder/main.go
```

### Default Seeded Data

**Roles:**
- admin - Administrator with full access
- user - Regular user with standard access
- moderator - Moderator with elevated access

**Users:**
- admin@example.com / !Password123 (Admin User)
  - Roles: admin, user
- user@example.com / !Password123 (Regular User)
  - Roles: user
- test@example.com / !Password123 (Test User)
  - Roles: user, moderator

**User-Role Assignments:**
All users are automatically assigned their respective roles during seeding.

## Adding New Seeder

### Step 1: Create Seeder File

Create `internal/seeder/post_seeder.go`:

```go
package seeder

import (
    "context"
    "fmt"
    "log"
    
    "github.com/alpinnz/go-rest-api-boilerplate/internal/domain/repository"
)

type PostSeeder struct {
    postRepo repository.PostRepository
}

func NewPostSeeder(postRepo repository.PostRepository) *PostSeeder {
    return &PostSeeder{
        postRepo: postRepo,
    }
}

func (s *PostSeeder) Seed(ctx context.Context) error {
    log.Println("Seeding posts...")
    
    // Add seeding logic here
    
    log.Println("Posts seeding completed!")
    return nil
}
```

### Step 2: Register in main.go

```go
// Initialize repository
postRepo := repository.NewPostRepository(db)

// Register seeder
runner.Register(seeder.NewPostSeeder(postRepo))
```

### Step 3: Add to Truncate (if needed)

Update `internal/seeder/runner.go`:

```go
tables := []string{
    "user_roles",
    "posts",     // Add new table
    "roles",
    "users",
}
```

## Configuration

### Fresh Mode

```go
options := seeder.RunOptions{
    Fresh: true, // Truncate tables before seeding
}
```

**Fresh: true**
- Truncates all tables
- Resets auto-increment IDs
- Starts with clean slate

**Fresh: false**
- Keeps existing data
- Skips duplicates
- Idempotent seeding

## Best Practices

1. **Check for Existence**: Always check if data exists before inserting
2. **Idempotent**: Seeder should be safe to run multiple times
3. **Order Matters**: Seed in correct order (roles before users)
4. **Error Handling**: Return descriptive errors
5. **Logging**: Log progress for debugging
6. **Transactions**: Consider using transactions for data consistency

## Architecture

```
cmd/seeder/main.go (Initialization)
        ↓
internal/seeder/runner.go (Orchestration)
        ↓
internal/seeder/*_seeder.go (Specific Seeders)
        ↓
internal/repository/ (Data Access)
        ↓
Database
```

## Examples

### Custom Seeding Order

```go
// Register seeders in specific order
runner.Register(seeder.NewRoleSeeder(roleRepo))
runner.Register(seeder.NewUserSeeder(userRepo))
runner.Register(seeder.NewPostSeeder(postRepo))
```

### Conditional Seeding

```go
if cfg.App.Env == "development" {
    runner.Register(seeder.NewTestDataSeeder(db))
}
```

### Custom Run Options

```go
options := seeder.RunOptions{
    Fresh: cfg.App.Env == "development",
}
```

## Testing

Test seeders individually:

```go
func TestUserSeeder_Seed(t *testing.T) {
    // Setup mock repository
    mockRepo := &mockUserRepository{}
    seeder := NewUserSeeder(mockRepo)
    
    // Run seeder
    err := seeder.Seed(context.Background())
    
    // Assert
    assert.NoError(t, err)
    assert.Equal(t, 3, mockRepo.createCount)
}
```

## Troubleshooting

**"User already exists"**
- Expected behavior when Fresh: false
- Set Fresh: true to start clean

**"Failed to truncate tables"**
- Check foreign key constraints
- Ensure CASCADE is used

**"Failed to create user"**
- Check database connection
- Verify migration has run
- Check for unique constraint violations

