package auth

import "golang.org/x/crypto/bcrypt"

// HashPassword generates a bcrypt hash from plaintext password.
// Uses bcrypt cost factor of 14 for strong security (recommended for 2024+).
// Higher cost = more secure but slower. Cost 14 takes ~1 second on modern hardware.
// Returns hashed password string or error if hashing fails.
//
// Example:
//
//	hashedPassword, err := auth.HashPassword("user_password")
//	if err != nil {
//	    return err
//	}
//	// Store hashedPassword in database
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash verifies plaintext password against bcrypt hash.
// Performs constant-time comparison to prevent timing attacks.
// Returns true if password matches hash, false otherwise.
// Safe to use in authentication flows.
//
// Example:
//
//	isValid := auth.CheckPasswordHash("user_password", storedHash)
//	if !isValid {
//	    return errors.New("invalid credentials")
//	}
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
