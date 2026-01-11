package middleware

import (
	"github.com/gin-gonic/gin"
)

func Logging() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Use gin.Recovery() or custom recovery
		// TODO: Log panics
		// TODO: Return 500 error on panic

		c.Next()
	}
}
