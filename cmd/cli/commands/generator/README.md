# Code Generator

CLI code generator utilities for rapid feature development.

## Structure

```
generator/
├── utils.go      # Generator utility functions
└── templates/    # Code generation templates
    ├── README.md
    ├── doc.go
    ├── *.tmpl    # Template files
    └── templates.go
```

## Purpose

This package provides:
- Template-based code generation
- Placeholder replacement
- File creation utilities
- Feature configuration

## Usage

Used internally by CLI commands in `cmd/cli/commands/gen.go`.

Invoked via Makefile:
```bash
make gen-module product
make gen-handler user
make gen-repository order
```

## Template Files

All template files use `.tmpl` extension to:
- Avoid package conflicts
- Prevent gofmt errors
- Clearly identify as templates

See `templates/README.md` for template documentation.

