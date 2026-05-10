package utils

import (
	"crypto/sha256"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	sum := sha256.Sum256([]byte(password))

	hash, err := bcrypt.GenerateFromPassword(sum[:], bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func VerifyPassword(password string, hash string) error {
	sum := sha256.Sum256([]byte(password))
	return bcrypt.CompareHashAndPassword([]byte(hash), sum[:])
}