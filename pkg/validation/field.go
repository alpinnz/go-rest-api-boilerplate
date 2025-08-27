package validation

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

func jsonFieldName(fe validator.FieldError, s interface{}) string {
	t := reflect.TypeOf(s)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if field, ok := t.FieldByName(fe.StructField()); ok {
		jsonTag := field.Tag.Get("json")
		if jsonTag != "" && jsonTag != "-" {
			parts := strings.Split(jsonTag, ",")
			return parts[0] // example: "last_name"
		}
	}

	return fe.Field() // fallback → "LastName"
}
