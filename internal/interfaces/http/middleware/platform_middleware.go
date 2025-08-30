package middleware

import (
	"github.com/alpinnz/go-rest-api-boilerplate/pkg/constants"
	"github.com/alpinnz/go-rest-api-boilerplate/pkg/logger"
	"github.com/gin-gonic/gin"
)

func PlatformMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		device := c.GetHeader(constants.XPlatform)
		logger.Log.Debug("%s : %s\n", constants.XPlatform, device)
		c.Next()
	}
}
