package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const RequestIDKey = "request_id"

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate request ID if not exists
		if _, exists := c.Get(RequestIDKey); !exists {
			c.Set(RequestIDKey, uuid.New().String())
		}

		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method
		requestID := GetRequestID(c)
		clientIP := c.ClientIP()
		userAgent := c.Request.UserAgent()

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()

		// Structured log format with request_id, client IP, and user agent
		log.Printf("[request_id=%s] [ip=%s] [method=%s] [path=%s] [status=%d] [latency=%v] [user_agent=%s]",
			requestID, clientIP, method, path, statusCode, latency, userAgent)

		// Log errors separately for better visibility
		if len(c.Errors) > 0 {
			log.Printf("[request_id=%s] [errors] %v", requestID, c.Errors.String())
		}
	}
}

// GetRequestID retrieves request ID from Gin context.
func GetRequestID(c *gin.Context) string {
	if requestID, exists := c.Get(RequestIDKey); exists {
		if id, ok := requestID.(string); ok {
			return id
		}
	}
	return ""
}
