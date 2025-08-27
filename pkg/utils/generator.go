package utils

import (
	"context"
	"errors"
	"time"

	"github.com/alpinnz/go-rest-api-boilerplate/pkg/constants"
	"github.com/alpinnz/go-rest-api-boilerplate/pkg/encrypt"
	normalizer "github.com/dimuska139/go-email-normalizer/v4"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// GenerateUUID function to generate uuid
func GenerateUUID() uuid.UUID {
	id, err := uuid.NewV7()
	if err != nil {
		return uuid.New()
	}
	return id
}

// ValidateUUID function to validate uuid
func ValidateUUID(id string) bool {
	err := uuid.Validate(id)
	if err != nil {
		return false
	}
	return true
}

// EmailNormalizer function to normalizer email
func EmailNormalizer(email string) string {
	n := normalizer.NewNormalizer()
	return n.Normalize(email)
}

func ParseUUID(s string) (uuid.UUID, error) {
	if s == "" {
		return uuid.Nil, errors.New("UUID string is empty")
	}

	id, err := uuid.Parse(s)
	if err != nil {
		return uuid.Nil, errors.New("invalid UUID format")
	}

	return id, nil
}

func NullableTime(t gorm.DeletedAt) *time.Time {
	if t.Valid {
		return &t.Time
	}
	return nil
}

func Fallback[T comparable](val, fallback T, zero T) T {
	if val != zero {
		return val
	}
	return fallback
}

func FallbackPtr[T any](val *T, fallback T) T {
	if val != nil {
		return *val
	}
	return fallback
}

// SetClaims impairment claims ke context
func SetClaims(ctx context.Context, claims *encrypt.Claims) context.Context {
	return context.WithValue(ctx, constants.AuthSession, claims)
}

// GetClaims mengambil claims dari context
func GetClaims(ctx context.Context) (*encrypt.Claims, error) {
	claims, ok := ctx.Value(constants.AuthSession).(*encrypt.Claims)
	if !ok || claims == nil {
		return nil, errors.New("unauthorized: claims not found or invalid")
	}
	return claims, nil
}

// GetSession parse gin context to get current session
//func GetSession(c *gin.Context) *entity.Session {
//	val, exist := c.Get(constant.AuthCookieSession)
//	if !exist {
//		return nil
//	}
//
//	if sess, ok := val.(*entity.Session); ok {
//		return sess
//	}
//
//	return nil
//}
