package controllers

import (
	"strconv"

	"github.com/alpinnz/go-rest-api-boilerplate/internal/interfaces/http/dto"
	"github.com/alpinnz/go-rest-api-boilerplate/internal/usecase"
	"github.com/alpinnz/go-rest-api-boilerplate/pkg/errors"
	"github.com/alpinnz/go-rest-api-boilerplate/pkg/response"
	"github.com/alpinnz/go-rest-api-boilerplate/pkg/utils"
	"github.com/alpinnz/go-rest-api-boilerplate/pkg/validation"

	"github.com/gin-gonic/gin"
)

// RoleController is responsible for handling HTTP requests related to roles.
// It acts as a bridge between the HTTP layer and the role use case layer.
type RoleController struct {
	validator   *validation.Validator // Used for validating request payloads
	roleUsecase usecase.RoleUsecase   // Encapsulates business logic for roles
}

// NewRoleController initializes and returns a new RoleController instance.
func NewRoleController(validator *validation.Validator, roleUsecase usecase.RoleUsecase) *RoleController {
	return &RoleController{
		validator:   validator,
		roleUsecase: roleUsecase,
	}
}

// GetAllRoles handles GET /roles endpoint.
// It retrieves a paginated list of roles from the database.
func (h *RoleController) GetAllRoles(c *gin.Context) {
	// Parse pagination parameters from query (defaults: per_page=10, page=1)
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	offset := (page - 1) * perPage

	// Call use case to fetch roles with pagination
	result, err := h.roleUsecase.FindAllWithPagination(c.Request.Context(), perPage, offset)
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

	// Build pagination metadata for response
	pagination := utils.NewPagination(page, perPage, result.TotalData)

	// Type assert the result data into []dto.Role
	roles, ok := result.Data.([]dto.Role)
	if !ok {
		response.InternalError(c, "Failed to parse role data", nil)
		return
	}

	// Respond with paginated role data
	response.OKWithPagination[dto.Role](c, "Roles retrieved successfully", "roles", roles, pagination)
}
