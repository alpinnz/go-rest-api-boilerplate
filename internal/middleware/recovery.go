package middleware

import (
	"github.com/alpinnz/go-rest-api-boilerplate/pkg/logger"
	"github.com/alpinnz/go-rest-api-boilerplate/pkg/response"
	"github.com/gin-gonic/gin"
)

func Recovery(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				requestID := GetRequestID(c)
				log.Error().
					Str("request_id", requestID).
					Interface("panic", err).
					Msg("Panic Recovered")

				response.InternalServerError(c, "Unexpected error occurred")
				c.Abort()
			}
		}()
		c.Next()
	}
}
