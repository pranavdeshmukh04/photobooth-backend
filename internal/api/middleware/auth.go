package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/photobooth/backend/config"
	"github.com/photobooth/backend/pkg/jwt"
)

// JWTAuth is a middleware that validates the JWT token
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		// Extract token from Authorization header
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		token := parts[1]

		// Validate token
		claims, err := jwt.ValidateToken(token, config.GetConfig().JWTSecret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		// Set user info in context
		c.Set("userID", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("tier", claims.Tier)

		c.Next()
	}
}
