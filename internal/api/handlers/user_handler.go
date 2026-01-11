package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetProfile(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "get profile"})
}

func UpdateProfile(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "update profile"})
}

func GetStats(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "get stats"})
}