package main

import (
	"docker-test/initializers"
	"docker-test/internal/handler"
	"docker-test/internal/services"
	"docker-test/internal/store"
	"docker-test/middleware"
	"log"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvFile()
	initializers.ConnectToDatabase()
	initializers.SyncDatabase()
}

func main() {
	router := gin.Default()

	//connect to db
	db, err := initializers.ConnectToDatabase()

	if err != nil {
		log.Fatal("Failed to connect database")
	}

	userStore := store.NewUserStore(db)
	userService := services.NewUserService(userStore)
	userHandler := handler.NewUserHandler(userService)

	postStore := store.NewPostStore(db)
	postService := services.NewPostService(postStore)
	postHandler := handler.NewPostHandler(postService)

	commentStore := store.NewCommentStore(db)
	commentService := services.NewCommentService(commentStore)
	commentHandler := handler.NewCommentHandler(commentService)

	router.Use(middleware.CORSMiddleware())

	v1 := router.Group("/v1")
	{
		v1.POST("/signup", userHandler.Signup)
		v1.POST("/login", userHandler.Login)
		v1.GET("/validate", middleware.RequireAuth, middleware.ValidateUser)
		user := v1.Group("/user")
		{
			user.PUT("/", middleware.RequireAuth, userHandler.UpdateUser)
			user.GET("/", userHandler.GetAllUsers)
			user.GET("/:id", userHandler.GetUserById)
			user.DELETE("/:id", userHandler.DeleteUserById)
			user.PUT("/:id", userHandler.UpdateUserById)
		}

		post := v1.Group("/post")
		{
			post.POST("/", middleware.RequireAuth, postHandler.CreatePost)
			post.GET("/", postHandler.GetAllPost)
			post.GET("/:id", postHandler.GetPostById)
		}

		comment := v1.Group("/comment")
		{
			comment.POST("/:id", middleware.RequireAuth, commentHandler.CreateComment)
			// comment.GET("/", middleware.RequireAuth, commentHandler.GetUserID)
		}

	}

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
}
