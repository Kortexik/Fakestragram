package handlers

import (
	"backend/src/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUnseenNotifications(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	notifications, err := models.GetUnseenNotifications(userID.(int), DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch notifications"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"notifications": notifications})
}

func MarkNotificationsAsSeen(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	err := models.MarkNotificationsAsSeen(userID.(int), DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark notifications as seen"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notifications marked as seen"})
}
