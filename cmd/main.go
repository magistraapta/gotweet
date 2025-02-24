package main

import (
	"docker-test/controllers"
	"docker-test/initializers"
	"docker-test/middleware"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

const (
	MaxUploadSize   = 8 << 20
	FileDirectory   = "static/upload"
	filePermissions = 0o755
)

func init() {
	initializers.LoadEnvFile()
	initializers.ConnectToDatabase()
	initializers.SyncDatabase()
}

func main() {
	router := gin.Default()

	if err := os.MkdirAll(FileDirectory, filePermissions); err != nil {
		panic(err)
	}

	router.Static("/static", "./static")

	v1 := router.Group("/v1")
	{
		v1.POST("/signup", controllers.Signup)
		v1.POST("/upload", controllers.HandleImageUpload)
		v1.POST("/login", controllers.Login)
		v1.GET("/validate", middleware.RequireAuth, controllers.Validate)
		user := v1.Group("/user")
		{
			user.PUT("/", middleware.RequireAuth, controllers.UpdateUser)
			user.GET("/", controllers.GetAllUsers)
			user.GET("/:id", controllers.GetUserById)
			user.DELETE("/:id", controllers.DeleteUserById)
			user.PUT("/:id", controllers.EditUserById)
		}

		post := v1.Group("/post")
		{
			post.POST("/", middleware.RequireAuth, controllers.CreatePost)
		}

	}

	router.Use(middleware.CorsMiddleware())
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
}
