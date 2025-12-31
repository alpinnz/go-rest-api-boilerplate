package handler

import (
	"github.com/alpinnz/go-rest-api-boilerplate/pkg/response"
	"github.com/gin-gonic/gin"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// Check godoc
// @Summary      Health check
// @Description  Check if the API is running
// @Tags         Health
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response
// @Router       /health [get]
func (h *HealthHandler) Check(c *gin.Context) {
	response.Success(c, true)
}
