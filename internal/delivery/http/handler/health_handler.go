package handler

import (
	"github.com/alpinnz/go-rest-api-boilerplate/pkg/response"
	"github.com/gin-gonic/gin"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Check(c *gin.Context) {
	response.Success(c, gin.H{
		"status":  "ok",
		"message": "Service is running",
	})
}
