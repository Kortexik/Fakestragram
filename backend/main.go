package main

import (
	"backend/src/handlers"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {

	err := handlers.ConnectDatabase()
	handlers.CheckErr(err)
	r := gin.Default()
	r.Use(CORSMiddleware())

	vUsers := r.Group("/users")
	{
		vUsers.GET("", handlers.GetUsers)
		vUsers.GET("/:id", handlers.GetUserById)
		vUsers.POST("adduser", handlers.AddUser)
		vUsers.PUT("updateuser/:id", handlers.UpdateUser)
		vUsers.DELETE("deleteuser/:id", handlers.DeleteUser)
		vUsers.OPTIONS("", handlers.Options)
	}

	// By default it serves on :8080 unless a
	// PORT environment variable was defined.
	vPosts := r.Group("/posts")
	{
		vPosts.GET("", handlers.GetUserPosts)
		vPosts.POST("addpost", handlers.AddUserPost)
		vPosts.POST("upload", handlers.UploadPost)
	}
	r.Run()
}
