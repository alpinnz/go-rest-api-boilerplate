# API Documentation

OpenAPI/Swagger documentation files.

## Files

- `docs.go` - auto-generated Go code (gitignored; regenerate with `make swagger`)
- `swagger.json` - OpenAPI spec (gitignored; auto-generated)
- `swagger.yaml` - YAML spec (gitignored; auto-generated)
- `swagger.html` - Swagger UI (static file; committed)

**Note:** Only `swagger.html` should be committed. All other files are generated.

## Access

When the API is running, access:
- Swagger UI: `http://localhost:8080/docs`
- OpenAPI spec: `http://localhost:8080/docs/swagger.json`

If you change the server host/port, update the base URL accordingly.

## Auto Token Management

The Swagger UI (`swagger.html`) includes automatic token management for Bearer authorization.

Features:
- Auto-save tokens after login/register
- Auto-fill Swagger's "Authorize" with `Bearer {access_token}`
- Auto-inject `Authorization: Bearer ...` for protected endpoints
- Auto-fill refresh_token input for refresh-token endpoint
- Auto-clear tokens on logout (localStorage + Swagger auth)
- Helper buttons to show/clear tokens

## Generation

The OpenAPI spec is generated from code annotations.

Regenerate using Makefile (recommended):
```bash
make swagger
```

Or manually:
```bash
swag init -g cmd/api/main.go -o internal/delivery/http/docs --parseDependency --parseInternal
```

## Annotations

Docs are generated from comments in:
- `cmd/api/main.go` (general API info)
- `internal/delivery/http/handler/*.go` (endpoints)

## Best Practices

- Regenerate docs after changing endpoints
- Keep annotations aligned with behavior
- Do not commit generated files (`docs.go`, `swagger.json`, `swagger.yaml`)
