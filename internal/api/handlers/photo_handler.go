package handlers

import (
	"net/http"
	
	"github.com/gin-gonic/gin"
)

func UploadPhoto(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "upload photo"})
}

func ListPhotos(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "list photos"})
}

func GetPhoto(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "get photo"})
}

func DeletePhoto(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "delete photo"})
}