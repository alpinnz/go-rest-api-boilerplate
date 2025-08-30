package controllers

import (
	"net/http"
	"strconv"

	"github.com/alpinnz/go-rest-api-boilerplate/internal/interfaces/http/dto"
	"github.com/alpinnz/go-rest-api-boilerplate/internal/usecase"
	"github.com/alpinnz/go-rest-api-boilerplate/pkg/errors"
	"github.com/alpinnz/go-rest-api-boilerplate/pkg/helper"
	"github.com/alpinnz/go-rest-api-boilerplate/pkg/response"
	"github.com/alpinnz/go-rest-api-boilerplate/pkg/translations"
	"github.com/alpinnz/go-rest-api-boilerplate/pkg/utils"
	"github.com/alpinnz/go-rest-api-boilerplate/pkg/validation"

	"github.com/gin-gonic/gin"
)

// UserController handles HTTP requests related to users
type UserController struct {
	validator   *validation.Validator
	tr          *translations.Store
	userUsecase usecase.UserUsecase
}

// NewUserController initializes and returns a new UserController instance
func NewUserController(validator *validation.Validator, tr *translations.Store, userUsecase usecase.UserUsecase) *UserController {
	return &UserController{
		validator:   validator,
		tr:          tr,
		userUsecase: userUsecase,
	}
}

// GetAllUsers handles GET /users
// Retrieves a paginated list of users from the system
func (h *UserController) GetAllUsers(c *gin.Context) {
	// Get pagination parameters from query string
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	offset := (page - 1) * perPage

	// Call usecase to fetch users with pagination
	result, err := h.userUsecase.FindAllWithPagination(c.Request.Context(), perPage, offset)
	if err != nil {
		// Handle custom application error
		if appErr, ok := err.(*errors.BaseError); ok {
			response.RespondError(c, appErr.Status, appErr.Code, appErr.Message, appErr.Errors)
			return
		}
		// Handle unexpected internal error
		response.InternalError(c, err.Error(), err)
		return
	}

	// Build pagination metadata
	pagination := utils.NewPagination(page, perPage, result.TotalData)

	// Convert result data to DTO
	users, ok := result.Data.([]dto.User)
	if !ok {
		msg := h.tr.TGin(c, translations.APP_FAILED_TO_PARSE_DATA, nil)
		response.InternalError(c, msg, err)
		return
	}
	// Send paginated success response
	msg := h.tr.TGin(c, translations.APP_RETRIEVED_SUCCESS, &map[string]any{"Name": "Users"})
	response.OKWithPagination[dto.User](c, msg, "users", users, pagination)
}

// GetUserByID handles GET /users/:id
// Retrieves a specific user by its ID
func (h *UserController) GetUserByID(c *gin.Context) {
	// Parse UUID from request path
	userID, err := utils.ParseUUID(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid user ID", nil)
		return
	}

	// Call usecase to get user by ID
	user, err := h.userUsecase.FindByID(c.Request.Context(), userID)
	if err != nil {
		if appErr, ok := err.(*errors.BaseError); ok {
			response.RespondError(c, appErr.Status, appErr.Code, appErr.Message, appErr.Errors)
			return
		}
		response.InternalError(c, err.Error(), err)
		return
	}

	// Send success response
	msg := h.tr.TGin(c, translations.APP_RETRIEVED_SUCCESS, &map[string]any{"Name": "User"})
	response.RespondSuccess(c, http.StatusOK, msg, user)
}

// CreateUser handles POST /users
// Creates a new user in the system
func (h *UserController) CreateUser(c *gin.Context) {
	var userInput dto.Register

	// Bind JSON request body into DTO
	if err := helper.ShouldBindJSON(c, &userInput); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	// Validate input data
	if err := h.validator.ValidateStructCtx(c, &userInput); err != nil {
		response.BadRequest(c, err.Error(), err)
		return
	}

	// Call usecase to create user
	createdUser, err := h.userUsecase.Create(c.Request.Context(), &userInput)
	if err != nil {
		if appErr, ok := err.(*errors.BaseError); ok {
			response.RespondError(c, appErr.Status, appErr.Code, appErr.Message, appErr.Errors)
			return
		}
		response.InternalError(c, err.Error(), err)
		return
	}

	// Send success response with created user
	msg := h.tr.TGin(c, translations.APP_CREATE_SUCCESS, &map[string]any{"Name": "User"})
	response.RespondSuccess(c, http.StatusCreated, msg, createdUser)
}

// UpdateUser handles PUT /users/:id
// Updates an existing user's details
func (h *UserController) UpdateUser(c *gin.Context) {
	// Parse UUID from request path
	userID, err := utils.ParseUUID(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid user ID", nil)
		return
	}

	// Bind JSON payload
	var payload dto.UserUpdate
	if err := helper.ShouldBindJSON(c, &payload); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	// Validate input payload
	if err := h.validator.ValidateStructCtx(c, &payload); err != nil {
		response.BadRequest(c, err.Error(), err)
		return
	}

	// Call usecase to update user
	if err := h.userUsecase.Update(c.Request.Context(), userID, &payload); err != nil {
		if appErr, ok := err.(*errors.BaseError); ok {
			response.RespondError(c, appErr.Status, appErr.Code, appErr.Message, appErr.Errors)
			return
		}
		response.InternalError(c, err.Error(), err)
		return
	}

	// Send success response
	msg := h.tr.TGin(c, translations.APP_UPDATE_SUCCESS, &map[string]any{"Name": "User"})
	response.RespondSuccess(c, http.StatusOK, msg, nil)
}

// DeleteUser handles DELETE /users/:id
// Deletes a user by ID
func (h *UserController) DeleteUser(c *gin.Context) {
	// Parse UUID from request path
	userID, err := utils.ParseUUID(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid user ID", nil)
		return
	}

	// Call usecase to delete user
	if err := h.userUsecase.Delete(c.Request.Context(), userID); err != nil {
		if appErr, ok := err.(*errors.BaseError); ok {
			response.RespondError(c, appErr.Status, appErr.Code, appErr.Message, appErr.Errors)
			return
		}
		response.InternalError(c, err.Error(), err)
		return
	}

	// Send success response
	msg := h.tr.TGin(c, translations.APP_DELETE_SUCCESS, &map[string]any{"Name": "User"})
	response.RespondSuccess(c, http.StatusOK, msg, nil)
}
