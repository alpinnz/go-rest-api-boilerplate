package encrypt

import (
	"crypto/sha256"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password, secret string) (string, error) {
	// Gabungkan password + secret
	salted := []byte(password + secret)

	// Hash dengan SHA256 → 32 byte
	digest := sha256.Sum256(salted)

	// Lanjut bcrypt
	hash, err := bcrypt.GenerateFromPassword(digest[:], bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func ComparePassword(hash string, password string, secret string) (bool, error) {
	var (
		salted        = password + secret
		passwordBytes = []byte(salted)
		hashBytes     = []byte(hash)
	)

	err := bcrypt.CompareHashAndPassword(hashBytes, passwordBytes)
	if err != nil {
		return false, nil
	}

	return true, nil
}
