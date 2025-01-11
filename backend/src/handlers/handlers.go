package handlers

import (
	"backend/src/models"
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

var DB *sql.DB

func ConnectDatabase() error {
	db, err := sql.Open("sqlite3", "./database/data.db")
	if err != nil {
		return err
	}

	DB = db
	return nil
}

func CheckErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

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

func UpdateUser(c *gin.Context) {
	var json models.User

	id := c.Param("id")

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	success, err := models.UpdateUser(id, json, DB)

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

func Options(c *gin.Context) {

	ourOptions := "some options"

	c.String(200, ourOptions)
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
	// Parse the file and other form values
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to upload file: " + err.Error()})
		return
	}
	defer file.Close()

	// Read the file content into a byte slice
	fileContent, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file content"})
		return
	}

	// Retrieve the caption
	caption := c.PostForm("caption")

	// Assuming user ID comes from authentication or a request parameter
	hardcodedUserId := 1
	// Construct the UserPost struct
	newPost := models.UserPost{
		UserID:  hardcodedUserId,
		Media:   fileContent,
		Caption: caption,
	}

	// Add the post to the database
	success, err := models.AddUserPost(newPost, DB)
	if success {
		c.JSON(http.StatusOK, gin.H{"message": "Post uploaded successfully"})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save post: " + err.Error()})
	}
}

func GetUserPosts(c *gin.Context) {
	posts, err := models.GetUserPosts(10, DB)
	CheckErr(err)

	if posts == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No Records Found"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"Posts": posts})
	}
}
