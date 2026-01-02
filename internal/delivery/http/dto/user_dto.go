package dto

import "github.com/alpinnz/go-rest-api-boilerplate/internal/domain/entity"

// UpdateUserRequest represents user update request from HTTP
type UpdateUserRequest struct {
	Email     string `json:"email" validate:"omitempty,email"`
	FirstName string `json:"first_name" validate:"omitempty,min=2"`
	LastName  string `json:"last_name" validate:"omitempty,min=2"`
}

// AssignRoleRequest represents role assignment request from HTTP
type AssignRoleRequest struct {
	RoleID string `json:"role_id" validate:"required,uuid"`
}

// UserListResponse represents paginated users response to HTTP
type UserListResponse struct {
	Users   []*UserDetailDTO `json:"users"`
	Total   int64            `json:"total"`
	Limit   int              `json:"limit"`
	Offset  int              `json:"offset"`
	HasMore bool             `json:"has_more"`
}

// UserDetailDTO represents detailed user data with roles in HTTP responses
type UserDetailDTO struct {
	ID        string     `json:"id"`
	Email     string     `json:"email"`
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
	FullName  string     `json:"full_name"`
	Roles     []*RoleDTO `json:"roles"`
	CreatedAt string     `json:"created_at"`
	UpdatedAt string     `json:"updated_at"`
}

// ToUserDetailDTO converts model.User to UserDetailDTO for HTTP response
func ToUserDetailDTO(user *entity.User) *UserDetailDTO {
	if user == nil {
		return nil
	}

	var roles []*RoleDTO
	for _, role := range user.Roles {
		roles = append(roles, ToRoleDTO(role))
	}

	return &UserDetailDTO{
		ID:        user.ID.String(),
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		FullName:  user.FullName(),
		Roles:     roles,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// ToUserListResponse converts slice of model.User to UserListResponse for HTTP response
func ToUserListResponse(users []*entity.User, total int64, limit, offset int, hasMore bool) *UserListResponse {
	var userDTOs []*UserDetailDTO
	for _, user := range users {
		userDTOs = append(userDTOs, ToUserDetailDTO(user))
	}

	return &UserListResponse{
		Users:   userDTOs,
		Total:   total,
		Limit:   limit,
		Offset:  offset,
		HasMore: hasMore,
	}
}
