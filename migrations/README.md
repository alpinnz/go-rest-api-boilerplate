# Database Migrations

This directory contains SQL migration files for database schema management.

## Structure

Migrations are organized in pairs:
- `{timestamp}_{name}.up.sql` - Apply migration
- `{timestamp}_{name}.down.sql` - Rollback migration

## Current Migrations

```
20260102053510_create_tabel_users.up.sql        # Create users table
20260102053510_create_tabel_users.down.sql      # Drop users table

20260102053516_create_tabel_roles.up.sql        # Create roles table
20260102053516_create_tabel_roles.down.sql      # Drop roles table

20260102053628_create_tabel_user_roles.up.sql   # Create user_roles junction table
20260102053628_create_tabel_user_roles.down.sql # Drop user_roles table
```

## Schema Overview

### Users Table

Stores user account information with soft delete support.

**Columns:**
- `id` (UUID, Primary Key)
- `email` (VARCHAR, Unique, Not Null)
- `password` (VARCHAR, Not Null) - bcrypt hashed
- `first_name` (VARCHAR)
- `last_name` (VARCHAR)
- `created_at` (TIMESTAMP)
- `updated_at` (TIMESTAMP)
- `deleted_at` (TIMESTAMP) - NULL for active users

**Indexes:**
- Primary key on `id`
- Unique index on `email`
- Index on `deleted_at` for soft delete queries

### Roles Table

Stores role definitions with soft delete support.

**Columns:**
- `id` (UUID, Primary Key)
- `name` (VARCHAR, Unique, Not Null)
- `description` (TEXT)
- `created_at` (TIMESTAMP)
- `updated_at` (TIMESTAMP)
- `deleted_at` (TIMESTAMP) - NULL for active roles

**Indexes:**
- Primary key on `id`
- Unique index on `name`
- Index on `deleted_at` for soft delete queries

### User Roles Table

Junction table for many-to-many relationship between users and roles.

**Columns:**
- `user_id` (UUID, Foreign Key -> users.id)
- `role_id` (UUID, Foreign Key -> roles.id)
- `created_at` (TIMESTAMP)

**Constraints:**
- Composite primary key on (`user_id`, `role_id`)
- Foreign key to users table with CASCADE delete
- Foreign key to roles table with CASCADE delete

## Running Migrations

### Prerequisites

Install development tools using CLI:
```bash
# Build CLI tool first
go build -o app cmd/cli/main.go

# Install tools (includes golang-migrate)
make install tools
```

### Apply Migrations

Run all pending migrations:
```bash
make migrate up
```

This connects to database using configuration from `.env` file.

### Rollback Migrations

Rollback all migrations:
```

## Migration Guidelines
- Always create both up and down migrations
- Test rollback before committing
- Use transactions for data migrations
- Keep migrations small and focused
- Never modify existing migrations in production

Naming Convention:
- Use descriptive names: create_table_orders, add_user_status
- Use snake_case for consistency
- Be specific: add_email_index not add_index

Common Patterns:
```sql
-- Add Column
ALTER TABLE users ADD COLUMN phone VARCHAR(20);

-- Add Index
CREATE INDEX idx_users_email ON users(email);

-- Create Table
CREATE TABLE profiles (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Add Foreign Key
ALTER TABLE orders ADD CONSTRAINT fk_user 
FOREIGN KEY (user_id) REFERENCES users(id);
```

## Troubleshooting

Migration Failed:
1. Check database state
2. Fix migration file
3. Force to previous version: `make migrate-down`
4. Re-run: `make migrate-up`

Dirty State:
1. Manually fix database if needed
2. Force clean state with last good version

