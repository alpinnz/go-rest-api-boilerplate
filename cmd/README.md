# cmd/

Application entry points.

## Directory Structure

```
cmd/
├── api/          Main application
└── seeder/       Database seeder
```

## api/

Main REST API application.

### Files

- `main.go` - Application entry point
- `main_test.go` - Main package tests

### Responsibilities

- Load configuration
- Initialize dependencies
- Wire up components
- Start HTTP server
- Handle graceful shutdown
- Automatic port finding if port is in use

### Features

**Automatic Port Finding:**

If the configured port is already in use, the application automatically finds the next available port:

```bash
# If port 8080 is in use:
2025/12/31 15:18:18 Port 8080 is already in use, trying next port...
2025/12/31 15:18:18 Port 8081 is already in use, trying next port...
Server starting on port 8082
```

The application tries up to 10 consecutive ports before failing.

### Usage

```bash
# Run application
make run

# Or directly
go run cmd/api/main.go

# Build binary
make build

# Run binary
./bin/api
```

## seeder/

Database seeder for development.

### Files

- `main.go` - Seeder implementation

### Responsibilities

- Connect to database
- Truncate tables (if fresh seed)
- Insert test data
- Handle errors gracefully

### Usage

```bash
# Run seeder
make seed

# Or directly
go run cmd/seeder/main.go
```

### Default Seeded Data

- Users: admin, user, test accounts
- Passwords: Hashed with bcrypt

See [Database Migrations](../migrations/README.md) for details.

## Adding New Command

Create new directory under `cmd/`:

```bash
# Example: Add CLI tool
mkdir cmd/cli
touch cmd/cli/main.go
```

```go
// cmd/cli/main.go
package main

func main() {
    // CLI implementation
}
```

```bash
# Run new command
go run cmd/cli/main.go
```

## References

- [Application entry point best practices](https://peter.bourgon.org/go-best-practices-2016/#program-design)
- [Go project layout](https://github.com/golang-standards/project-layout)

