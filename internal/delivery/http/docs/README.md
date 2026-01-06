# API Documentation

OpenAPI/Swagger documentation files.

## Files

- `docs.go` - Auto-generated Go code (gitignored, regenerate with `make swagger`)
- `swagger.json` - OpenAPI 2.0 specification (gitignored, auto-generated)
- `swagger.yaml` - YAML format specification (gitignored, auto-generated)
- `swagger.html` - Swagger UI interface (static file, committed to repo)

**Note:** Only `swagger.html` should be committed to the repository. All other files are auto-generated and should be regenerated in each environment.

## Access

When the API is running, access the documentation at:

- Swagger UI: http://localhost:8080/docs
- OpenAPI Spec: http://localhost:8080/docs/swagger.json

### Auto Token Management

The Swagger UI includes automatic token management with Bearer Token authorization:

**Features:**
- Auto-save tokens after login/register
- Auto-authorize Swagger's built-in "Authorize" button with Bearer token
- Auto-inject access token in Authorization header (Bearer {access_token}) for protected endpoints
- Auto-fill refresh_token field in refresh-token endpoint
- Auto-clear tokens on logout from both localStorage and Swagger authorization
- Status indicator showing authentication state
- Helper buttons to show/clear tokens

**How It Works:**
- All protected endpoints use `BearerAuth` security scheme
- After login, access_token is automatically set in Swagger authorization
- All requests to protected endpoints include: `Authorization: Bearer {access_token}`
- The 🔒 icon on endpoints indicates Bearer Token is required

**Workflow:**
1. Use `/auth/login` or `/auth/register` endpoint
2. Tokens automatically saved to localStorage
3. Swagger "Authorize" button automatically filled with Bearer token
4. Test protected endpoints - token automatically included
5. Check "Authorize" button (🔓 → 🔒) to confirm authorization
6. Use "Show Tokens" button to view saved tokens
7. Use "Clear Tokens" button to logout/reset

## Generation

Only `swagger.json` is automatically generated from code annotations using Swag.

To regenerate:
```bash
# Using Makefile (recommended)
make swagger

# Or manually
swag init -g cmd/api/main.go -o internal/delivery/http/docs --parseDependency --parseInternal
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

## Best Practices

- Regenerate documentation after adding or modifying API endpoints
- Keep annotations up to date with implementation
- Test documentation in Swagger UI before committing
- Only commit `swagger.html` to repository
- Auto-generated files (`docs.go`, `swagger.json`, `swagger.yaml`) should be gitignored
- Each environment should regenerate documentation with `make swagger`

