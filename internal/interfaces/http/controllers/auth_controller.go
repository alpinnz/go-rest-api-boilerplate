package controllers

import (
	"github.com/alpinnz/go-rest-api-boilerplate/internal/interfaces/http/dto"
	"github.com/alpinnz/go-rest-api-boilerplate/internal/usecase"
	"github.com/alpinnz/go-rest-api-boilerplate/pkg/errors"
	"github.com/alpinnz/go-rest-api-boilerplate/pkg/helper"
	"github.com/alpinnz/go-rest-api-boilerplate/pkg/response"
	"github.com/alpinnz/go-rest-api-boilerplate/pkg/translations"
	"github.com/alpinnz/go-rest-api-boilerplate/pkg/validation"

	"github.com/gin-gonic/gin"
)

// AuthController is responsible for handling HTTP requests related to authentication.
// It acts as a bridge between the HTTP layer and the authentication use case layer.
type AuthController struct {
	validator   *validation.Validator // Used for validating request payloads
	tr          *translations.Store
	authUsecase usecase.AuthUsecase // Encapsulates business logic for authentication
}

// NewAuthController initializes and returns a new AuthController instance.
func NewAuthController(validator *validation.Validator, tr *translations.Store, authUsecase usecase.AuthUsecase) *AuthController {
	return &AuthController{
		validator:   validator,
		tr:          tr,
		authUsecase: authUsecase,
	}
}

// Login handles POST /auth/login endpoint.
// It authenticates the user using email & password, then returns JWT tokens.
func (h *AuthController) Login(c *gin.Context) {
	var input dto.Login

	// Step 1: Bind & validate request body
	if err := helper.ShouldBindJSON(c, &input); err != nil {
		msg := h.tr.TGin(c, translations.APP_INVALID_REQUEST_PAYLOAD, nil)
		response.BadRequest(c, msg, err)
		return
	}
	if err := h.validator.ValidateStructCtx(c, &input); err != nil {
		response.BadRequest(c, err.Error(), err)
		return
	}

	// Step 2: Call use case to authenticate user
	authSession, err := h.authUsecase.Login(c.Request.Context(), input)
	if err != nil {
		// Handle custom application errors
		if appErr, ok := err.(*errors.BaseError); ok {
			response.RespondError(c, appErr.Status, appErr.Code, appErr.Message, appErr.Errors)
			return
		}
		// Handle unexpected/internal errors
		response.InternalError(c, err.Error(), err)
		return
	}

	// Step 3: Respond with login success
	response.Ok(c, authSession)
}

// Register handles POST /auth/register endpoint.
// It authenticates the user using email & password, then returns JWT tokens.
func (h *AuthController) Register(c *gin.Context) {
	var input dto.Register

	// Step 1: Bind & validate request body
	if err := helper.ShouldBindJSON(c, &input); err != nil {
		msg := h.tr.TGin(c, translations.APP_INVALID_REQUEST_PAYLOAD, nil)
		response.BadRequest(c, msg, err)
		return
	}
	if err := h.validator.ValidateStructCtx(c, &input); err != nil {
		response.BadRequest(c, err.Error(), err)
		return
	}

	// Step 2: Call use case to authenticate user
	authSession, err := h.authUsecase.Register(c.Request.Context(), input)
	if err != nil {
		// Handle custom application errors
		if appErr, ok := err.(*errors.BaseError); ok {
			response.RespondError(c, appErr.Status, appErr.Code, appErr.Message, appErr.Errors)
			return
		}
		// Handle unexpected/internal errors
		response.InternalError(c, err.Error(), err)
		return
	}

	// Step 3: Respond with register success
	response.Created(c, authSession)
}
