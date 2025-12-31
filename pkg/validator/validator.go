package validator

import (
	"errors"
	"reflect"
	"regexp"
	"strings"
)

func Validate(s interface{}) error {
	v := reflect.ValueOf(s)
	t := reflect.TypeOf(s)

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)
		tag := fieldType.Tag.Get("validate")

		if tag == "" {
			continue
		}

		rules := strings.Split(tag, ",")
		for _, rule := range rules {
			if err := validateRule(field, fieldType.Name, rule); err != nil {
				return err
			}
		}
	}

	return nil
}

func validateRule(field reflect.Value, fieldName, rule string) error {
	parts := strings.Split(rule, "=")
	ruleName := parts[0]

	switch ruleName {
	case "required":
		if isZero(field) {
			return errors.New(fieldName + " is required")
		}
	case "email":
		if field.Kind() == reflect.String {
			email := field.String()
			emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
			if !emailRegex.MatchString(email) {
				return errors.New(fieldName + " must be a valid email")
			}
		}
	case "min":
		if len(parts) == 2 {
			if field.Kind() == reflect.String {
				if len(field.String()) < parseMinLength(parts[1]) {
					return errors.New(fieldName + " must be at least " + parts[1] + " characters")
				}
			}
		}
	}

	return nil
}

func isZero(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String:
		return v.String() == ""
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Bool:
		return !v.Bool()
	}
	return false
}

func parseMinLength(s string) int {
	length := 0
	for _, c := range s {
		if c >= '0' && c <= '9' {
			length = length*10 + int(c-'0')
		}
	}
	return length
}
