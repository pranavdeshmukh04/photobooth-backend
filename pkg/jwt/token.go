package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID string `json:"sub"`
	Email  string `json:"email"`
	Tier   string `json:"tier"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(userID, email, tier, secret string, expiry time.Duration) (string, error) {
	// TODO: Create claims
	// TODO: Set expiry
	// TODO: Sign token
	// TODO: Return token string

	return "", nil
}

func GenerateRefreshToken(userID, tokenID, secret string, expiry time.Duration) (string, error) {
	// TODO: Create refresh token claims
	// TODO: Sign token
	// TODO: Return token string

	return "", nil
}

func ValidateToken(tokenString, secret string) (*Claims, error) {
	// TODO: Parse token
	// TODO: Validate signature
	// TODO: Check expiry
	// TODO: Return claims

	return nil, nil
}
