# Go REST API Boilerplate

Production-ready REST API boilerplate built with Go, following Clean Architecture and industry best practices.

## Features

### Core
- Clean Architecture with dependency inversion
- JWT-based authentication with Redis sessions
- User registration & login
- Password hashing with bcrypt
- RESTful API design
- Health check endpoint

### Database
- PostgreSQL for persistent data
- Redis for session storage & caching
- Migration system with versioning
- Database seeder for development

### Security
- JWT token authentication
- Session management via Redis
- Password hashing (bcrypt)
- Session invalidation on logout
- CORS middleware
- Input validation with go-playground/validator
- SQL injection prevention

### Localization
- Multi-language support (English, Indonesian)
- Accept-Language header detection
- Code-based error translation
- Parameter interpolation
- Extensible language system

### Development
- Hot reload ready
- Docker Compose for local development
- Make commands for all operations
- Linter configuration (golangci-lint)
- Environment-based configuration
- OpenAPI/Swagger JSON auto-generation

### Production Ready
- Graceful shutdown
- Connection pooling
- Structured logging middleware
- Recovery middleware (panic handling)
- Docker support

## Tech Stack

- **Go 1.21+** - Programming language
- **Gin** - HTTP web framework
- **go-playground/validator** - Struct validation
- **PostgreSQL 15** - Primary database
- **Redis 7** - Session storage & caching
- **Docker Compose** - Container orchestration
- **JWT** - Authentication tokens

## Quick Start

```bash
# 1. Setup environment
cp .env.example .env

# 2. Start services
make docker-up

# 3. Run migrations
make migrate-up

# 4. Seed database (optional)
make seed

# 5. Run application
make run

# 6. Test endpoints
curl http://localhost:8080/api/v1/health
curl http://localhost:8080/docs/swagger.json
```

## Using This Boilerplate

### For New Project

```bash
# Clone and setup
git clone <this-repo> my-new-project
cd my-new-project
rm -rf .git
git init

# Update module name in go.mod
# Update container names in docker-compose.yml
# Update import paths in all Go files
```

### Project Structure

```
cmd/
  api/              Application entry point
  seeder/           Database seeder
config/             Configuration management
internal/
  domain/           Business entities & repository interfaces
  usecase/          Business logic implementation
  delivery/http/    HTTP handlers, routes & middleware
  repository/       Data access layer (by domain)
  infrastructure/   Database connections (PostgreSQL, Redis)
  middleware/       HTTP middleware (auth, cors, logger, recovery)
pkg/                Public reusable packages
migrations/         SQL migration files
```


## Commands

```bash
# Application
make run           # Run application
make build         # Build binary
make test          # Run tests
make clean         # Clean artifacts

# Docker
make docker-up     # Start Docker services
make docker-down   # Stop Docker services

# Database
make migrate-up    # Run migrations
make migrate-down  # Rollback migrations
make seed          # Run database seeder

# Code Quality
make fmt           # Format code
make vet           # Vet code
make lint          # Run linter
```

## API Endpoints

### Public
- `GET  /api/v1/health` - Health check
- `GET  /docs/swagger.json` - OpenAPI specification
- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/login` - User login (returns access_token & refresh_token)
- `POST /api/v1/auth/refresh` - Refresh access token using refresh_token

### Protected (Requires JWT Access Token)
- `POST /api/v1/auth/logout` - User logout
- `GET  /api/v1/users/me` - Get current user profile
- `GET  /api/v1/users/:id` - Get user by ID

## Authentication

This API uses JWT-based authentication with access and refresh tokens for optimal security and user experience.

### Token Types

**Access Token**
- Lifetime: 15 minutes
- Used for API requests
- Sent in `Authorization: Bearer <token>` header
- Short-lived to minimize risk if compromised

**Refresh Token**
- Lifetime: 7 days
- Used to obtain new access tokens
- Sent in request body to `/auth/refresh` endpoint
- Single-use (token rotation for security)

### How It Works

1. **Login**: User provides credentials, receives both tokens
2. **API Request**: Client uses access token in Authorization header
3. **Token Expired**: When 401 received, use refresh token to get new pair
4. **Logout**: Both tokens invalidated immediately via Redis

### Login Response
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": 1,
    "email": "user@example.com",
    "name": "John Doe"
  }
}
```

### Client Implementation

```javascript
class AuthClient {
  async apiRequest(url, options = {}) {
    // Try with current access token
    let response = await fetch(url, {
      ...options,
      headers: {
        ...options.headers,
        'Authorization': `Bearer ${this.accessToken}`
      }
    })

    // If 401, refresh and retry
    if (response.status === 401) {
      const refreshed = await this.refreshToken()
      if (refreshed) {
        response = await fetch(url, {
          ...options,
          headers: {
            ...options.headers,
            'Authorization': `Bearer ${this.accessToken}`
          }
        })
      }
    }

    return response
  }

  async refreshToken() {
    const response = await fetch('/api/v1/auth/refresh', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ refresh_token: this.refreshToken })
    })

    if (response.ok) {
      const { data } = await response.json()
      this.accessToken = data.access_token
      this.refreshToken = data.refresh_token
      return true
    }

    // Refresh failed, redirect to login
    window.location.href = '/login'
    return false
  }
}
```

### Security Features

- **Token Rotation**: Refresh tokens are single-use, new pair generated each refresh
- **Short-Lived Access**: 15-minute expiration minimizes exposure
- **Server-Side Sessions**: All tokens stored in Redis, can be immediately revoked
- **Token Type Validation**: Access tokens can't be used as refresh tokens
- **HTTPS Required**: Always use HTTPS in production

### Configuration

Edit `pkg/auth/jwt.go` to customize token expiration:

```go
const (
    AccessTokenExpiration  = 15 * time.Minute  // Change as needed
    RefreshTokenExpiration = 7 * 24 * time.Hour // Change as needed
)
```

### Redis Storage

```
Access Token:  session:<token> -> user_id (TTL: 15 min)
Refresh Token: refresh:<token> -> access_token (TTL: 7 days)
```

### Usage Examples
```bash
# Use access token in Authorization header
curl -H "Authorization: Bearer <access_token>" \
     http://localhost:8080/api/v1/users/me

# Refresh access token when expired
curl -X POST http://localhost:8080/api/v1/auth/refresh \
     -H "Content-Type: application/json" \
     -d '{"refresh_token": "<refresh_token>"}'
```

### Troubleshooting

**Access Token Expired (401)**
- Use refresh token to get new access token

**Refresh Token Expired (401)**
- User must login again (7 days passed)

**Both Tokens Invalid**
- User logged out or server restarted (Redis cleared)
- User must login again

## Localization

Multi-language support with automatic detection from `Accept-Language` header.

### Supported Languages
- English (en)
- Indonesian (id)

### Usage

```bash
# English response
curl -H "Accept-Language: en" http://localhost:8080/api/v1/auth/login

# Indonesian response
curl -H "Accept-Language: id" http://localhost:8080/api/v1/auth/login
```

### Adding New Language

1. Create new JSON file in `internal/localization/lang/`:
```bash
cp internal/localization/lang/en.json internal/localization/lang/fr.json
```

2. Translate all keys in the new file

3. Load in router (automatic via middleware)

See [internal/localization/README.md](internal/localization/README.md) for details.

## Documentation

- [OpenAPI Specification](docs/README.md) - Auto-generated `swagger.json`
- [cmd/](cmd/README.md) - Application entry points
- [internal/](internal/README.md) - Application code
- [internal/domain/](internal/domain/README.md) - Domain entities and errors
- [internal/localization/](internal/localization/README.md) - Multi-language support
- [internal/middleware/](internal/middleware/README.md) - HTTP middleware
- [pkg/](pkg/README.md) - Public packages
- [pkg/auth/](pkg/auth/README.md) - Authentication utilities
- [pkg/response/](pkg/response/README.md) - Response formatting
- [pkg/validator/](pkg/validator/README.md) - Input validation
- [migrations/](migrations/README.md) - Database migrations

## Environment Variables

```bash
# Application
APP_PORT=8080

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=go-rest-api-boilerplate
DB_SSLMODE=disable

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# JWT
JWT_SECRET=change-this-secret-key
JWT_EXPIRATION=24h

# Server
READ_TIMEOUT=10s
WRITE_TIMEOUT=10s
SHUTDOWN_TIMEOUT=5s
```

## Session Management

Sessions are stored in Redis for fast access and automatic expiration.

**Flow:**
1. **Login**: Generate JWT token, store in Redis with user ID
2. **Request**: Validate token from Redis, extract user ID
3. **Logout**: Delete token from Redis (immediate invalidation)
4. **Expire**: Auto-expire based on JWT_EXPIRATION

**Why Redis Only:**
- Fast in-memory access
- Built-in TTL (time-to-live)
- Automatic expiration
- Easier horizontal scaling
- No need for cleanup jobs

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for:
- Development workflow
- Branch naming conventions
- Commit message format
- Code style guidelines
- Testing requirements

## License

MIT License - see [LICENSE](LICENSE) file

## Support

- **Issues**: [GitHub Issues](https://github.com/username/repo/issues)
- **Discussions**: [GitHub Discussions](https://github.com/username/repo/discussions)

