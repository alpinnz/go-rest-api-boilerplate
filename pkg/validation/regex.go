package validation

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

// NewRegexValidator returns a validator function that matches a regex pattern.
func NewRegexValidator(pattern string) validator.Func {
	re := regexp.MustCompile(pattern)
	return func(fl validator.FieldLevel) bool {
		value := fl.Field().String()
		if value == "" {
			return true
		}
		return re.MatchString(value)
	}
}
