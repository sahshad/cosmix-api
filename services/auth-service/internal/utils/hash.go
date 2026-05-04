package utils

import (
	"crypto/sha256"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes passwords or tokens safely
// If the input is longer than 72 bytes (like refresh tokens), this still works
func HashPassword(password string) (string, error) {
	// Step 1: SHA256 pre-hash
	sum := sha256.Sum256([]byte(password))

	// Step 2: bcrypt hash of SHA256 sum
	hash, err := bcrypt.GenerateFromPassword(sum[:], bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// VerifyPassword compares a plaintext password or token with a stored bcrypt hash
func VerifyPassword(password string, hash string) error {
	sum := sha256.Sum256([]byte(password))
	return bcrypt.CompareHashAndPassword([]byte(hash), sum[:])
}