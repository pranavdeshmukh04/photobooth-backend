package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "register"})
}

func Login(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "login"})
}

func Refresh(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "refresh"})
}

func Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "logout"})
}

func Me(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "me"})
}