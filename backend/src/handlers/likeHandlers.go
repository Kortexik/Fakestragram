package handlers

import (
	"backend/src/models"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AddLike(c *gin.Context) {
	var json struct {
		UserID int `json:"userID"`
		PostID int `json:"postId"`
	}

	// Bind incoming JSON to struct
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	log.Printf("Request Data: %+v\n", json) // Add this line

	if json.UserID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	like := models.Like{
		UserID: json.UserID,
		PostID: json.PostID,
	}

	success, err := models.AddLike(like, DB)
	if success {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func DeleteLike(c *gin.Context) {
	// Extract postId from URL
	postID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	// Extract username from JWT token (assuming JWT middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	success, err := models.DeleteLike(userID.(int), postID, DB)
	if err != nil || !success {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove like"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Like removed successfully"})
}
