package main

import (
	"docker-test/controllers"
	"docker-test/initializers"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type PageData struct {
	Title    string
	ImageURL string
	Error    string
}

const (
	MaxUploadSize   = 8 << 20
	FileDirectory   = "files"
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

	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "./static")

	router.POST("/signup", controllers.Signup)
	router.POST("/upload", handleImageUpload)
	router.GET("/", renderView)

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
}

func handleImageUpload(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		c.HTML(http.StatusBadRequest, "index.html", gin.H{
			"error": "Failed to upload file",
		})
	}

	if file.Size > MaxUploadSize {
		c.HTML(http.StatusInternalServerError, "index.html", gin.H{
			"error": "file size above limit",
		})
	}

	if err := c.SaveUploadedFile(file, "static/upload/"+file.Filename); err != nil {
		c.HTML(http.StatusBadRequest, "index.html", gin.H{
			"error": "Failed to upload file",
		})
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"image": "static/upload/" + file.Filename,
	})
}

func renderView(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{"Title": "Upload file"})
}
