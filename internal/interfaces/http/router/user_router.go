package router

import (
	"github.com/alpinnz/go-rest-api-boilerplate/internal/interfaces/http/controllers"
	"github.com/gin-gonic/gin"
)

func NewUserRoutes(rg *gin.RouterGroup, controller *controllers.UserController) {
	users := rg.Group("/users")

	users.GET("", controller.GetAllUsers)
	users.GET("/:id", controller.GetUserByID)
	users.POST("", controller.CreateUser)
	users.PATCH("/:id", controller.UpdateUser)
	users.DELETE("/:id", controller.DeleteUser)
}
