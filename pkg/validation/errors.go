package validation

type ValidatorError struct {
	Field   string `json:"field,omitempty"`
	Params  string `json:"params,omitempty"`
	Tag     string `json:"tag,omitempty"`
	Message string `json:"message,omitempty"`
}

type ValidatorErrors []ValidatorError

//func (ve ValidatorErrors) Error() string {
//	msg := ""
//	for _, e := range ve {
//		msg += e.Field + ": " + e.Message + "; "
//	}
//	return msg
//}

func (ve ValidatorErrors) Error() string {
	if len(ve) == 0 {
		return ""
	}
	e := ve[0]
	return e.Message
}
