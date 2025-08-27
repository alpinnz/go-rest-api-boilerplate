package validation

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// tagMessageMap maps validator tags to custom human-readable errors messages.
// Each tag corresponds to a rule defined by go-playground/validator.
// Use msg() if only the field name is needed, or msgWithParam() if Param() is required.
var tagMessageMap = map[string]func(fe validator.FieldError, s interface{}) string{
	// Common rules
	"omitempty": func(fe validator.FieldError, s interface{}) string { return "" }, // skip empty field
	"required":  msg("%s is required"),
	"eqfield":   msgWithParam("%s must equal %s"),
	"gt":        msgWithParam("%s must be greater than %s"),
	"gte":       msgWithParam("%s must be greater than or equal to %s"),
	"lt":        msgWithParam("%s must be less than %s"),
	"lte":       msgWithParam("%s must be less than or equal to %s"),

	// Format validators
	"email":       msg("%s must be a valid email"),
	"url":         msg("%s must be a valid URL"),
	"alphanum":    msg("%s must contain only alphanumeric characters"),
	"alpha":       msg("%s must contain only alphabetic characters"),
	"contains":    msgWithParam("%s must contain '%s'"),
	"excludes":    msgWithParam("%s must not contain '%s'"),
	"ip":          msg("%s must be a valid IP address"),
	"ipv4":        msg("%s must be a valid IPv4 address"),
	"ipv6":        msg("%s must be a valid IPv6 address"),
	"hostname":    msg("%s must be a valid hostname"),
	"fqdn":        msg("%s must be a fully qualified domain name"),
	"cidr":        msg("%s must be a valid CIDR notation"),
	"credit_card": msg("%s must be a valid credit card number"),
	"uuid":        msg("%s must be a valid UUID"),

	// Custom regex-based rules
	"regex_name":  msg("%s must be a valid name format"),
	"regex_uuid7": msg("%s must be a valid UUIDv7"),

	// Kind/type validators
	"kind_string": msg("%s must be a string"),
	"kind_bool":   msg("%s must be a boolean"),
	"kind_int":    msg("%s must be an integer"),
	"kind_array":  msg("%s must be an array"),

	// Password strength rule
	"strong_password": msg("%s must contain at least 8 characters, uppercase, lowercase, number, and special char"),
}

// messageForTag resolves the errors message based on the validator tag.
// If the tag is not mapped, it falls back to a generic "is not valid" message.
func messageForTag(fe validator.FieldError, s interface{}) string {
	if fn, ok := tagMessageMap[fe.Tag()]; ok {
		if msg := fn(fe, s); msg != "" {
			return msg
		}
	}
	return fmt.Sprintf("%s is not valid (%s)", jsonFieldName(fe, s), fe.Tag())
}
