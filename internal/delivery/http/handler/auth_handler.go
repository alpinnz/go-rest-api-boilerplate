package handler

import (
	"github.com/alpinnz/go-rest-api-boilerplate/internal/delivery/http/dto"
	"github.com/alpinnz/go-rest-api-boilerplate/internal/domain"
	"github.com/alpinnz/go-rest-api-boilerplate/internal/usecase"
	"github.com/alpinnz/go-rest-api-boilerplate/pkg/response"
	"github.com/alpinnz/go-rest-api-boilerplate/pkg/validator"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authUseCase *usecase.AuthUseCase
}

func NewAuthHandler(authUseCase *usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{
		authUseCase: authUseCase,
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
func (h *AuthHandler) Register(c *gin.Context) {
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
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	user, err := h.authUseCase.Register(c.Request.Context(), input)
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
func (h *AuthHandler) Login(c *gin.Context) {
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

	result, err := h.authUseCase.Login(c.Request.Context(), input)
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
func (h *AuthHandler) Logout(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	if err := h.authUseCase.Logout(c.Request.Context(), token); err != nil {
		response.InternalServerError(c, domain.NewAppError("common.error"))
		return
	}

	response.Success(c, true)
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
// @Router       /auth/refresh-token [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
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

	result, err := h.authUseCase.RefreshToken(c.Request.Context(), req.RefreshToken)
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
