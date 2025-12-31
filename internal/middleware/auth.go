package middleware

import (
	"strings"

	"github.com/alpinnz/go-rest-api-boilerplate/internal/usecase"
	"github.com/alpinnz/go-rest-api-boilerplate/pkg/response"
	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	userUseCase *usecase.UserUseCase
}

func NewAuthMiddleware(userUseCase *usecase.UserUseCase) *AuthMiddleware {
	return &AuthMiddleware{
		userUseCase: userUseCase,
	}
}

func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "Authorization header required")
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Unauthorized(c, "Invalid authorization format")
			c.Abort()
			return
		}

		token := parts[1]
		userID, err := m.userUseCase.ValidateSession(c.Request.Context(), token)
		if err != nil {
			response.Unauthorized(c, "Invalid or expired token")
			c.Abort()
			return
		}

		c.Set("userID", userID)
		c.Next()
	}
}
