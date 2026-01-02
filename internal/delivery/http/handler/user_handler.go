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

type UserHandler struct {
	userUseCase *usecase.UserUseCase
}

func NewUserHandler(userUseCase *usecase.UserUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
	}
}

// GetProfile godoc
// @Summary      Get current user profile
// @Description  Get authenticated user's profile information
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

	id, ok := userID.(uuid.UUID)
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
// @Summary      Get user by ID
// @Description  Get user information by user ID
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
	id, err := uuid.Parse(idParam)
	if err != nil {
		response.BadRequest(c, domain.NewAppError("request.invalid_id"))
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

// List godoc
// @Summary      List users
// @Description  Get paginated list of users
// @Tags         Users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        Accept-Language  header  string  false  "Language preference (en, id)"  default(en)
// @Param        limit   query     int     false  "Limit"  default(10)
// @Param        offset  query     int     false  "Offset" default(0)
// @Success      200  {object}  response.Response{data=dto.UserListResponse}
// @Failure      401  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Router       /users [get]
func (h *UserHandler) List(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	input := usecase.ListUsersInput{
		Limit:  limit,
		Offset: offset,
	}

	result, err := h.userUseCase.List(c.Request.Context(), input)
	if err != nil {
		response.InternalServerError(c, domain.NewAppError("common.error"))
		return
	}

	responseDTO := dto.ToUserListResponse(result.Users, result.Total, result.Limit, result.Offset, result.HasMore)
	response.Success(c, responseDTO)
}

// Update godoc
// @Summary      Update user
// @Description  Update user information
// @Tags         Users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        Accept-Language  header  string  false  "Language preference (en, id)"  default(en)
// @Param        id   path      int  true  "User ID"
// @Param        request body dto.UpdateUserRequest true "Update user request"
// @Success      200  {object}  response.Response{data=dto.UserDTO}
// @Failure      400  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Failure      404  {object}  response.Response
// @Failure      409  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Router       /users/{id} [put]
func (h *UserHandler) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		response.BadRequest(c, domain.NewAppError("request.invalid_id"))
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, domain.NewAppError("request.invalid_id"))
		return
	}

	if errs := validator.ValidateStruct(req); len(errs) > 0 {
		response.ValidationError(c, errs)
		return
	}

	input := usecase.UpdateUserInput{
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	user, err := h.userUseCase.Update(c.Request.Context(), id, input)
	if err != nil {
		if err == domain.ErrNotFound {
			response.NotFound(c, domain.NewAppError("user.not_found"))
			return
		}
		if err == domain.ErrConflict {
			response.Conflict(c, domain.NewAppError("user.email_exists"))
			return
		}
		response.InternalServerError(c, domain.NewAppError("common.error"))
		return
	}

	response.Success(c, dto.ToUserDTO(user))
}

// Delete godoc
// @Summary      Delete user
// @Description  Delete user by ID
// @Tags         Users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        Accept-Language  header  string  false  "Language preference (en, id)"  default(en)
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  response.Response
// @Failure      400  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Failure      404  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Router       /users/{id} [delete]
func (h *UserHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		response.BadRequest(c, domain.NewAppError("request.invalid_id"))
		return
	}

	if err := h.userUseCase.Delete(c.Request.Context(), id); err != nil {
		if err == domain.ErrNotFound {
			response.NotFound(c, domain.NewAppError("user.not_found"))
			return
		}
		response.InternalServerError(c, domain.NewAppError("common.error"))
		return
	}

	response.Success(c, true)
}

// AssignRole godoc
// @Summary      Assign role to user
// @Description  Assign a role to a user
// @Tags         Users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        Accept-Language  header  string  false  "Language preference (en, id)"  default(en)
// @Param        id   path      int  true  "User ID"
// @Param        request body dto.AssignRoleRequest true "Assign role request"
// @Success      200  {object}  response.Response
// @Failure      400  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Failure      404  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Router       /users/{id}/roles [post]
func (h *UserHandler) AssignRole(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := uuid.Parse(idParam)
	if err != nil {
		response.BadRequest(c, domain.NewAppError("request.invalid_id"))
		return
	}

	var req dto.AssignRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, domain.NewAppError("request.invalid_id"))
		return
	}

	if errs := validator.ValidateStruct(req); len(errs) > 0 {
		response.ValidationError(c, errs)
		return
	}

	roleID, err := uuid.Parse(req.RoleID)
	if err != nil {
		response.BadRequest(c, domain.NewAppError("request.invalid_id"))
		return
	}

	if err := h.userUseCase.AssignRole(c.Request.Context(), userID, roleID); err != nil {
		if err == domain.ErrNotFound {
			response.NotFound(c, domain.NewAppError("user.not_found"))
			return
		}
		response.InternalServerError(c, domain.NewAppError("common.error"))
		return
	}

	response.Success(c, true)
}

// RemoveRole godoc
// @Summary      Remove role from user
// @Description  Remove a role from a user
// @Tags         Users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        Accept-Language  header  string  false  "Language preference (en, id)"  default(en)
// @Param        id   path      int  true  "User ID"
// @Param        roleId path    int  true  "Role ID"
// @Success      200  {object}  response.Response
// @Failure      400  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Failure      404  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Router       /users/{id}/roles/{roleId} [delete]
func (h *UserHandler) RemoveRole(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := uuid.Parse(idParam)
	if err != nil {
		response.BadRequest(c, domain.NewAppError("request.invalid_id"))
		return
	}

	roleIDParam := c.Param("roleId")
	roleID, err := uuid.Parse(roleIDParam)
	if err != nil {
		response.BadRequest(c, domain.NewAppError("request.invalid_id"))
		return
	}

	if err := h.userUseCase.RemoveRole(c.Request.Context(), userID, roleID); err != nil {
		if err == domain.ErrNotFound {
			response.NotFound(c, domain.NewAppError("user.not_found"))
			return
		}
		response.InternalServerError(c, domain.NewAppError("common.error"))
		return
	}

	response.Success(c, true)
}
