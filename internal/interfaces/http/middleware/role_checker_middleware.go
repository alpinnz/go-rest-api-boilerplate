package middleware

import (
	"github.com/alpinnz/go-rest-api-boilerplate/pkg/constants"
	"github.com/alpinnz/go-rest-api-boilerplate/pkg/helper"
	"github.com/alpinnz/go-rest-api-boilerplate/pkg/response"
	"github.com/alpinnz/go-rest-api-boilerplate/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func AuthorizeRole(roleID uuid.UUID, roleName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := utils.GetClaims(c.Request.Context())
		if err != nil || !helper.ContainsRole(claims.Roles, roleID) {
			response.Unauthorized(c, "Unauthorized: "+roleName+" role required")
			c.Abort()
			return
		}
		c.Next()
	}
}

func IsUser() gin.HandlerFunc {
	return AuthorizeRole(constants.RoleIDUser, "User")
}

func IsAdmin() gin.HandlerFunc {
	return AuthorizeRole(constants.RoleIDAdmin, "Admin")
}
