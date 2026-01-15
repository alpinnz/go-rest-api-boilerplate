# Go REST API Boilerplate

Production-ready REST API boilerplate built with Go, following Clean Architecture, SOLID, and DDD.

> Note: Folder-level documentation lives in:
> - `cmd/README.md` (binaries)
> - `internal/README.md` (architecture/layers)
> - `pkg/README.md` (reusable packages)
> - `internal/websocket/README.md` (WebSocket details)

## Table of Contents

- [Tech Stack](#tech-stack)
- [Features](#features)
- [Prerequisites](#prerequisites)
- [Quick Start](#quick-start)
- [Project Structure](#project-structure)
- [Development Guide](#development-guide)
- [API Endpoints](#api-endpoints)

## Tech Stack

- **Go 1.24+** - Programming language
- **Gin** - HTTP web framework
- **go-playground/validator** - Struct validation
- **PostgreSQL 15** - Primary database
- **Redis 7** - Session storage and caching
- **Docker** - API containerization
- **JWT** - Authentication tokens
- **Air** - Hot reload for development
- **golangci-lint** - Code linting
- **Swag** - OpenAPI/Swagger generation

## Features

### Architecture & Design
- Clean Architecture (layered separation)
- Dependency Injection (centralized container)
- Domain-Driven Design (entities + domain errors)
- Repository pattern (interfaces in domain, implementations in `internal/repository`)

### Platform Features
- Swagger UI at `/docs`
- JWT auth (access + refresh) with Redis-backed sessions
- RBAC (roles + user-role)
- Localization (EN/ID)
- WebSocket server (see `internal/websocket/README.md`)

### Production Readiness
- Structured logging (request_id propagation)
- Graceful shutdown with cleanup
- Configurable CORS + rate limiting
- Health check endpoints (live/ready)
- Connection pooling
- Built-in migrations runner

## Prerequisites

- Go 1.24+
- PostgreSQL
- Redis

## Quick Start

```bash
make install

# Setup environment
if [ -f .env.example ]; then cp .env.example .env; else echo "create a .env file with required environment variables"; fi

# Start everything (Docker + Migrations + Dev server)
make start

curl http://localhost:8080/api/v1/health
```

## Project Structure

High level overview (details in `internal/README.md` and `pkg/README.md`):

```
cmd/                Application entry points (see cmd/README.md)
config/             Configuration management
internal/           Application code (see internal/README.md)
pkg/                Reusable packages (see pkg/README.md)
migrations/         SQL migration files
```

Design decisions (project-wide):
- Prefer interfaces in `internal/domain/repository`, implementations in `internal/repository`
- No hardcoded config values (use env + validation)
- Consistent request tracing via `request_id` in logs

## Development Guide

Use `make help` to list all commands.

Common commands:
- `make dev` (hot reload)
- `make test`
- `make lint`
- `make migrate-up`
- `make seed`

## API Endpoints

### Public
- `GET  /api/v1/health`
- `GET  /api/v1/health/live`
- `GET  /api/v1/health/ready`
- `GET  /docs`

(See swagger for full list.)
