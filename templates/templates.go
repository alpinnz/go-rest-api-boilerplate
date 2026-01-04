package templates

import _ "embed"

// Embedded Go templates
//
//go:embed dto_template.go.tmpl
var DTOTemplate string

//go:embed entity_template.go.tmpl
var EntityTemplate string

//go:embed handler_template.go.tmpl
var HandlerTemplate string

//go:embed repository_interface_template.go.tmpl
var RepositoryInterfaceTemplate string

//go:embed repository_implementation_template.go.tmpl
var RepositoryImplementationTemplate string

//go:embed usecase_template.go.tmpl
var UsecaseTemplate string

// Embedded SQL templates
//
//go:embed migration_up_template.sql
var MigrationUpTemplate string

//go:embed migration_down_template.sql
var MigrationDownTemplate string
