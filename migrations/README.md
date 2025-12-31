# Database Management

SQL migrations and database seeding for schema management.

## Overview

This directory contains SQL migration files for managing database schema changes. Migrations use timestamp-based versioning and support both up (apply) and down (rollback) operations.

## File Structure

```
migrations/
├── README.md
├── 20251231062946_create_table_users.up.sql
└── 20251231062946_create_table_users.down.sql
```

## Naming Convention

```
<timestamp>_<description>.<direction>.sql
```

**Format:**
- **Timestamp:** `YYYYMMDDHHMMSS` format
- **Description:** Snake case, descriptive name
- **Direction:** `up` (apply) or `down` (rollback)

**Examples:**
- `20251231062946_create_table_users.up.sql`
- `20251231062946_create_table_users.down.sql`
- `20260101120000_add_user_avatar.up.sql`
- `20260101120000_add_user_avatar.down.sql`

## Running Migrations

### Using Makefile

```bash
# Run all pending migrations
make migrate-up

# Rollback last migration
make migrate-down

# Create new migration
make migrate-create NAME=add_user_avatar
```

### Using golang-migrate CLI

```bash
# Install golang-migrate
brew install golang-migrate  # macOS
# or download from https://github.com/golang-migrate/migrate

# Run all pending migrations
migrate -path migrations \
  -database "postgresql://postgres:postgres@localhost:5432/go-rest-api-boilerplate?sslmode=disable" \
  up

# Rollback last migration
migrate -path migrations \
  -database "postgresql://postgres:postgres@localhost:5432/go-rest-api-boilerplate?sslmode=disable" \
  down 1

# Check version
migrate -path migrations \
  -database "postgresql://postgres:postgres@localhost:5432/go-rest-api-boilerplate?sslmode=disable" \
  version
```

## Current Schema

### Users Table

**File:** `20251231062946_create_table_users.up.sql`

```sql
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_email ON users(email);
```

**Columns:**
- `id` - Primary key, auto-increment
- `email` - Unique email address
- `password` - Bcrypt hashed password
- `name` - User full name
- `created_at` - Account creation timestamp
- `updated_at` - Last update timestamp

**Indexes:**
- `idx_users_email` - For fast email lookups during login

## Creating New Migrations

### Step 1: Generate Files

```bash
# Using Makefile
make migrate-create NAME=add_user_avatar

# This creates two files:
# migrations/<timestamp>_add_user_avatar.up.sql
# migrations/<timestamp>_add_user_avatar.down.sql
```

### Step 2: Write Up Migration

```sql
-- migrations/<timestamp>_add_user_avatar.up.sql

-- Add avatar column to users table
ALTER TABLE users ADD COLUMN IF NOT EXISTS avatar VARCHAR(500);

-- Create index for avatar lookups
CREATE INDEX IF NOT EXISTS idx_users_avatar ON users(avatar);
```

### Step 3: Write Down Migration

```sql
-- migrations/<timestamp>_add_user_avatar.down.sql

-- Remove avatar index
DROP INDEX IF EXISTS idx_users_avatar;

-- Remove avatar column
ALTER TABLE users DROP COLUMN IF EXISTS avatar;
```

### Step 4: Test Migrations

```bash
# Apply migration
make migrate-up

# Verify schema
docker exec -it go-rest-api-postgres psql -U postgres -d api_db -c "\d users"

# Test rollback
make migrate-down

# Verify rollback
docker exec -it go-rest-api-postgres psql -U postgres -d api_db -c "\d users"

# Reapply
make migrate-up
```

## Database Seeder

Seeder for populating test data during development.

### Running Seeder

```bash
# Run seeder via Makefile
make seed

# Run seeder manually
go run cmd/seeder/main.go
```

### Default Seeded Data

**Test Users:**
- `admin@example.com` / `admin123`
- `user@example.com` / `user123`
- `test@example.com` / `test123`

All passwords are hashed with bcrypt.

### Seeder Behavior

**Fresh Mode:** Truncates tables before inserting

```go
// cmd/seeder/main.go
fresh := true  // Truncate tables first
seeder.SeedAll(ctx, fresh)
```

**Safe Mode:** Skips if data exists

```go
// cmd/seeder/main.go
fresh := false  // Skip if data exists
seeder.SeedAll(ctx, fresh)
```

### Customizing Seed Data

Edit `cmd/seeder/main.go`:

```go
func (s *Seeder) SeedUsers(ctx context.Context) error {
    users := []struct {
        Email    string
        Password string
        Name     string
    }{
        {
            Email:    "admin@example.com",
            Password: "admin123",
            Name:     "Admin User",
        },
        // Add your test users here
        {
            Email:    "custom@example.com",
            Password: "custom123",
            Name:     "Custom User",
        },
    }
    
    // Implementation...
}
```

## Development Workflow

### Initial Setup

```bash
# 1. Start Docker services
make docker-up

# 2. Run migrations
make migrate-up

# 3. Seed database
make seed

# 4. Start application
make run
```

### Adding New Table

```bash
# 1. Create migration
make migrate-create NAME=create_products_table

# 2. Edit up migration
# migrations/<timestamp>_create_products_table.up.sql

# 3. Edit down migration
# migrations/<timestamp>_create_products_table.down.sql

# 4. Apply migration
make migrate-up

# 5. Update seeder (optional)
# cmd/seeder/main.go

# 6. Test
make seed
make run
```

### Resetting Database

```bash
# Option 1: Rollback and reapply
make migrate-down
make migrate-up
make seed

# Option 2: Complete reset (destructive)
make docker-down
docker volume prune
make docker-up
make migrate-up
make seed
```

## Migration Best Practices

### 1. Idempotent Operations

Always use conditional statements:

```sql
-- Good: Safe to run multiple times
CREATE TABLE IF NOT EXISTS users (...);
ALTER TABLE users ADD COLUMN IF NOT EXISTS avatar VARCHAR(500);
DROP TABLE IF EXISTS users;
DROP INDEX IF EXISTS idx_users_avatar;

-- Bad: Fails on second run
CREATE TABLE users (...);
ALTER TABLE users ADD COLUMN avatar VARCHAR(500);
DROP TABLE users;
DROP INDEX idx_users_avatar;
```

### 2. Reversible Migrations

Every up migration must have corresponding down:

```sql
-- Up migration
ALTER TABLE users ADD COLUMN avatar VARCHAR(500);
CREATE INDEX idx_users_avatar ON users(avatar);

-- Down migration (reverse order)
DROP INDEX IF EXISTS idx_users_avatar;
ALTER TABLE users DROP COLUMN IF EXISTS avatar;
```

### 3. Atomic Changes

One logical change per migration:

```bash
# Good: Separate migrations
20260101120000_add_user_avatar.up.sql
20260101120100_add_user_bio.up.sql

# Bad: Multiple unrelated changes
20260101120000_add_user_fields.up.sql  # Contains avatar + bio + phone
```

### 4. Data Migrations

Handle existing data carefully:

```sql
-- Add column with default
ALTER TABLE users ADD COLUMN role VARCHAR(50) DEFAULT 'user';

-- Update existing records
UPDATE users SET role = 'user' WHERE role IS NULL;

-- Make NOT NULL after data is populated
ALTER TABLE users ALTER COLUMN role SET NOT NULL;
```

### 5. Never Modify Existing

Once applied to production, never modify:

```bash
# Good: Create new migration
20260101120000_add_user_avatar.up.sql       # Original
20260101120100_modify_user_avatar.up.sql    # Fix/change

# Bad: Edit existing migration
20260101120000_add_user_avatar.up.sql       # Modified (breaks other environments)
```

## Troubleshooting

### Migration State is Dirty

```bash
# Check migration state
docker exec -it go-rest-api-postgres psql -U postgres -d api_db \
  -c "SELECT version, dirty FROM schema_migrations;"

# If dirty=true, force to last good version
migrate -path migrations \
  -database "postgresql://postgres:postgres@localhost:5432/api_db?sslmode=disable" \
  force <version_number>

# Then retry
make migrate-up
```

### Connection Refused

```bash
# Check if PostgreSQL is running
docker ps | grep postgres

# Check connection details
cat .env | grep DB_

# Restart database
make docker-down
make docker-up
```

### Table Already Exists

```bash
# Check existing tables
docker exec -it go-rest-api-postgres psql -U postgres -d api_db -c "\dt"

# Option 1: Rollback first
make migrate-down
make migrate-up

# Option 2: Drop table manually
docker exec -it go-rest-api-postgres psql -U postgres -d api_db \
  -c "DROP TABLE IF EXISTS users CASCADE;"
make migrate-up
```

### Seeder Fails

```bash
# Check database connection
make docker-up

# Check if migrations ran
docker exec -it go-rest-api-postgres psql -U postgres -d api_db -c "\dt"

# Run migrations if tables don't exist
make migrate-up

# Then retry seeder
make seed
```

## References

- [golang-migrate](https://github.com/golang-migrate/migrate)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [Migration Best Practices](https://www.prisma.io/dataguide/types/relational/migration-best-practices)

