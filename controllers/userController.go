package controllers

import (
	"docker-test/initializers"
	"docker-test/model"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

const (
	MaxUploadSize   = 8 << 20
	FileDirectory   = "files"
	filePermissions = 0o755
	bcryptCost      = 10
)

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

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create token: " + err.Error(),
		})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
	})
}

func Validate(c *gin.Context) {

	user, _ := c.Get("user")

	c.JSON(http.StatusOK, gin.H{
		"message": "Im logged in!",
		"user":    user,
	})
}

func UpdateUser(c *gin.Context) {
	// get current user
	user, exist := c.Get("user")

	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Not Authenticated",
		})

		return
	}

	// get user req body
	var body struct {
		Username string
	}

	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	// update user
	currentUser := user.(model.User)

	result := initializers.DB.Model(&currentUser).Updates(map[string]interface{}{
		"username": body.Username,
	})

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to update user" + result.Error.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User updated",
		"user":    currentUser,
	})

}

func HandleImageUpload(c *gin.Context) {
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

func GetAllUsers(c *gin.Context) {
	var users []model.User

	result := initializers.DB.Find(&users)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":   "Failed to retrieve users",
			"error": result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":   "Success retrieved all users",
		"users": users, // Return the users array, not the result object
	})
}

func GetUserById(c *gin.Context) {
	// get user id
	id := c.Param("id")

	// get user from db by id
	var user model.User

	result := initializers.DB.First(&user, id)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	// response
	c.JSON(http.StatusOK, gin.H{
		"msg":  "Successfully getting user",
		"user": user,
	})
}

func DeleteUserById(c *gin.Context) {
	// get user id
	id := c.Param("id")

	// query to delete user
	var user model.User
	result := initializers.DB.Delete(&user, id)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to delete user: " + result.Error.Error(),
		})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "User not found",
		})
	}
	// send response
	c.JSON(http.StatusOK, gin.H{
		"msg": "Successfully delete user",
	})
}

func EditUserById(c *gin.Context) {
	// Get user ID from URL parameter
	id := c.Param("id")

	// Parse request body
	var body struct {
		Username string
	}
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	// Find user first
	var user model.User
	result := initializers.DB.First(&user, id)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	// Update the user
	result = initializers.DB.Model(&user).Update("username", body.Username)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	// Send response
	c.JSON(http.StatusOK, gin.H{
		"msg":  "Successfully edited user",
		"user": user,
	})
}
