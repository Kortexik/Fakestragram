package handlers

import (
	"backend/src/models"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetUsers(c *gin.Context) {
	users, err := models.GetUsers(10, DB)
	CheckErr(err)

	if users == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No Records Found"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"data": users})
	}
}

func GetUserById(c *gin.Context) {
	id := c.Param("id")
	user, err := models.GetUserById(id, DB)
	CheckErr(err)

	if user.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("No user with id '%s' found", id)})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"data": user})
	}
}

func GetUsernameById(c *gin.Context) {
	id := c.Param("id")
	user, err := models.GetUserById(id, DB)
	CheckErr(err)

	if user.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("No user with id '%s' found", id)})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"data": user.Username})
	}
}

func GetUserProfileDataByUsername(c *gin.Context) {
	username := c.Param("username")
	user, err := models.GetUserProfileDataByUsername(username, DB)
	CheckErr(err)
	userPosts, err := models.GetUserPostsByUserId(user.ID, DB)
	CheckErr(err)
	numberOfFollowers, err := models.NumberOfFollowers(user.ID, DB)
	CheckErr(err)
	numberOfFollowees, err := models.NumberOfFollowees(user.ID, DB)
	CheckErr(err)

	if user.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("No user with username '%s' found", username)})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"user": user, "userPosts": userPosts, "numberOfFollowers": numberOfFollowers, "numberOfFollowees": numberOfFollowees})
	}
}

func AddUser(c *gin.Context) {

	var json models.User

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	success, err := models.AddUser(json, DB)

	if success {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
}
func DeleteUser(c *gin.Context) {
	id := c.Param("id")

	result, err := models.DeleteUser(id, DB)
	CheckErr(err)
	if result {
		c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Success, removed user with id '%s'.", id)})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
}

func UpdateProfile(c *gin.Context) {
	_, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	file, _, err := c.Request.FormFile("avatar")
	if err != nil && err != http.ErrMissingFile {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to upload file: " + err.Error()})
		return
	}
	defer func() {
		if file != nil {
			file.Close()
		}
	}()

	var fileContent []byte
	if file != nil {
		fileContent, err = io.ReadAll(file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file content"})
			return
		}
	}

	bio := c.Request.FormValue("bio")
	userID := c.Request.FormValue("id")
	id, err := strconv.Atoi(userID)
	CheckErr(err)

	if fileContent != nil {
		success, err := models.UpdateProfilePic(models.User{
			ID:         id,
			ProfilePic: fileContent,
		}, DB)
		if err != nil || !success {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile picture"})
			return
		}
	}

	success, err := models.UpdateBio(models.User{
		ID:  id,
		Bio: &bio,
	}, DB)
	if err != nil || !success {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update bio"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})

}

func Options(c *gin.Context) {

	ourOptions := "some options"

	c.String(200, ourOptions)
}
