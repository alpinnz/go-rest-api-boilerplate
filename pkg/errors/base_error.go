package errors

import (
	"fmt"

	"github.com/alpinnz/go-rest-api-boilerplate/pkg/response"
)

type BaseError struct {
	Status  response.Status `json:"status"`
	Code    response.Code   `json:"code"`
	Message string          `json:"message"`
	Errors  any             `json:"errors,omitempty"` // optional, bisa kosong
}

// Error() implements the errors interface
func (e *BaseError) Error() string {
	return e.Message
}

func NewBadRequest(msg string, errs ...any) *BaseError {
	var errorsDetail any
	if len(errs) > 0 {
		errorsDetail = errs[0]
	}
	return &BaseError{
		Status:  response.StatusBadRequest,
		Code:    response.CodeBadRequest,
		Message: msg,
		Errors:  errorsDetail,
	}
}

func NewInternalError(msg string, errs any) *BaseError {
	return &BaseError{
		Status:  response.StatusInternalServerError,
		Code:    response.CodeInternalError,
		Message: msg,
		Errors:  errs,
	}
}

func NewNotFound(msg string, err string) *BaseError {
	return &BaseError{
		Status:  response.StatusNotFound,
		Code:    response.CodeNotFound,
		Message: msg,
		Errors:  err,
	}
}

func NewDatabaseError(details string, errs any) *BaseError {
	return &BaseError{
		Status:  response.StatusInternalServerError,
		Code:    response.CodeDatabaseError,
		Message: fmt.Sprintf("Database errors: %s\n", details),
		Errors:  errs,
	}
}
