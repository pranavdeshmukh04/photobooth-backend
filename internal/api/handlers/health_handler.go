package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func HealthLive(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "live"})
}

func HealthReady(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ready"})
}