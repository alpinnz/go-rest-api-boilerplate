# API Documentation

OpenAPI specification for the REST API.

## OpenAPI Specification File

- **Location**: `docs/swagger.json`
- **HTTP Endpoint**: `GET /docs/swagger.json`
- **Format**: OpenAPI 3.0.3
- **Maintenance**: Can be auto-generated from code annotations or manually edited

## Access Specification

### Via HTTP (When API Running)

```bash
# Direct access
curl http://localhost:8080/docs/swagger.json

# Import to Swagger Editor
https://editor.swagger.io/?url=http://localhost:8080/docs/swagger.json
```

### Via File System

Access the file directly at `docs/swagger.json` for offline use.

## Auto-Generate from Code

Try to generate `docs/swagger.json` from code annotations:

```bash
make swagger-generate
```

This uses Swaggo to scan annotations in:
- `cmd/api/main.go` - API metadata
- Handler files - Endpoint definitions

**Note**: If generation fails, you can edit `docs/swagger.json` manually. The file follows standard OpenAPI 3.0.3 format.

## Features

- Complete endpoint documentation
- Request/response schemas
- Authentication specification (Bearer token)
- Multi-language examples (English, Indonesian)
- Field validation rules
- Error response formats

## Usage

### View with Swagger Editor

1. Go to https://editor.swagger.io
2. File > Import File
3. Select `docs/swagger.json`
4. View interactive documentation

### Import to Postman

1. Open Postman
2. Import > Upload Files
3. Select `docs/swagger.json`
4. Collection auto-created with all endpoints

### Import to Insomnia

1. Open Insomnia
2. Application > Preferences > Data > Import Data
3. Select `docs/swagger.json`
4. All endpoints imported

### Generate Client SDK

Use OpenAPI Generator:

```bash
# TypeScript/Axios
npx @openapitools/openapi-generator-cli generate \
  -i docs/swagger.json \
  -g typescript-axios \
  -o client/typescript

# Python
openapi-generator-cli generate \
  -i docs/swagger.json \
  -g python \
  -o client/python

# Go
openapi-generator-cli generate \
  -i docs/swagger.json \
  -g go \
  -o client/go
```

### Validate Specification

```bash
npx @apidevtools/swagger-cli validate docs/swagger.json
```

## Documented Endpoints

### Public Endpoints

- `GET /api/v1/health` - Health check
- `POST /api/v1/auth/register` - Register new user
- `POST /api/v1/auth/login` - User login

### Protected Endpoints (Require Bearer Token)

- `POST /api/v1/auth/logout` - User logout
- `GET /api/v1/users/me` - Get current user profile
- `GET /api/v1/users/{id}` - Get user by ID

## Request Headers

### Accept-Language

Optional header for localized responses:

```
Accept-Language: en
Accept-Language: id
```

Supported languages:
- `en` - English (default)
- `id` - Indonesian (Bahasa Indonesia)

### Authorization

Required for protected endpoints:

```
Authorization: Bearer <jwt-token>
```

## Response Format

All responses follow consistent structure:

```json
{
  "code": null,
  "message": "Success",
  "data": { ... },
  "errors": null,
  "meta": null,
  "request_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

### Success Response

```json
{
  "code": null,
  "message": "Success",
  "data": {
    "id": 1,
    "email": "user@example.com",
    "name": "John Doe"
  },
  "errors": null,
  "meta": null,
  "request_id": "uuid"
}
```

### Validation Error

```json
{
  "code": "VALIDATION_ERROR",
  "message": "Validation error",
  "data": null,
  "errors": {
    "fields": [
      {
        "field": "email",
        "code": "validation.required",
        "message": "Email is required"
      }
    ]
  },
  "meta": null,
  "request_id": "uuid"
}
```

## Maintenance

When adding/modifying endpoints:

1. Update handler implementation
2. Update `docs/swagger.json` with new endpoint spec
3. Test using Swagger Editor or Postman
4. Validate specification: `npx @apidevtools/swagger-cli validate docs/swagger.json`
5. Commit changes

## Best Practices

1. **Keep swagger.json in sync** - Update when API changes
2. **Test changes** - Validate spec before committing
3. **Document all fields** - Include descriptions and examples
4. **Version API** - Use `/api/v1` prefix
5. **Document errors** - Include all possible error responses
6. **Add examples** - Show request/response samples

## Tools

- [Swagger Editor](https://editor.swagger.io/) - View and edit spec
- [Postman](https://www.postman.com/) - API testing
- [Insomnia](https://insomnia.rest/) - API testing
- [OpenAPI Generator](https://openapi-generator.tech/) - Generate client SDKs
- [Swagger CLI](https://apitools.dev/swagger-cli/) - Validate specification

## External Resources

- [OpenAPI Specification](https://swagger.io/specification/)
- [OpenAPI Guide](https://swagger.io/docs/specification/about/)
- [OpenAPI Examples](https://github.com/OAI/OpenAPI-Specification/tree/main/examples)

