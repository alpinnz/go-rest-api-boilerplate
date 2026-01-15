package validator

import (
	"fmt"

	"github.com/alpinnz/go-rest-api-boilerplate/internal/domain"
	v10 "github.com/go-playground/validator/v10"
)

var validate = v10.New()

// ValidateStruct validates struct fields using validator tags.
// Returns slice of AppError with translation keys for localization support.
// Returns nil if validation passes.
func ValidateStruct(s any) []domain.AppError {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	verrs, ok := err.(v10.ValidationErrors)
	if !ok {
		return []domain.AppError{{Code: "common.error"}}
	}

	out := make([]domain.AppError, 0, len(verrs))
	for _, fe := range verrs {
		field := fe.Field()
		tag := fe.Tag()

		// Build translation key: validation.{tag}
		code := fmt.Sprintf("validation.%s", tag)

		params := map[string]string{}
		if fe.Param() != "" {
			params[tag] = fe.Param()
		}

		out = append(out, domain.AppError{
			Code:   code,
			Field:  toSnakeLower(field),
			Params: params,
		})
	}

	return out
}

// toSnakeLower converts PascalCase to snake_case.
// Example: Email -> email, UserID -> user_id
func toSnakeLower(in string) string {
	if in == "" {
		return ""
	}

	out := make([]rune, 0, len(in)+4)
	for i, r := range in {
		// Add underscore before uppercase (except first char)
		if i > 0 && r >= 'A' && r <= 'Z' {
			out = append(out, '_')
		}

		// Convert to lowercase
		if r >= 'A' && r <= 'Z' {
			r = r - 'A' + 'a'
		}

		out = append(out, r)
	}

	return string(out)
}

// Validate is kept for backward compatibility.
// Use ValidateStruct for new code with localization support.
func Validate(s interface{}) []domain.AppError {
	errList := ValidateStruct(s)
	if errList == nil {
		return []domain.AppError{}
	}
	return errList
}
