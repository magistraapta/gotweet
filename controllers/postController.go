package controllers

import (
	"docker-test/initializers"
	"docker-test/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreatePost(c *gin.Context) {
	// get current user

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

	// get req body
	var body struct {
		Content string `json:"content" binding:"required"`
	}

	if c.ShouldBindBodyWithJSON(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": "failed to read take input",
		})
		return
	}

	// create post
	post := model.Post{UserID: user.ID, Content: body.Content}

	result := initializers.DB.Create(&post)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create post",
		})
		return
	}

	// send response

	c.JSON(http.StatusCreated, gin.H{
		"message": "Post successfully created",
		"post": gin.H{
			"id":        post.ID,
			"userId":    post.UserID,
			"content":   post.Content,
			"createdAt": post.CreatedAt,
			"updatedAt": post.UpdatedAt,
		},
	})
}
