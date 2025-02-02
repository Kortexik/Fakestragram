package handlers

import (
	"backend/src/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Follow(c *gin.Context) {
	var payload struct {
		FollowerID int `json:"followerId"`
		FolloweeID int `json:"followeeId"`
	}

	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	follow := models.Follow{
		FollowerID: payload.FollowerID,
		FolloweeID: payload.FolloweeID,
	}

	createdFollow, err := models.AddFollow(follow, DB)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, createdFollow)
}

func UnfollowUser(c *gin.Context) {
	var request struct {
		FollowerID int `json:"followerId"`
		FolloweeID int `json:"followeeId"`
	}

	// Bind JSON input
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Check if follow relationship exists before attempting to delete
	exists, err := models.FollowExists(request.FollowerID, request.FolloweeID, DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Follow relationship not found"})
		return
	}

	// Use `DeleteFollow` function to remove follow relationship
	success, err := models.DeleteFollow(request.FollowerID, request.FolloweeID, DB)
	if err != nil || !success {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unfollow user"})
		return
	}

	// Respond with success
	c.JSON(http.StatusOK, gin.H{"message": "Unfollowed successfully"})
}

func CheckFollowStatus(c *gin.Context) {
	var request struct {
		FollowerID int `json:"followerId"`
		FolloweeID int `json:"followeeId"`
	}

	// Bind JSON input
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	fmt.Println(request.FollowerID)
	fmt.Println(request.FolloweeID)

	// Check if follow relationship exists
	exists, err := models.FollowExists(request.FollowerID, request.FolloweeID, DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Respond with follow status
	c.JSON(http.StatusOK, gin.H{"isFollowing": exists})
}
