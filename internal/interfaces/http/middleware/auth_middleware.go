package middleware

import (
	"os"
	"strings"

	"github.com/alpinnz/go-rest-api-boilerplate/internal/usecase"
	"github.com/alpinnz/go-rest-api-boilerplate/pkg/constants"
	"github.com/alpinnz/go-rest-api-boilerplate/pkg/encrypt"
	"github.com/alpinnz/go-rest-api-boilerplate/pkg/response"
	"github.com/alpinnz/go-rest-api-boilerplate/pkg/utils"
	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	UserUsecase usecase.UserUsecase
}

func NewAuthMiddleware(userUsecase usecase.UserUsecase) *AuthMiddleware {
	return &AuthMiddleware{
		UserUsecase: userUsecase,
	}
}
func (m *AuthMiddleware) IsAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(constants.XAccessToken)
		if authHeader == "" {
			response.Unauthorized(c, "Missing authorization header")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, "Bearer ", 2)
		if len(parts) != 2 || strings.TrimSpace(parts[1]) == "" {
			response.Unauthorized(c, "Invalid authorization header format")
			c.Abort()
			return
		}

		claims, err := encrypt.ExtractHash(parts[1], os.Getenv("ACCESS_TOKEN_SECRET"))
		if err != nil {
			response.Unauthorized(c, err.Error()) // e.g., "invalid token" or "invalid claims"
			c.Abort()
			return
		}

		// (Optional) verify the user still exists
		if _, err := m.UserUsecase.FindByID(c.Request.Context(), claims.UserID); err != nil {
			response.Unauthorized(c, "User not found or inactive")
			c.Abort()
			return
		}

		// Store claims in context
		ctx := utils.SetClaims(c.Request.Context(), claims)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
