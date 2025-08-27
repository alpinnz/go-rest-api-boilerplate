package mapper

import (
	"github.com/alpinnz/go-rest-api-boilerplate/internal/domain/entities"
	"github.com/alpinnz/go-rest-api-boilerplate/internal/interfaces/http/dto"
	"github.com/alpinnz/go-rest-api-boilerplate/pkg/utils"
)

func ToAuthSessionResponse(authSession *entities.AuthSession) *dto.AuthSession {
	resp := &dto.AuthSession{
		AccessToken:      authSession.AccessToken,
		AccessExpiresAt:  authSession.AccessExpiresAt,
		RefreshToken:     authSession.RefreshToken,
		RefreshExpiresAt: authSession.RefreshExpiresAt,
	}

	// Conditional mapping between User and UserID
	if authSession.User != nil {
		resp.User = ToUserResponse(authSession.User)
		resp.UserID = nil
	} else {
		resp.UserID = utils.StringPtr(authSession.UserID.String())
		resp.User = nil
	}

	return resp
}

func ToAuthSessionResponses(authSessions []*entities.AuthSession) []*dto.AuthSession {
	out := make([]*dto.AuthSession, 0, len(authSessions))
	for _, r := range authSessions {
		if r == nil {
			continue
		}
		out = append(out, ToAuthSessionResponse(r))
	}
	return out
}
