package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupCORS is a middleware that sets up CORS
func SetupCORS(allowedOrigins []string) gin.HandlerFunc {
	// Check if allowed origins are provided
	if len(allowedOrigins) == 0 {
		allowedOrigins = []string{"*"}  // Allow all in development
	}

	// Create CORS configuration
	corsConfig := cors.DefaultConfig()

	// Set allowed origins, methods, headers and credentials
	corsConfig.AllowOrigins = allowedOrigins
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE"}
	corsConfig.AllowHeaders = []string{"Authorization", "Content-Type"}
	corsConfig.AllowCredentials = true

	// Create and return CORS middleware
	return cors.New(corsConfig)
}
