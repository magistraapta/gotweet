package handler

import (
	"docker-test/internal/services"
	"docker-test/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PostHandler struct {
	services *services.PostService
}

func NewPostHandler(services *services.PostService) *PostHandler {
	return &PostHandler{services: services}
}

func (h *PostHandler) CreatePost(c *gin.Context) {
	userObj, exist := c.Get("user")

	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{
			"err": "Not Authenticated",
		})
		return
	}

	user, ok := userObj.(model.User)

	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid user ID format",
		})
		return
	}

	var body struct {
		Content string `json:"content" binding:"required"`
	}

	if c.ShouldBindBodyWithJSON(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": "failed to read take input",
		})
		return
	}

	post, err := h.services.CreatePost(body.Content, user.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "Failed to create post"})
	}

	c.JSON(http.StatusOK, gin.H{"msg": "Post created successfully", "post": post})
}

func (h *PostHandler) GetAllPost(c *gin.Context) {
	posts, err := h.services.GetAllPost()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": err.Error(),
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":   "success",
		"posts": posts,
	})
}

func (h *PostHandler) GetPostById(c *gin.Context) {
	postID, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "Invalid post ID",
		})
		return
	}

	post, err := h.services.GetPostById(postID)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	c.JSON(http.StatusOK, post)
}
