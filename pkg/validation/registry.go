package validation

import (
	"reflect"

	"github.com/go-playground/validator/v10"
)

// RegisterAll registers all custom validation rules to the validator instance.
func RegisterAll(v *validator.Validate) error {
	// Register regex validators
	regexRules := map[string]string{
		"regex_name":  `^[\w\s-]*$`,
		"regex_uuid7": `^[0-9a-f]{8}-[0-9a-f]{4}-7[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`,
	}
	for tag, pattern := range regexRules {
		if err := v.RegisterValidation(tag, NewRegexValidator(pattern)); err != nil {
			return err
		}
	}

	// Register kind validators
	kindRules := map[string]reflect.Kind{
		"kind_string": reflect.String,
		"kind_bool":   reflect.Bool,
		"kind_int":    reflect.Int,
		"kind_array":  reflect.Array,
	}
	for tag, kind := range kindRules {
		if err := v.RegisterValidation(tag, NewKindValidator(kind)); err != nil {
			return err
		}
	}

	// Register password validator
	if err := v.RegisterValidation("strong_password", StrongPasswordValidator); err != nil {
		return err
	}

	return nil
}
