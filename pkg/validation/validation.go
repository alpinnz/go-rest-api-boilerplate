package validation

import (
	"fmt"

	"github.com/alpinnz/go-rest-api-boilerplate/pkg/translations"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Validator struct {
	tr       *translations.Store
	validate *validator.Validate
}

// NewValidator initializes validator with custom rules.
func NewValidator(tr *translations.Store) (*Validator, error) {
	v := validator.New(validator.WithRequiredStructEnabled())

	// register custom validation rules
	if err := RegisterAll(v); err != nil {
		return nil, err
	}

	return &Validator{tr: tr, validate: v}, nil
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

func (v *Validator) ValidateStructCtx(c *gin.Context, s interface{}) error {
	if err := v.validate.Struct(s); err != nil {
		var errs ValidatorErrors
		for _, fe := range err.(validator.ValidationErrors) {
			key := translations.Key(fmt.Sprintf("%s.%s", "validator", fe.Tag()))
			tr := v.tr.TGin(c, key, &map[string]any{})
			var msg string
			if fe.Param() == "" {
				// tidak ada param, cukup field name
				msg = fmt.Sprintf(tr, fmt.Sprintf("%s", jsonFieldName(fe, s)))
			} else {
				// ada param, gunakan param
				msg = fmt.Sprintf(tr, fmt.Sprintf("%s", jsonFieldName(fe, s)), fe.Param())
			}
			errs = append(errs, ValidatorError{
				Field:   jsonFieldName(fe, s),
				Tag:     fe.Tag(),
				Message: msg,
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
