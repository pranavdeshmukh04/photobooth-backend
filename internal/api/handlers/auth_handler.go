package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/photobooth/backend/internal/services"
	"github.com/photobooth/backend/pkg/validator"
)

// AuthHandler contains the dependencies for auth handlers
type AuthHandler struct {
	authService *services.AuthService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// RegisterRequest represents the request body for the register endpoint
type RegisterRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// LoginRequest represents the request body for the login endpoint
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// RefreshRequest represents the request body for the refresh endpoint
type RefreshRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

// LogoutRequest represents the request body for the logout endpoint
type LogoutRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

// Register handles the registration of a new user
func (h *AuthHandler) Register(c *gin.Context) {
	var registerRequest RegisterRequest
	if err := c.ShouldBindJSON(&registerRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate request
	if err := validator.ValidateEmail(registerRequest.Email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := validator.ValidatePassword(registerRequest.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := validator.ValidateName(registerRequest.Name); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create user
	user, err := h.authService.CreateUser(registerRequest.Email, registerRequest.Password, registerRequest.Name)
	if err != nil {
		if errors.Is(err, services.ErrEmailAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Generate tokens
	accessToken, refreshToken, err := h.authService.GenerateTokens(user.ID, user.Email, user.SubscriptionTier)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"user": user,
		"tokens": gin.H{
			"accessToken":  accessToken,
			"refreshToken": refreshToken,
		},
	})
}

// Login handles user login
func (h *AuthHandler) Login(c *gin.Context) {
	var loginRequest LoginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate request
	if err := validator.ValidateEmail(loginRequest.Email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := validator.ValidatePassword(loginRequest.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Login user
	user, accessToken, refreshToken, err := h.authService.Login(loginRequest.Email, loginRequest.Password)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
		"tokens": gin.H{
			"accessToken":  accessToken,
			"refreshToken": refreshToken,
		},
	})
}

// Refresh handles token refresh
func (h *AuthHandler) Refresh(c *gin.Context) {
	var refreshRequest RefreshRequest
	if err := c.ShouldBindJSON(&refreshRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Refresh tokens
	accessToken, refreshToken, err := h.authService.RefreshTokens(refreshRequest.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}

// Logout handles user logout
func (h *AuthHandler) Logout(c *gin.Context) {
	var logoutRequest LogoutRequest
	if err := c.ShouldBindJSON(&logoutRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from context (set by JWT middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Logout user
	if err := h.authService.Logout(userID.(string), logoutRequest.RefreshToken); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// Me returns the current user's information
func (h *AuthHandler) Me(c *gin.Context) {
	// Get user ID from context (set by JWT middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get user
	user, err := h.authService.GetUserByID(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}
