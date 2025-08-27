package validation

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

// StrongPasswordValidator ensures the password meets complexity requirements.
var StrongPasswordValidator validator.Func = func(fl validator.FieldLevel) bool {
	pass := fl.Field().String()
	if pass == "" {
		return true
	}
	return regexp.MustCompile(`[A-Z]`).MatchString(pass) &&
		regexp.MustCompile(`[a-z]`).MatchString(pass) &&
		regexp.MustCompile(`[0-9]`).MatchString(pass) &&
		regexp.MustCompile(`[#!?@$%^&*-]`).MatchString(pass) &&
		len(pass) >= 8
}
