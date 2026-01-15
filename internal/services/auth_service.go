package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/photobooth/backend/config"
	"github.com/photobooth/backend/internal/models"
	"github.com/photobooth/backend/internal/repositories"
	"github.com/photobooth/backend/pkg/jwt"
	"github.com/photobooth/backend/pkg/utils"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailAlreadyExists = errors.New("email already exists")
)

// AuthService handles authentication business logic
type AuthService struct {
	userRepo *repositories.UserRepository
	config   *config.Config
}

// NewAuthService creates a new auth service
func NewAuthService(userRepo *repositories.UserRepository, config *config.Config) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		config:   config,
	}
}

// CreateUser creates a new user
func (s *AuthService) CreateUser(email, password, name string) (*models.UserResponse, error) {
	// Check if user already exists
	existingUser, err := s.userRepo.GetByEmail(email)
	if err != nil && !errors.Is(err, repositories.ErrUserNotFound) {
		return nil, fmt.Errorf("failed to check user existence: %w", err)
	}
	if existingUser != nil {
		return nil, ErrEmailAlreadyExists
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &models.User{
		Email:        email,
		Name:         name,
		PasswordHash: hashedPassword,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	userResponse := user.ToUserResponse()
	return &userResponse, nil
}

// Login authenticates a user and returns tokens
func (s *AuthService) Login(email, password string) (*models.UserResponse, string, string, error) {
	// Get user by email
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			return nil, "", "", ErrInvalidCredentials
		}
		return nil, "", "", fmt.Errorf("failed to get user: %w", err)
	}

	// Compare password
	if err := utils.ComparePassword(user.PasswordHash, password); err != nil {
		return nil, "", "", ErrInvalidCredentials
	}

	// Generate tokens
	accessToken, refreshToken, err := s.GenerateTokens(user.ID, user.Email, user.SubscriptionTier)
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to generate tokens: %w", err)
	}

	userResponse := user.ToUserResponse()
	return &userResponse, accessToken, refreshToken, nil
}

// GenerateTokens generates access and refresh tokens for a user
func (s *AuthService) GenerateTokens(userID, email, tier string) (string, string, error) {
	// Generate access token
	accessToken, err := jwt.GenerateAccessToken(
		userID,
		email,
		tier,
		s.config.JWTSecret,
		s.config.JWTAccessExpiry,
	)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generate refresh token ID
	tokenID := fmt.Sprintf("refresh::%s", uuid.New().String()[:12])

	// Generate refresh token
	refreshToken, err := jwt.GenerateRefreshToken(
		userID,
		s.config.JWTSecret,
		s.config.JWTRefreshExpiry,
	)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Store refresh token in database
	refreshTokenModel := models.RefreshToken{
		TokenID:     tokenID,
		HashedToken: refreshToken,
		ExpiresAt:   time.Now().Add(s.config.JWTRefreshExpiry),
		CreatedAt:   time.Now(),
	}

	if err := s.userRepo.AddRefreshToken(userID, refreshTokenModel); err != nil {
		return "", "", fmt.Errorf("failed to store refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}

// RefreshTokens generates new tokens using a refresh token
func (s *AuthService) RefreshTokens(refreshToken string) (string, string, error) {
	// Validate refresh token
	claims, err := jwt.ValidateToken(refreshToken, s.config.JWTSecret)
	if err != nil {
		return "", "", fmt.Errorf("invalid refresh token: %w", err)
	}

	// Get user
	user, err := s.userRepo.GetByID(claims.UserID)
	if err != nil {
		return "", "", fmt.Errorf("failed to get user: %w", err)
	}

	// Verify refresh token exists in database
	tokenFound := false
	for _, token := range user.RefreshTokens {
		// Compare refresh token
		if token.HashedToken == refreshToken {
			// Check if token is expired
			if time.Now().After(token.ExpiresAt) {
				return "", "", errors.New("refresh token expired")
			}
			tokenFound = true
			break
		}
	}

	if !tokenFound {
		return "", "", errors.New("refresh token not found")
	}

	// Generate new tokens
	accessToken, newRefreshToken, err := s.GenerateTokens(user.ID, user.Email, user.SubscriptionTier)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate tokens: %w", err)
	}

	return accessToken, newRefreshToken, nil
}

// Logout invalidates a refresh token
func (s *AuthService) Logout(userID, refreshToken string) error {
	// Get user
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Find and remove refresh token
	var newTokens []models.RefreshToken
	for _, token := range user.RefreshTokens {
		// Compare refresh token
		if token.HashedToken != refreshToken {
			// Keep tokens that don't match
			newTokens = append(newTokens, token)
		}
	}

	user.RefreshTokens = newTokens

	if err := s.userRepo.Update(user); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// GetUserByID retrieves a user by ID
func (s *AuthService) GetUserByID(userID string) (*models.UserResponse, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	userResponse := user.ToUserResponse()
	return &userResponse, nil
}
