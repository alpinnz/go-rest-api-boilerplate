package handler

import (
	"strconv"

	"github.com/alpinnz/go-rest-api-boilerplate/internal/domain"
	"github.com/alpinnz/go-rest-api-boilerplate/internal/usecase"
	"github.com/alpinnz/go-rest-api-boilerplate/pkg/response"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userUseCase *usecase.UserUseCase
}

func NewUserHandler(userUseCase *usecase.UserUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
	}
}

func (h *UserHandler) Register(c *gin.Context) {
	var input usecase.RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "INVALID_INPUT", "Invalid request body", gin.H{"details": err.Error()})
		return
	}

	user, err := h.userUseCase.Register(c.Request.Context(), input)
	if err != nil {
		if err == domain.ErrConflict {
			response.Conflict(c, "email", "Email already exists")
			return
		}
		response.InternalServerError(c, "Failed to register user")
		return
	}

	response.Created(c, user)
}

func (h *UserHandler) Login(c *gin.Context) {
	var input usecase.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "INVALID_INPUT", "Invalid request body", gin.H{"details": err.Error()})
		return
	}

	result, err := h.userUseCase.Login(c.Request.Context(), input)
	if err != nil {
		if err == domain.ErrInvalidCredentials {
			response.Unauthorized(c, "Invalid email or password")
			return
		}
		response.InternalServerError(c, "Failed to process login")
		return
	}

	response.Success(c, result)
}

func (h *UserHandler) Logout(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	if err := h.userUseCase.Logout(c.Request.Context(), token); err != nil {
		response.InternalServerError(c, "Failed to logout")
		return
	}

	response.Success(c, gin.H{"message": "Logged out successfully"})
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "Authentication required")
		return
	}

	id, ok := userID.(int64)
	if !ok {
		response.InternalServerError(c, "Invalid user ID format")
		return
	}

	user, err := h.userUseCase.GetByID(c.Request.Context(), id)
	if err != nil {
		response.NotFound(c, "user")
		return
	}

	response.Success(c, user)
}

func (h *UserHandler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		response.BadRequest(c, "INVALID_ID", "Invalid user ID format", gin.H{"details": err.Error()})
		return
	}

	user, err := h.userUseCase.GetByID(c.Request.Context(), id)
	if err != nil {
		response.NotFound(c, "user")
		return
	}

	response.Success(c, user)
}
