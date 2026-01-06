# Code Generation Templates

Templates for rapid feature development via CLI generator.

## Location

These templates are embedded within the CLI tool at:
```
cmd/cli/commands/generator/templates/
```

## Available Templates

```
templates/
├── doc.go                                  # Package documentation
├── entity_template.go.tmpl                 # Domain entity
├── repository_interface_template.go.tmpl   # Repository interface
├── repository_implementation_template.go.tmpl # Repository implementation
├── usecase_template.go.tmpl                # Use case
├── handler_template.go.tmpl                # HTTP handler
├── dto_template.go.tmpl                    # DTO
├── migration_up_template.sql              # Migration up
├── migration_down_template.sql            # Migration down
└── templates.go                            # Template exports
```

## Usage

Templates are used automatically by the CLI generator via Makefile:

```bash
# Generate complete module (6 files)
make gen-module name=product

# Creates:
# - internal/domain/entity/product.go
# - internal/domain/repository/product_repository.go
# - internal/repository/product_repository.go
# - internal/usecase/product_usecase.go
# - internal/delivery/http/dto/product_dto.go
# - internal/delivery/http/handler/product_handler.go

# Generate individual components
make gen-handler name=user
make gen-repository name=order
make gen-usecase name=payment
make gen-migration name=create_table_products
```

## Placeholder Pattern

Templates use `PLACEHOLDER_` prefix for type-safe placeholders:

```
PLACEHOLDER_Entity   -> Product  (capitalized)
PLACEHOLDER_entity   -> product  (lowercase)
PLACEHOLDER_ENTITY   -> PRODUCT  (uppercase)
PLACEHOLDER_Entities -> Products (plural capitalized)
PLACEHOLDER_entities -> products (plural lowercase)
```

## Customizing Templates

To customize templates for your project:

1. Edit template files in this directory
2. Modify structure or add fields
3. Change imports or add custom methods
4. Update documentation

Templates are read by the CLI tool during code generation.

## Template Files

Template files use `.tmpl` extension to:
- Prevent package conflicts (templates contain placeholder code)
- Avoid compilation errors during go build
- Clearly identify as template files, not source code
- Enable proper embedding via `//go:embed` directives

The CLI generator reads `.tmpl` files and generates proper `.go` files in their respective directories.

## Best Practices

**Naming Conventions:**
- Entity: Singular, PascalCase (User, Product, Order)
- Repository: EntityRepository (UserRepository)
- Use Case: EntityUseCase (UserUseCase)
- Handler: EntityHandler (UserHandler)
- DTO: EntityRequest/EntityResponse (CreateUserRequest)

**File Organization:**
```
internal/
├── domain/
│   ├── entity/product.go
│   └── repository/product_repository.go
├── repository/product_repository.go
├── usecase/product_usecase.go
└── delivery/http/
    ├── dto/product_dto.go
    └── handler/product_handler.go
```

