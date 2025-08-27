package router

import (
	"github.com/alpinnz/go-rest-api-boilerplate/internal/interfaces/http/controllers"
	"github.com/gin-gonic/gin"
)

func NewRoleRoutes(rg *gin.RouterGroup, controller *controllers.RoleController) {
	roles := rg.Group("/roles")

	roles.GET("", controller.GetAllRoles)
}
