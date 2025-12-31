package response

import (
	"strings"

	"github.com/alpinnz/go-rest-api-boilerplate/internal/domain"
	"github.com/alpinnz/go-rest-api-boilerplate/internal/localization"
)

// getByPath retrieves value from nested dictionary using dot notation path.
// Example: "validation.required" -> dict["validation"]["required"]
func getByPath(dict localization.Dictionary, path string) (string, bool) {
	parts := strings.Split(path, ".")
	var cur any = dict

	for _, p := range parts {
		// Try map[string]any first
		if m, ok := cur.(map[string]any); ok {
			cur, ok = m[p]
			if !ok {
				return "", false
			}
			continue
		}

		// Try map[string]interface{} (from JSON unmarshal)
		if m, ok := cur.(map[string]interface{}); ok {
			cur, ok = m[p]
			if !ok {
				return "", false
			}
			continue
		}

		// Not a map type
		return "", false
	}

	s, ok := cur.(string)
	return s, ok
}

// applyParams replaces placeholders in message with actual values.
// Example: "Password must be at least {min} characters" + {"min": "8"}
// Returns: "Password must be at least 8 characters"
func applyParams(msg string, params map[string]string) string {
	for k, v := range params {
		msg = strings.ReplaceAll(msg, "{"+k+"}", v)
	}
	return msg
}

// TranslateAppError translates AppError code to localized message.
// For validation errors, looks up field label and injects it into message template.
// Falls back to common.error if translation key not found.
func TranslateAppError(dict localization.Dictionary, e domain.AppError) string {
	if msg, ok := getByPath(dict, e.Code); ok {
		// Copy params to avoid modifying original
		params := make(map[string]string)
		for k, v := range e.Params {
			params[k] = v
		}

		// If this is a validation error and has a field, lookup field label
		if e.Field != "" && len(e.Code) > 11 && e.Code[:11] == "validation." {
			fieldLabel := getFieldLabel(dict, e.Field)
			params["field"] = fieldLabel
		}

		return applyParams(msg, params)
	}

	// Fallback to common.error
	if msg, ok := getByPath(dict, "common.error"); ok {
		return msg
	}

	return "error"
}

// getFieldLabel retrieves human-readable field label from dictionary.
// Falls back to the field name itself if label not found.
func getFieldLabel(dict localization.Dictionary, field string) string {
	labelPath := "fields." + field
	if label, ok := getByPath(dict, labelPath); ok {
		return label
	}
	// Fallback: convert snake_case to Title Case
	return snakeToTitle(field)
}

// snakeToTitle converts snake_case to Title Case.
// Example: first_name -> First Name
func snakeToTitle(s string) string {
	if s == "" {
		return ""
	}

	var result []rune
	capitalize := true

	for _, r := range s {
		if r == '_' {
			result = append(result, ' ')
			capitalize = true
			continue
		}

		if capitalize && r >= 'a' && r <= 'z' {
			result = append(result, r-'a'+'A')
			capitalize = false
		} else {
			result = append(result, r)
		}
	}

	return string(result)
}

// TranslateValidationErrors converts slice of AppError to localized FieldError slice.
func TranslateValidationErrors(dict localization.Dictionary, errs []domain.AppError) []FieldError {
	out := make([]FieldError, 0, len(errs))
	for _, e := range errs {
		out = append(out, FieldError{
			Field:   e.Field,
			Code:    e.Code,
			Message: TranslateAppError(dict, e),
		})
	}
	return out
}
