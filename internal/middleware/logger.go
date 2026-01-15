package middleware

import (
	"time"

	"github.com/alpinnz/go-rest-api-boilerplate/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const RequestIDKey = "request_id"

func Logger(l *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate request ID if not exists
		requestID := uuid.New().String()
		c.Set(RequestIDKey, requestID)

		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method
		clientIP := c.ClientIP()
		userAgent := c.Request.UserAgent()

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()

		// Structured logging
		l.Info().
			Str("request_id", requestID).
			Str("client_ip", clientIP).
			Str("method", method).
			Str("path", path).
			Int("status", statusCode).
			Dur("latency", latency).
			Str("user_agent", userAgent).
			Msg("HTTP Request")

		// Log errors separately
		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				l.Error().
					Str("request_id", requestID).
					Err(err).
					Msg("Request Error")
			}
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
