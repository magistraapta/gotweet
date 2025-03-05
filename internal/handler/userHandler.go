package handler

import (
	"docker-test/internal/services"
	"docker-test/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	services *services.UserService
}

func NewUserHandler(services *services.UserService) *UserHandler {
	return &UserHandler{services: services}
}

func (h *UserHandler) GetUserById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": "failed to get requestbody",
		})
		return
	}
	user, err := h.services.GetUserById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": "failed to get user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":  "Success getting user",
		"user": user,
	})
}

func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.services.GetAllUsers()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": err.Error(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":   "success",
		"users": users,
	})
}

func (h *UserHandler) DeleteUserById(c *gin.Context) {
	idString := c.Param("id")
	id, err := strconv.Atoi(idString)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": "Failed to get user id",
		})
	}

	result := h.services.DeleteUserById(id)

	if result != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": "Failed to delete user" + result.Error(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "Success delete user",
	})
}

func (h *UserHandler) UpdateUserById(c *gin.Context) {
	idString := c.Param("id")
	id, err := strconv.Atoi(idString)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": "Invalid user ID",
		})

		return
	}

	var updatedUser model.User
	if err := c.ShouldBind(&updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "Invalid req body",
		})
		return
	}

	if err := h.services.UpdateUserById(id, &updatedUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": "Failed to update user",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "Successfully update user",
	})

}

func (h *UserHandler) Signup(c *gin.Context) {
	// get req body

	var userRequest model.User

	if err := c.ShouldBindJSON(&userRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": "Invalid req body",
		})

		return
	}

	// hash user password
	hash, err := bcrypt.GenerateFromPassword([]byte(userRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to hash password",
		})
		return
	}

	userRequest.Password = string(hash)
	if err := h.services.Signup(&userRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User successfully created!",
	})
}

func (h *UserHandler) Login(c *gin.Context) {
	// get req body
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": "Invalid request body",
		})

		return
	}

	token, err := h.services.Login(body.Email, body.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", token, 3600*24*30, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{"token": token})

}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	user, exist := c.Get("user")

	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Not Authenticated",
		})

		return
	}

	var body struct {
		Username string
	}

	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	currentUser := user.(model.User)

	updatedUser, err := h.services.UpdateUser(currentUser, body.Username)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to update user" + err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User updated",
		"user":    updatedUser,
	})
}
