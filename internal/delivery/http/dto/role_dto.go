package dto

import "github.com/alpinnz/go-rest-api-boilerplate/internal/domain/entity"

// CreateRoleRequest represents role creation request from HTTP
type CreateRoleRequest struct {
	Name        string `json:"name" validate:"required,min=2"`
	Description string `json:"description" validate:"omitempty"`
}

// UpdateRoleRequest represents role update request from HTTP
type UpdateRoleRequest struct {
	Name        string `json:"name" validate:"omitempty,min=2"`
	Description string `json:"description" validate:"omitempty"`
}

// RoleDTO represents role data in HTTP responses
type RoleDTO struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// RoleListResponse represents paginated roles response to HTTP
type RoleListResponse struct {
	Roles   []*RoleDTO `json:"roles"`
	Total   int64      `json:"total"`
	Limit   int        `json:"limit"`
	Offset  int        `json:"offset"`
	HasMore bool       `json:"has_more"`
}

// ToRoleDTO converts model.Role to RoleDTO for HTTP response
func ToRoleDTO(role *entity.Role) *RoleDTO {
	if role == nil {
		return nil
	}
	return &RoleDTO{
		ID:          role.ID.String(),
		Name:        role.Name,
		Description: role.Description,
		CreatedAt:   role.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   role.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// ToRoleListResponse converts slice of model.Role to RoleListResponse for HTTP response
func ToRoleListResponse(roles []*entity.Role, total int64, limit, offset int, hasMore bool) *RoleListResponse {
	var roleDTOs []*RoleDTO
	for _, role := range roles {
		roleDTOs = append(roleDTOs, ToRoleDTO(role))
	}

	return &RoleListResponse{
		Roles:   roleDTOs,
		Total:   total,
		Limit:   limit,
		Offset:  offset,
		HasMore: hasMore,
	}
}
