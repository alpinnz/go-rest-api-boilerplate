package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func setupTest() {
	SetJWTConfig(
		"test-access-secret-key-min-32chars",
		15*time.Minute,
		"test-refresh-secret-key-min-32char",
		7*24*time.Hour,
	)

	// Ensure tests are isolated even if a previous test changed timeNow.
	timeNow = time.Now
}

func TestGenerateAccessToken(t *testing.T) {
	setupTest()

	userID := uuid.New()
	token, err := GenerateAccessToken(userID)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestGenerateRefreshToken(t *testing.T) {
	setupTest()

	userID := uuid.New()
	token, err := GenerateRefreshToken(userID)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestValidateToken_AccessToken_Success(t *testing.T) {
	setupTest()

	userID := uuid.New()
	token, err := GenerateAccessToken(userID)
	assert.NoError(t, err)

	claims, err := ValidateToken(token)
	assert.NoError(t, err)
	assert.Equal(t, userID.String(), claims.UserID)
	assert.Equal(t, TokenTypeAccess, claims.TokenType)
}

func TestValidateToken_RefreshToken_Success(t *testing.T) {
	setupTest()

	userID := uuid.New()
	token, err := GenerateRefreshToken(userID)
	assert.NoError(t, err)

	claims, err := ValidateToken(token)
	assert.NoError(t, err)
	assert.Equal(t, userID.String(), claims.UserID)
	assert.Equal(t, TokenTypeRefresh, claims.TokenType)
}

func TestValidateToken_InvalidToken(t *testing.T) {
	setupTest()

	_, err := ValidateToken("invalid.token.here")
	assert.Error(t, err)
}

func TestValidateToken_ExpiredToken(t *testing.T) {
	setupTest()

	base := time.Date(2026, 1, 15, 10, 0, 0, 0, time.UTC)
	userID := uuid.New()

	// 1) Generate token at "base" with a short expiration
	SetJWTConfig(
		"test-access-secret-key-min-32chars",
		1*time.Minute,
		"test-refresh-secret-key-min-32char",
		7*24*time.Hour,
	)

	timeNow = func() time.Time { return base }
	token, err := GenerateAccessToken(userID)
	assert.NoError(t, err)

	// 2) Validate after expiration
	timeNow = func() time.Time { return base.Add(2 * time.Minute) }
	_, err = ValidateToken(token)
	assert.Error(t, err)
	assert.ErrorIs(t, err, jwt.ErrTokenExpired)
}

func TestHashPassword(t *testing.T) {
	password := "MySecurePassword123!"
	hash, err := HashPassword(password)

	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.NotEqual(t, password, hash)
}

func TestCheckPasswordHash_Success(t *testing.T) {
	password := "MySecurePassword123!"
	hash, err := HashPassword(password)
	assert.NoError(t, err)

	match := CheckPasswordHash(password, hash)
	assert.True(t, match)
}

func TestCheckPasswordHash_WrongPassword(t *testing.T) {
	password := "MySecurePassword123!"
	hash, err := HashPassword(password)
	assert.NoError(t, err)

	match := CheckPasswordHash("WrongPassword", hash)
	assert.False(t, match)
}
