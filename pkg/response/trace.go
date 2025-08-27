package response

import "github.com/gin-gonic/gin"

// getTraceID extracts trace ID from context (if exists)
func getTraceID(c *gin.Context) string {
	return c.GetString("trace_id")
}
