# Code Generation Templates

Templates for rapid feature development.

## Available Templates

Template files use `//go:build ignore` directive and are excluded from `go fmt` to avoid package conflicts.

```
templates/
├── doc.go                                  # Package documentation
├── entity_template.go                      # Domain entity
├── repository_interface_template.go        # Repository interface
├── repository_implementation_template.go   # Repository implementation
├── usecase_template.go                     # Use case
├── handler_template.go                     # HTTP handler
├── dto_template.go                         # DTO
├── migration_up_template.sql              # Migration up
└── migration_down_template.sql            # Migration down
```

## Why Exclude Templates from gofmt?

- Each template declares different package names (dto, entity, handler, etc.)
- This is intentional - templates generate code for different packages
- Go's gofmt expects one package per directory, but templates need multiple
- Templates are placeholder code, not meant to be compiled directly

Exclusion Methods:
1. `make fmt` - Automatically excludes templates/ directory
2. `.gofmtignore` - Git pre-commit hook exclusion
3. `//go:build ignore` - Prevents compilation

## Placeholder Pattern

Templates use `PLACEHOLDER_` prefix for type-safe placeholders:

```
PLACEHOLDER_Entity   -> Product  (capitalized)
PLACEHOLDER_entity   -> product  (lowercase)
PLACEHOLDER_ENTITY   -> PRODUCT  (uppercase)
PLACEHOLDER_Entities -> Products (plural capitalized)
PLACEHOLDER_entities -> products (plural lowercase)
```

Why PLACEHOLDER_ prefix?
- Valid Go syntax (no IDE errors)
- Type-safe (templates are valid Go code)
- Easy to search and replace
- Clear and explicit

## Usage with Code Generator

Generate code using Makefile:

```bash
# Generate complete module (6 files)
make gen-module product

# Creates:
# - internal/domain/entity/product.go
# - internal/domain/repository/product_repository.go
# - internal/repository/product_repository.go
# - internal/usecase/product_usecase.go
# - internal/delivery/http/dto/product_dto.go
# - internal/delivery/http/handler/product_handler.go

# Generate individual components
make gen-handler user
make gen-repository order
make gen-service payment
make gen-migration create_table_products
```

## Manual Usage

Copy template and replace placeholders:

```bash
# Copy entity template
cp templates/entity_template.go internal/domain/entity/product.go

# Replace placeholders
sed -i 's/PLACEHOLDER_Entity/Product/g' internal/domain/entity/product.go
sed -i 's/PLACEHOLDER_entity/product/g' internal/domain/entity/product.go
```

## Customizing Templates

Edit template files in templates/ directory to match your coding style:

1. Modify structure or add fields
2. Change imports
3. Add custom methods
4. Update documentation

Templates are read by CLI tool during generation.
cp templates/usecase_template.go internal/usecase/product_usecase.go

# Handler
## Manual Usage

For manual customization:

```bash
# Copy template
cp templates/entity_template.go internal/domain/entity/product.go

# Replace placeholders
sed -i 's/PLACEHOLDER_Entity/Product/g' internal/domain/entity/product.go
sed -i 's/PLACEHOLDER_entity/product/g' internal/domain/entity/product.go
```

## Template Structure

Each template includes:
- Standard CRUD operations
- Soft delete support
- Context handling
- Error handling
- Validation tags (for DTOs)

Customize after generation:
- Add feature-specific fields
- Implement business logic
- Add custom methods
- Update validation rules

## Best Practices

Naming Conventions:
- Entity: Singular, PascalCase (User, Product, Order)
- Repository: EntityRepository (UserRepository)
- Use Case: EntityUseCase (UserUseCase)
- Handler: EntityHandler (UserHandler)
- DTO: EntityRequest/EntityResponse (CreateUserRequest)

File Organization:
```
internal/
  domain/entity/product.go
  domain/repository/product_repository.go
  repository/product_repository.go
  usecase/product_usecase.go
  delivery/http/handler/product_handler.go
  delivery/http/dto/product_dto.go
```


