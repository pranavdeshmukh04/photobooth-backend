package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Claims represents the claims for the JWT token
type Claims struct {
	UserID string `json:"sub"`
	Email  string `json:"email"`
	Tier   string `json:"tier"`
	jwt.RegisteredClaims
}

// GenerateAccessToken generates an access token for the user
func GenerateAccessToken(userID, email, tier, secret string, expiry time.Duration) (string, error) {
	// Check if userID, email, tier, secret is empty
	if userID == "" || email == "" || tier == "" || secret == "" {
		return "", errors.New("userID, email, tier, secret cannot be empty")
	}

	// Create new claims
	newClaims := &Claims{
		UserID: userID,
		Email:  email,
		Tier:   tier,
		RegisteredClaims: jwt.RegisteredClaims{
			ID: uuid.New().String(),
			Subject: userID,
			IssuedAt: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			Issuer: "photobooth",
		},
	}

	// Create new token with claims	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
	// Sign token with secret
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// GenerateRefreshToken generates a refresh token for the user
func GenerateRefreshToken(userID, secret string, expiry time.Duration) (string, error) {
	// Check if userID, secret is empty
	if userID == "" || secret == "" {
		return "", errors.New("userID, secret cannot be empty")
	}

	// Create new claims
	refreshClaims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ID: uuid.New().String(),
			Subject: userID,
			IssuedAt: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			Issuer: "photobooth",
		},
	}

	// Create new token with claims	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	// Sign token with secret
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// ValidateToken validates a token
func ValidateToken(tokenString, secret string) (*Claims, error) {
    // Parse token with claims
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

    // Type assert and return
    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims, nil
    }
    
    return nil, jwt.ErrSignatureInvalid
}
