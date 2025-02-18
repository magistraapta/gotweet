package controllers

import (
	"docker-test/initializers"
	"docker-test/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

const bcryptCost = 10

func Signup(c *gin.Context) {
	// get input from body
	var body struct {
		Username string
		Email    string
		Password string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed read take input from body",
		})

		return
	}

	// hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcryptCost)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to hash password",
		})
	}

	// create user
	user := model.User{Username: body.Username, Email: body.Email, Password: string(hash)}
	results := initializers.DB.Create(&user)

	if results.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to create user",
		})
	}

	// respond
	c.JSON(http.StatusCreated, gin.H{
		"message": "User successfully created!",
	})
}

func Login(c *gin.Context) {
	// take email and password from reqbody
	var body struct {
		Email    string
		Password string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to take req body",
		})

		return
	}

	// look up requested user
	var user model.User

	initializers.DB.First(&user, "email = ?", body.Email)

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})

		return
	}

	// compare the password with hashed password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid email or password",
		})
	}
}
