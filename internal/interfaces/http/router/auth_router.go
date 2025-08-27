package router

import (
	"github.com/alpinnz/go-rest-api-boilerplate/internal/interfaces/http/controllers"
	"github.com/gin-gonic/gin"
)

func NewAuthRoutes(rg *gin.RouterGroup, controller *controllers.AuthController) {
	roles := rg.Group("/auth")

	roles.POST("/login", controller.Login)
	roles.POST("/register", controller.Register)
}
