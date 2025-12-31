package middleware

import (
	"log"

	"github.com/alpinnz/go-rest-api-boilerplate/pkg/response"
	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				requestID := GetRequestID(c)
				log.Printf("[request_id=%s] [PANIC] %v", requestID, err)
				response.InternalServerError(c, "Unexpected error occurred")
				c.Abort()
			}
		}()
		c.Next()
	}
}
