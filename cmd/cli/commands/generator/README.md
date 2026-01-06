# Code Generator

CLI code generator utilities for rapid feature development.

## Structure

```
generator/
├── utils.go         # Generator utility functions
└── templates/       # Code generation templates
    ├── README.md
    ├── doc.go
    ├── templates.go # Embedded template exports
    └── *.tmpl       # Template files
```

## Purpose

This package provides the core generator functionality:
- Template-based code generation
- Placeholder replacement engine
- File creation and validation
- Feature configuration management

## Usage

Generator is used internally by CLI commands in `cmd/cli/commands/gen.go`.

**Via Makefile:**
```bash
make gen-module name=product       # Generate complete module
make gen-handler name=user         # Generate HTTP handler
make gen-repository name=order     # Generate repository
make gen-usecase name=payment      # Generate use case
make gen-migration name=add_column # Generate migration
```

## Template System

All template files use `.tmpl` extension to:
- Prevent package conflicts
- Avoid compilation errors
- Enable embedding via `//go:embed`

See [templates/README.md](templates/README.md) for detailed template documentation.

## Utility Functions

**Available utilities:**
- `FileExists(path)` - Check if file exists
- `WriteFromTemplate(template, dst)` - Write template to file
- `ReplaceInFile(path, config)` - Replace placeholders
- `NewFeatureConfig(name)` - Create configuration from name
