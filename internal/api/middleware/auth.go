package middleware

import (
	"github.com/gin-gonic/gin"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Extract token from Authorization header
		// TODO: Validate token
		// TODO: Extract user info from claims
		// TODO: Set user info in context
		// TODO: Continue to next handler or abort with 401

		c.Next()
	}
}
