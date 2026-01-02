package handler

import (
	"strconv"

	"github.com/alpinnz/go-rest-api-boilerplate/internal/delivery/http/dto"
	"github.com/alpinnz/go-rest-api-boilerplate/internal/domain"
	"github.com/alpinnz/go-rest-api-boilerplate/internal/usecase"
	"github.com/alpinnz/go-rest-api-boilerplate/pkg/response"
	"github.com/alpinnz/go-rest-api-boilerplate/pkg/validator"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RoleHandler struct {
	roleUseCase *usecase.RoleUseCase
}

func NewRoleHandler(roleUseCase *usecase.RoleUseCase) *RoleHandler {
	return &RoleHandler{
		roleUseCase: roleUseCase,
	}
}

// Create godoc
// @Summary      Create role
// @Description  Create a new role
// @Tags         Roles
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        Accept-Language  header  string  false  "Language preference (en, id)"  default(en)
// @Param        request body dto.CreateRoleRequest true "Create role request"
// @Success      201  {object}  response.Response{data=dto.RoleDTO}
// @Failure      400  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Failure      409  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Router       /roles [post]
func (h *RoleHandler) Create(c *gin.Context) {
	var req dto.CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, domain.NewAppError("request.invalid_id"))
		return
	}

	if errs := validator.ValidateStruct(req); len(errs) > 0 {
		response.ValidationError(c, errs)
		return
	}

	input := usecase.CreateRoleInput{
		Name:        req.Name,
		Description: req.Description,
	}

	role, err := h.roleUseCase.Create(c.Request.Context(), input)
	if err != nil {
		if err == domain.ErrConflict {
			response.Conflict(c, domain.NewAppError("role.name_exists"))
			return
		}
		response.InternalServerError(c, domain.NewAppError("common.error"))
		return
	}

	response.Created(c, dto.ToRoleDTO(role), "role.created")
}

// GetByID godoc
// @Summary      Get role by ID
// @Description  Get role information by role ID
// @Tags         Roles
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        Accept-Language  header  string  false  "Language preference (en, id)"  default(en)
// @Param        id   path      int  true  "Role ID"
// @Success      200  {object}  response.Response{data=dto.RoleDTO}
// @Failure      400  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Failure      404  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Router       /roles/{id} [get]
func (h *RoleHandler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		response.BadRequest(c, domain.NewAppError("request.invalid_id"))
		return
	}

	role, err := h.roleUseCase.GetByID(c.Request.Context(), id)
	if err != nil {
		if err == domain.ErrNotFound {
			response.NotFound(c, domain.NewAppError("role.not_found"))
			return
		}
		response.InternalServerError(c, domain.NewAppError("common.error"))
		return
	}

	response.Success(c, dto.ToRoleDTO(role))
}

// List godoc
// @Summary      List roles
// @Description  Get paginated list of roles
// @Tags         Roles
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        Accept-Language  header  string  false  "Language preference (en, id)"  default(en)
// @Param        limit   query     int     false  "Limit"  default(10)
// @Param        offset  query     int     false  "Offset" default(0)
// @Success      200  {object}  response.Response{data=dto.RoleListResponse}
// @Failure      401  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Router       /roles [get]
func (h *RoleHandler) List(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	input := usecase.ListRolesInput{
		Limit:  limit,
		Offset: offset,
	}

	result, err := h.roleUseCase.List(c.Request.Context(), input)
	if err != nil {
		response.InternalServerError(c, domain.NewAppError("common.error"))
		return
	}

	responseDTO := dto.ToRoleListResponse(result.Roles, result.Total, result.Limit, result.Offset, result.HasMore)
	response.Success(c, responseDTO)
}

// Update godoc
// @Summary      Update role
// @Description  Update role information
// @Tags         Roles
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        Accept-Language  header  string  false  "Language preference (en, id)"  default(en)
// @Param        id   path      int  true  "Role ID"
// @Param        request body dto.UpdateRoleRequest true "Update role request"
// @Success      200  {object}  response.Response{data=dto.RoleDTO}
// @Failure      400  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Failure      404  {object}  response.Response
// @Failure      409  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Router       /roles/{id} [put]
func (h *RoleHandler) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		response.BadRequest(c, domain.NewAppError("request.invalid_id"))
		return
	}

	var req dto.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, domain.NewAppError("request.invalid_id"))
		return
	}

	if errs := validator.ValidateStruct(req); len(errs) > 0 {
		response.ValidationError(c, errs)
		return
	}

	input := usecase.UpdateRoleInput{
		Name:        req.Name,
		Description: req.Description,
	}

	role, err := h.roleUseCase.Update(c.Request.Context(), id, input)
	if err != nil {
		if err == domain.ErrNotFound {
			response.NotFound(c, domain.NewAppError("role.not_found"))
			return
		}
		if err == domain.ErrConflict {
			response.Conflict(c, domain.NewAppError("role.name_exists"))
			return
		}
		response.InternalServerError(c, domain.NewAppError("common.error"))
		return
	}

	response.Success(c, dto.ToRoleDTO(role))
}

// Delete godoc
// @Summary      Delete role
// @Description  Delete role by ID
// @Tags         Roles
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        Accept-Language  header  string  false  "Language preference (en, id)"  default(en)
// @Param        id   path      int  true  "Role ID"
// @Success      200  {object}  response.Response
// @Failure      400  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Failure      404  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Router       /roles/{id} [delete]
func (h *RoleHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		response.BadRequest(c, domain.NewAppError("request.invalid_id"))
		return
	}

	if err := h.roleUseCase.Delete(c.Request.Context(), id); err != nil {
		if err == domain.ErrNotFound {
			response.NotFound(c, domain.NewAppError("role.not_found"))
			return
		}
		response.InternalServerError(c, domain.NewAppError("common.error"))
		return
	}

	response.Success(c, true)
}
