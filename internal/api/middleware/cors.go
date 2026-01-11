package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupCORS(allowedOrigins []string) gin.HandlerFunc {
	// TODO: Configure CORS with allowed origins
	// TODO: Set allowed methods (GET, POST, PUT, PATCH, DELETE, OPTIONS)
	// TODO: Set allowed headers (Authorization, Content-Type)
	// TODO: Return CORS middleware

	return cors.Default()
}
