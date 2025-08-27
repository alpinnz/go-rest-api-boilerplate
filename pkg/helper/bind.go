package helper

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func ShouldBindJSON(c *gin.Context, obj interface{}) error {
	// Bind JSON
	if err := c.ShouldBindJSON(obj); err != nil {
		if err.Error() == "EOF" {
			return errors.New("request body cannot be empty")
		}
		return err
	}
	return nil
}
