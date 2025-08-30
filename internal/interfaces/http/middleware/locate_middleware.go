package middleware

import (
	"github.com/alpinnz/go-rest-api-boilerplate/pkg/constants"
	"github.com/alpinnz/go-rest-api-boilerplate/pkg/logger"
	"github.com/gin-gonic/gin"
)

func LocateMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		locate := c.GetHeader(constants.XLocate)
		logger.Log.Debug("%s : %s\n", constants.XLocate, locate)
		c.Next()
	}
}
