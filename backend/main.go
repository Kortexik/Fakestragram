package main

import (
	"backend/src/handlers"
	"backend/src/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	// gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.LoggingMiddleware())

	vUsers := r.Group("/users")
	{
		vUsers.GET("", handlers.GetUsers)
		vUsers.GET("/:id", handlers.GetUserById)
		vUsers.GET("/getuserprofile/:username", handlers.GetUserProfileDataByUsername) //dodaj do data, liczbe followersow, liczbe followees, liczbe postów
		vUsers.GET("/username/:id", handlers.GetUsernameById)
		vUsers.POST("adduser", handlers.AddUser)
		vUsers.DELETE("deleteuser/:id", handlers.DeleteUser)
		vUsers.OPTIONS("", handlers.Options)
	}

	vPosts := r.Group("/posts")
	{
		vPosts.GET("", handlers.GetUserPosts)
		vPosts.GET("/:id", handlers.GetPostById)
		vPosts.POST("addpost", handlers.AddUserPost)
	}

	vAuth := r.Group("/auth")
	{
		vAuth.POST("/login", middleware.LoginHandler)
		vAuth.POST("/register", handlers.Register)
		vAuth.GET("/check-username", handlers.CheckUsernameExists)
		vAuth.GET("/check-email", handlers.CheckEmailExists)
	}

	vProtected := r.Group("/protected")
	vProtected.Use(middleware.AuthMiddleware())
	{
		vProtected.POST("/upload", handlers.UploadPost)
		vProtected.GET("/home", handlers.ProtectedHandler)
		vProtected.POST("/like", handlers.AddLike)
		vProtected.DELETE("/like/:id", handlers.DeleteLike)
		vProtected.POST("/comment", handlers.AddComment)
		vProtected.DELETE("/comment/:id", handlers.DeleteComment)
		vProtected.GET("/me", handlers.GetCurrentUser)
		vProtected.POST("/isFollowing", handlers.CheckFollowStatus)
		vProtected.POST("/follow", handlers.Follow)
		vProtected.DELETE("/unfollow", handlers.UnfollowUser)
		vProtected.GET("/notifications", handlers.GetUnseenNotifications)
		vProtected.POST("/notifications/mark-seen", handlers.MarkNotificationsAsSeen)
		vProtected.PUT("/update-profile", handlers.UpdateProfile)
	}

	r.Run(":8080")
}
