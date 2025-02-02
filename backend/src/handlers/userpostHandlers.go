package handlers

import (
	"backend/src/models"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUserPosts(c *gin.Context) {
	posts, err := models.GetUserPosts(DB)
	CheckErr(err)

	if posts == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No Records Found"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"Posts": posts})
	}
}

func GetPostById(c *gin.Context) {
	id := c.Param("id")
	post, err := models.GetUserPostById(id, DB)
	CheckErr(err)

	if post.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("No post with id '%s' found", id)})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": post})

}

func AddUserPost(c *gin.Context) {
	var json models.UserPost

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	success, err := models.AddUserPost(json, DB)

	if success {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
}

func UploadPost(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to upload file: " + err.Error()})
		return
	}
	defer file.Close()

	fileContent, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file content"})
		return
	}

	caption := c.PostForm("caption")
	userID, err := models.GetUserIdByUsername(username.(string), DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user ID"})
		return
	}

	newPost := models.UserPost{
		UserID:  userID,
		Media:   fileContent,
		Caption: caption,
	}

	success, err := models.AddUserPost(newPost, DB)
	if err != nil || !success {
		log.Println("Error adding user post:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save post: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post uploaded successfully"})
}
