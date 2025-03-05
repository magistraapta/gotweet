package handler

import (
	"docker-test/internal/services"
	"docker-test/model"
	"net/http"

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

	c.JSON(http.StatusOK, gin.H{"msg": "Post created successfully", "post": gin.H{
		"id":        post.ID,
		"userId":    post.UserID,
		"content":   post.Content,
		"createdAt": post.CreatedAt,
		"updatedAt": post.UpdatedAt,
	},
	})
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
