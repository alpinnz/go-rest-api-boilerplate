package response

import (
	"log"

	"github.com/gin-gonic/gin"
)

// RespondSuccess sends a success response with custom status
func RespondSuccess(c *gin.Context, status Status, message string, data interface{}) {
	c.JSON(int(status), BaseResponse{
		TraceID: getTraceID(c),
		Message: message,
		Data:    data,
	})
}

// RespondError sends an errors response with status and optional code
func RespondError(c *gin.Context, status Status, code Code, message string, errors interface{}) {
	traceID := getTraceID(c)

	// log to stdout or file
	log.Printf("[ERROR] trace_id=%s | code=%s | message=%s | errors=%v", traceID, code, message, errors)

	c.JSON(int(status), BaseResponse{
		TraceID: traceID,
		Code:    string(code),
		Message: message,
		Errors:  errors,
	})
}
