# Command Line Applications

Entry points (binaries) for this project. Most day-to-day commands are exposed via the `Makefile`.

## Structure

```
cmd/
├── api/       # Main REST API server
├── seeder/    # Database seeder utility
└── cli/       # CLI tool (used by Makefile)
```

## API (cmd/api/main.go)

Main REST API server entry point.

**Run (recommended):**
- `make dev` (hot reload)
- `make run` (no hot reload)
- `make build` (build to `bin/api`)

**Initialization (high level):**
1. Load configuration from environment (`.env`)
2. Build the dependency injection container (`internal/container`)
3. Setup HTTP router (handlers + middleware)
4. Start the HTTP server
5. Graceful shutdown and resource cleanup

Notes:
- Dependency wiring is centralized in `internal/container` (avoid manual wiring in `main.go`).
- For environment variables and defaults, keep a single source of truth in the root `README.md`.

## Seeder (cmd/seeder/main.go)

Database seeding utility for development/testing.

**Purpose:**
- Populate database with initial data
- Create users/roles
- Setup user-role relationships

**Run (preferred):**
- `make seed`

**What it does (high level):**
- Load configuration
- Connect to database
- Run migrations (if needed)
- Seed data

## CLI (cmd/cli/main.go)

CLI tool used by the `Makefile` for development, testing, code generation, and database tasks.

Use `make help` to discover available commands.

## Building Applications

Binaries are created in `bin/`.

See root `README.md` for the full command reference (development, migrations, docker, tests).
