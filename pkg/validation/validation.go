package validation

import (
	"github.com/go-playground/validator/v10"
)

type Validator struct {
	validate *validator.Validate
}

// NewValidator initializes validator with custom rules.
func NewValidator() (*Validator, error) {
	v := validator.New(validator.WithRequiredStructEnabled())

	// register custom validation rules
	if err := RegisterAll(v); err != nil {
		return nil, err
	}

	return &Validator{validate: v}, nil
}

// ValidateStruct validates a struct.
func (v *Validator) ValidateStruct(s interface{}) error {
	if err := v.validate.Struct(s); err != nil {
		var errs ValidatorErrors
		for _, fe := range err.(validator.ValidationErrors) {
			errs = append(errs, ValidatorError{
				Field:   jsonFieldName(fe, s),
				Tag:     fe.Tag(),
				Message: messageForTag(fe, s),
				Params:  fe.Param(),
			})
		}
		return errs
	}
	return nil
}

// ValidateField validates a single field.
func (v *Validator) ValidateField(field interface{}, tag string) error {
	return v.validate.Var(field, tag)
}
