package encrypt

import (
	"fmt"
	"time"

	"github.com/alpinnz/go-rest-api-boilerplate/internal/domain/entities"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	UserID uuid.UUID       `json:"user_id"`
	Roles  []entities.Role `json:"roles"`
	jwt.RegisteredClaims
}

// BuildClaims is a helper to construct JWT claims with common fields.
func BuildClaims(userID uuid.UUID, roles []entities.Role, issuedAt, expiresAt time.Time) Claims {
	return Claims{
		UserID: userID,
		Roles:  roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(issuedAt),
			NotBefore: jwt.NewNumericDate(issuedAt),
		},
	}
}

func GenerateHash(claims Claims, secretKey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

func ExtractHash(hash string, secretKey string) (*Claims, error) {
	// Parse dan versifier token
	token, err := jwt.ParseWithClaims(hash, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, fmt.Errorf("invalid token: %v", err)
	}

	// Validasi dan casting claims
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}
