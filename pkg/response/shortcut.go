package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Ok Shortcut for 200 OK
func Ok(c *gin.Context, data interface{}) {
	RespondSuccess(c, http.StatusOK, "Success", data)
}

func OKWithPagination[T any](c *gin.Context, message string, key string, items []T, pagination Pagination) {
	resp := ListWithPagination[T]{
		Key:        key,
		Items:      items,
		Pagination: pagination,
	}

	c.JSON(int(StatusOK), BaseResponse{
		TraceID: getTraceID(c),
		Message: message,
		Data:    resp.ToMap(),
	})
}

// Created Shortcut for 201 Created
func Created(c *gin.Context, data interface{}) {
	RespondSuccess(c, http.StatusCreated, "Success", data)
}

// BadRequest Shortcut for 400 Bad Request
func BadRequest(c *gin.Context, message string, errors interface{}) {
	RespondError(c, http.StatusBadRequest, CodeBadRequest, message, errors)
}

// Unauthorized Shortcut for 401 Unauthorized
func Unauthorized(c *gin.Context, message string) {
	RespondError(c, http.StatusUnauthorized, CodeUnauthorized, message, nil)
}

// Forbidden Shortcut for 403 Forbidden
func Forbidden(c *gin.Context, message string) {
	RespondError(c, http.StatusForbidden, CodeForbidden, message, nil)
}

// NotFound Shortcut for 404 Not Found
func NotFound(c *gin.Context, message string) {
	RespondError(c, http.StatusNotFound, CodeNotFound, message, nil)
}

// Conflict Shortcut for 409 Conflict
func Conflict(c *gin.Context, message string) {
	RespondError(c, http.StatusConflict, CodeConflict, message, nil)
}

// Unprocessable Shortcut for 422 Validation Error
func Unprocessable(c *gin.Context, message string, errors interface{}) {
	RespondError(c, http.StatusUnprocessableEntity, CodeValidationError, message, errors)
}

// InternalError Shortcut for 500 Internal Server Error
func InternalError(c *gin.Context, message string, errors interface{}) {
	RespondError(c, http.StatusInternalServerError, CodeInternalError, message, errors)
}
