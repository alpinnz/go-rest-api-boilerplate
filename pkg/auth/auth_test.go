package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGenerateAccessToken(t *testing.T) {
	SetJWTSecret("test-secret-key")

	userID := uuid.New().String()
	token, err := GenerateAccessToken(userID)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestValidateAccessToken_Success(t *testing.T) {
	SetJWTSecret("test-secret-key")

	userID := uuid.New().String()
	token, err := GenerateAccessToken(userID)
	assert.NoError(t, err)

	claims, err := ValidateAccessToken(token)
	assert.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
}

func TestValidateAccessToken_InvalidToken(t *testing.T) {
	SetJWTSecret("test-secret-key")

	_, err := ValidateAccessToken("invalid.token.here")
	assert.Error(t, err)
}

func TestValidateAccessToken_ExpiredToken(t *testing.T) {
	SetJWTSecret("test-secret-key")

	// This would require mocking time, skipped for simplicity
	// In real tests, you'd use a library like github.com/benbjohnson/clock
	t.Skip("Requires time mocking")
}

func TestGenerateRefreshToken(t *testing.T) {
	SetJWTSecret("test-secret-key")

	userID := uuid.New().String()
	token, err := GenerateRefreshToken(userID)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestHashPassword(t *testing.T) {
	password := "MySecurePassword123!"
	hash, err := HashPassword(password)

	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.NotEqual(t, password, hash)
}

func TestCheckPassword_Success(t *testing.T) {
	password := "MySecurePassword123!"
	hash, err := HashPassword(password)
	assert.NoError(t, err)

	err = CheckPassword(hash, password)
	assert.NoError(t, err)
}

func TestCheckPassword_WrongPassword(t *testing.T) {
	password := "MySecurePassword123!"
	hash, err := HashPassword(password)
	assert.NoError(t, err)

	err = CheckPassword(hash, "WrongPassword")
	assert.Error(t, err)
}

func TestCheckPassword_EmptyPassword(t *testing.T) {
	hash, err := HashPassword("password")
	assert.NoError(t, err)

	err = CheckPassword(hash, "")
	assert.Error(t, err)
}
