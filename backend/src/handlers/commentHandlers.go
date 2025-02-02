package handlers

import (
	"backend/src/models"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddComment(c *gin.Context) {
	var json struct {
		UserID         int    `json:"userID"`
		PostID         int    `json:"postId"`
		CommentContent string `json:"commentContent"`
	}

	// Bind incoming JSON to struct
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	log.Printf("Request Data: %+v\n", json)

	if json.UserID == 0 || json.CommentContent == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	comment := models.Comment{
		UserID:  json.UserID,
		PostID:  json.PostID,
		Content: json.CommentContent,
	}

	createdComment, err := models.AddComment(comment, DB)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Return the full comment object as a response
	c.JSON(http.StatusOK, createdComment)
}

func DeleteComment(c *gin.Context) {
	id := c.Param("id")

	result, err := models.DeleteComment(id, DB)
	CheckErr(err)
	if result {
		c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Success, removed comment with id '%s'.", id)})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
}
