# Command Line Applications

Entry points for different applications in the project.

## Structure

```
cmd/
├── api/       # Main REST API server
├── seeder/    # Database seeder utility
└── cli/       # CLI tool (accessed via Makefile)
```

## Applications

### API (cmd/api/main.go)

Main REST API application server.

Run Commands:
```bash
make dev                # With hot reload (recommended)
make run                # Without hot reload
make build              # Build binary to bin/api
```

Configuration:
- Uses environment variables from `.env` file
- APP_PORT, DB_*, REDIS_*, JWT_SECRET

Initialization Flow:
1. Load configuration from environment
2. Setup logger
3. Connect to PostgreSQL database
4. Connect to Redis
5. Initialize repositories, use cases, handlers
6. Setup middleware and routes
7. Start HTTP server
8. Handle graceful shutdown

### Seeder (cmd/seeder/main.go)


Database seeding utility for development and testing.

**Purpose:**
- Populate database with initial data
- Create test users and roles
- Setup user-role relationships

**Run:**
```bash
# Run seeder (includes migration)
make seed

# Direct run
go run cmd/seeder/main.go

# Build and run
make build-seeder
./bin/seeder
```

**Default Data Created:**
- 3 roles: admin, user, moderator
- 3 users with different role combinations
- User-role assignments

**Configuration:**
Uses same database configuration as API from `.env` file.

**Seeder Flow:**
1. Load configuration
2. Connect to database
3. Run migrations (if needed)
4. Truncate existing data (fresh mode)
5. Seed roles
6. Seed users
7. Seed user-role relationships
8. Close connections

**Example seeder structure:**
```go
func main() {
    // Load config
    cfg := config.Load()
    
    // Connect to database
    db, err := database.NewPostgresDB(cfg.Database.DSN())
    if err != nil {
        log.Fatal("Failed to connect", err)
    }
    defer db.Close()
    
    // Initialize seeder runner
    runner := seeder.NewRunner(db)
    
    // Run seeders
    if err := runner.Run(ctx); err != nil {
        log.Fatal("Seeding failed", err)
    }
    
    log.Info("Seeding completed successfully")
}
```

### CLI (cmd/cli/main.go)

CLI tool for code generation, development, testing, and deployment tasks.

Access via Makefile:
```bash
make help               # Show all available commands
make gen-module product # Generate complete module
make dev                # Start development server
make test               # Run tests
make migrate-up         # Run migrations
```

## Building Applications

Build Commands:
```bash
make build              # Build API binary
make build-all          # Build all binaries (API, Seeder, CLI)
```

Binaries created in `bin/` directory.

## Makefile Commands

All development commands accessed through Makefile:

```bash
make help               # Show all available commands
make dev                # Start with hot reload
make start              # Quick start (docker + migrate + dev)
make gen-module name    # Generate complete module
make test               # Run tests
make migrate-up         # Run migrations
make seed               # Seed database
make docker-up          # Start containers
```

See main README.md for complete command reference.

