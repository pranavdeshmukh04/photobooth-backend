package utils

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes a password using bcrypt
func HashPassword(password string) (string, error) {
	// Check if password is empty
	if password == "" {
		return "", errors.New("password is empty")
	}
	// Generate bcrypt hash with cost 14
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// ComparePassword compares a hashed password with a plain password
func ComparePassword(hashedPassword, password string) error {
	// Check if hashed password is empty
	if hashedPassword == "" {
		return errors.New("hashed password is empty")
	}
	// Check if password is empty
	if password == "" {
		return errors.New("password is empty")
	}
	// Compare hashed password with plain password
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return err
	}
	return nil
}
