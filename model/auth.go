package model

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type Login struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type Register struct {
	Username        string `json:"username" binding:"required"`
	Email           string `json:"email"`
	Password        string `json:"password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
}

type TokenOwner struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	ID       string `json:"id"`
}

// VerifyPassword verifies if the given password matches the stored hash.
func ValidateUserPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// HashPassword generates a bcrypt hash for the given password.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	return string(bytes), nil
}
