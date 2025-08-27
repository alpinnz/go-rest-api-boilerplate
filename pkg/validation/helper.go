package validation

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// msg is a helper for rules that only depend on the field name (no param needed).
// Example: "required" -> "email is required"
func msg(tpl string) func(fe validator.FieldError, s interface{}) string {
	return func(fe validator.FieldError, s interface{}) string {
		return fmt.Sprintf(tpl, fmt.Sprintf("'%s'", jsonFieldName(fe, s)))
	}
}

// msgWithParam is a helper for rules that require an extra parameter from validator.
// Example: "gte=18" -> "age must be greater than or equal to 18"
func msgWithParam(tpl string) func(fe validator.FieldError, s interface{}) string {
	return func(fe validator.FieldError, s interface{}) string {
		return fmt.Sprintf(tpl, fmt.Sprintf("'%s'", jsonFieldName(fe, s)), fe.Param())
	}
}
