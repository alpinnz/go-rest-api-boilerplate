package handler

import (
	"strconv"

	"github.com/alpinnz/go-rest-api-boilerplate/internal/delivery/http/dto"
	"github.com/alpinnz/go-rest-api-boilerplate/internal/domain"
	"github.com/alpinnz/go-rest-api-boilerplate/internal/usecase"
	"github.com/alpinnz/go-rest-api-boilerplate/pkg/response"
	"github.com/alpinnz/go-rest-api-boilerplate/pkg/validator"
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

// Register godoc
// @Summary      Register new user
// @Description  Create a new user account
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        Accept-Language  header  string  false  "Language preference (en, id)"  default(en)
// @Param        request body dto.RegisterRequest true "Register request"
// @Success      201  {object}  response.Response{data=dto.UserDTO}
// @Failure      400  {object}  response.Response
// @Failure      409  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Router       /auth/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, domain.NewAppError("request.invalid_json"))
		return
	}

	// Validate input
	if errs := validator.ValidateStruct(req); len(errs) > 0 {
		response.ValidationError(c, errs)
		return
	}

	// Map DTO to usecase input
	input := usecase.RegisterInput{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
	}

	user, err := h.userUseCase.Register(c.Request.Context(), input)
	if err != nil {
		if err == domain.ErrConflict {
			response.Conflict(c, domain.NewAppError("user.email_exists"))
			return
		}
		response.InternalServerError(c, domain.NewAppError("common.error"))
		return
	}

	// Map domain entity to DTO for response
	response.Created(c, dto.ToUserDTO(user), "user.created")
}

// Login godoc
// @Summary      User login
// @Description  Authenticate user and get JWT token
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        Accept-Language  header  string  false  "Language preference (en, id)"  default(en)
// @Param        request body dto.LoginRequest true "Login credentials"
// @Success      200  {object}  response.Response{data=dto.LoginResponse}
// @Failure      400  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Router       /auth/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, domain.NewAppError("request.invalid_json"))
		return
	}

	// Validate input
	if errs := validator.ValidateStruct(req); len(errs) > 0 {
		response.ValidationError(c, errs)
		return
	}

	// Map DTO to usecase input
	input := usecase.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	}

	result, err := h.userUseCase.Login(c.Request.Context(), input)
	if err != nil {
		if err == domain.ErrInvalidCredentials {
			response.Unauthorized(c, domain.NewAppError("auth.invalid_credentials"))
			return
		}
		response.InternalServerError(c, domain.NewAppError("common.error"))
		return
	}

	// Map usecase result to DTO for response
	responseDTO := &dto.LoginResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		User:         dto.ToUserDTO(result.User),
	}

	response.Success(c, responseDTO)
}

// Logout godoc
// @Summary      User logout
// @Description  Logout user and invalidate JWT token
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Router       /auth/logout [post]
func (h *UserHandler) Logout(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	if err := h.userUseCase.Logout(c.Request.Context(), token); err != nil {
		response.InternalServerError(c, domain.NewAppError("common.error"))
		return
	}

	response.Success(c, true)
}

// GetProfile godoc
// @Summary      GetAccessToken current user profile
// @Description  GetAccessToken authenticated user's profile information
// @Tags         Users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        Accept-Language  header  string  false  "Language preference (en, id)"  default(en)
// @Success      200  {object}  response.Response{data=dto.UserDTO}
// @Failure      401  {object}  response.Response
// @Failure      404  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Router       /users/me [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, domain.NewAppError("auth.unauthorized"))
		return
	}

	id, ok := userID.(int64)
	if !ok {
		response.InternalServerError(c, domain.NewAppError("common.error"))
		return
	}

	user, err := h.userUseCase.GetByID(c.Request.Context(), id)
	if err != nil {
		if err == domain.ErrNotFound {
			response.NotFound(c, domain.NewAppError("user.not_found"))
			return
		}
		response.InternalServerError(c, domain.NewAppError("common.error"))
		return
	}

	response.Success(c, dto.ToUserDTO(user))
}

// GetByID godoc
// @Summary      GetAccessToken user by ID
// @Description  GetAccessToken user information by user ID
// @Tags         Users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        Accept-Language  header  string  false  "Language preference (en, id)"  default(en)
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  response.Response{data=dto.UserDTO}
// @Failure      400  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Failure      404  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Router       /users/{id} [get]
func (h *UserHandler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		response.BadRequest(c, domain.NewAppError("request.invalid_json"))
		return
	}

	user, err := h.userUseCase.GetByID(c.Request.Context(), id)
	if err != nil {
		if err == domain.ErrNotFound {
			response.NotFound(c, domain.NewAppError("user.not_found"))
			return
		}
		response.InternalServerError(c, domain.NewAppError("common.error"))
		return
	}

	response.Success(c, dto.ToUserDTO(user))
}

// RefreshToken godoc
// @Summary      Refresh access token
// @Description  Generate new access token using refresh token
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        Accept-Language  header  string  false  "Language preference (en, id)"  default(en)
// @Param        request body dto.RefreshTokenRequest true "Refresh token"
// @Success      200  {object}  response.Response{data=dto.LoginResponse}
// @Failure      400  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Router       /auth/refresh [post]
func (h *UserHandler) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, domain.NewAppError("request.invalid_json"))
		return
	}

	// Validate input
	if errs := validator.ValidateStruct(req); len(errs) > 0 {
		response.ValidationError(c, errs)
		return
	}

	result, err := h.userUseCase.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		if err == domain.ErrInvalidCredentials || err == domain.ErrSessionExpired {
			response.Unauthorized(c, domain.NewAppError("auth.invalid_refresh_token"))
			return
		}
		response.InternalServerError(c, domain.NewAppError("common.error"))
		return
	}

	// Map usecase result to DTO for response
	responseDTO := &dto.LoginResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		User:         dto.ToUserDTO(result.User),
	}

	response.Success(c, responseDTO)
}
