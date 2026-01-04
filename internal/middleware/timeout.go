package middleware

import (
	"context"
	"time"

	"github.com/alpinnz/go-rest-api-boilerplate/pkg/response"
	"github.com/gin-gonic/gin"
)

// Timeout middleware adds a timeout to request context
func Timeout(duration time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create a context with timeout
		ctx, cancel := context.WithTimeout(c.Request.Context(), duration)
		defer cancel()

		// Replace request context with timeout context
		c.Request = c.Request.WithContext(ctx)

		// Channel to signal when request is done
		done := make(chan struct{})

		go func() {
			c.Next()
			close(done)
		}()

		select {
		case <-done:
			// Request completed successfully
			return
		case <-ctx.Done():
			// Context timeout or cancelled
			if ctx.Err() == context.DeadlineExceeded {
				response.RequestTimeout(c, "Request timeout exceeded")
				c.Abort()
			}
			return
		}
	}
}
