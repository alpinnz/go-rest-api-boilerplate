package validation

import (
	"reflect"

	"github.com/go-playground/validator/v10"
)

// NewKindValidator returns a validator function that checks if a field has a specific reflect.Kind.
func NewKindValidator(kind reflect.Kind) validator.Func {
	return func(fl validator.FieldLevel) bool {
		return fl.Field().Kind() == kind
	}
}
