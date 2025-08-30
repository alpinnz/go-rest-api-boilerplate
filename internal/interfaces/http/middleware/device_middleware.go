package middleware

import (
	"github.com/alpinnz/go-rest-api-boilerplate/pkg/constants"
	"github.com/alpinnz/go-rest-api-boilerplate/pkg/logger"
	"github.com/gin-gonic/gin"
)

func DeviceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		device := c.GetHeader(constants.XDevice)
		logger.Log.Debug("%s : %s\n", constants.XDevice, device)
		c.Next()
	}
}
