# API Documentation

OpenAPI/Swagger documentation files.

## Files

- `swagger.json` - OpenAPI 2.0 specification (auto-generated)
- `swagger.html` - Swagger UI interface (static file)

## Access

When the API is running, access the documentation at:

- Swagger UI: http://localhost:8080/docs
- OpenAPI Spec: http://localhost:8080/docs/swagger.json

## Generation

Only `swagger.json` is automatically generated from code annotations using Swag.

To regenerate:
```bash
# Install swag if not already installed
go install github.com/swaggo/swag/cmd/swag@latest

# Generate swagger.json
swag init -g cmd/api/main.go -o internal/delivery/http/docs
```

**Note:** `swagger.html` is a static HTML file that loads Swagger UI from CDN and should be committed to the repository.

## Annotations

API documentation is generated from comments in:
- `cmd/api/main.go` - General API info
- `internal/delivery/http/handler/*.go` - Endpoint documentation

Example annotation:
```go
// @Summary User login
// @Description Authenticate user and return JWT token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login credentials"
// @Success 200 {object} dto.LoginResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
    // Implementation
}
```

## Note

- `swagger.json` is gitignored as it's a generated file and should be regenerated in each environment
- `swagger.html` is committed to the repository as it's a static UI wrapper

