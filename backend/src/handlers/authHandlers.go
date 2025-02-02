package handlers

import (
	"backend/src/models"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var creds struct {
		Username  string
		Password  string
		FirstName string
		LastName  string
		Email     string
	}

	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := models.User{
		Username:  creds.Username,
		Password:  []byte(creds.Password),
		FirstName: creds.FirstName,
		LastName:  creds.LastName,
		Email:     creds.Email}

	success, err := models.AddUser(user, DB)
	if !success || err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to register user", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

func CheckUsernameExists(c *gin.Context) {
	username := c.Query("username")

	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE username = ? LIMIT 1)`
	err := DB.QueryRow(query, username).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"exists": exists})
}

func CheckEmailExists(c *gin.Context) {
	username := c.Query("email")

	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = ? LIMIT 1)`
	err := DB.QueryRow(query, username).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"exists": exists})
}

func GetCurrentUser(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"username": username, "userID": userID})
}

func ProtectedHandler(c *gin.Context) {
	username, _ := c.Get("username")
	c.JSON(http.StatusOK, gin.H{"message": "Welcome " + username.(string)})
}
