package mapper

import (
	"github.com/alpinnz/go-rest-api-boilerplate/internal/domain/entities"
	"github.com/alpinnz/go-rest-api-boilerplate/internal/interfaces/http/dto"
)

func ToRoleResponse(role *entities.Role, withMeta bool) *dto.Role {
	if role == nil {
		return nil
	}
	dtoRole := dto.Role{
		ID:   role.ID.String(),
		Name: role.Name,
	}
	if withMeta {
		dtoRole.CreatedAt = &role.CreatedAt
		dtoRole.UpdatedAt = &role.UpdatedAt
	}
	return &dtoRole
}

func ToRoleResponses(roles []*entities.Role, withMeta bool) []*dto.Role {
	if roles == nil {
		return nil
	}
	out := make([]*dto.Role, 0, len(roles))
	for _, r := range roles {
		if r == nil {
			continue
		}
		out = append(out, ToRoleResponse(r, withMeta))
	}
	return out
}
