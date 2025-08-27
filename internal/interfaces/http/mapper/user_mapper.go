package mapper

import (
	"github.com/alpinnz/go-rest-api-boilerplate/internal/domain/entities"
	"github.com/alpinnz/go-rest-api-boilerplate/internal/interfaces/http/dto"
)

func ToUserResponse(user *entities.User) *dto.User {
	if user == nil {
		return nil
	}
	return &dto.User{
		ID:          user.ID.String(),
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Email:       user.Email,
		Roles:       ToRoleResponses(user.Roles, false),
		ActivatedAt: user.ActivatedAt,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}
}

func ToUserResponses(users []*entities.User) []dto.User {
	out := make([]dto.User, 0, len(users))
	for _, r := range users {
		if r == nil {
			continue
		}
		out = append(out, *ToUserResponse(r)) // dereference pointer
	}
	return out
}
