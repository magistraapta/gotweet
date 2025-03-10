package handler

import (
	"docker-test/internal/services"
	"docker-test/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CommentHandler struct {
	services *services.CommentService
}

func NewCommentHandler(services *services.CommentService) *CommentHandler {
	return &CommentHandler{services: services}
}

func (h *CommentHandler) CreateComment(c *gin.Context) {
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

	// get postID

	postID, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": "Failed to get postID from req body",
		})
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

	comment, err := h.services.CreateComment(user.ID, uint(postID), body.Content)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "Failed to create post"})
	}

	c.JSON(http.StatusOK, gin.H{"msg": "Post created successfully", "comment": gin.H{
		"id":      comment.ID,
		"userId":  comment.UserID,
		"postID":  comment.PostID,
		"content": comment.Content,
	},
	})

}
