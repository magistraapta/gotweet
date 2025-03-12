package test

import (
	"docker-test/internal/services"
	"docker-test/internal/store"
	"docker-test/model"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestCreateUser(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&model.User{})

	repo := store.NewUserStore(db)
	services := services.NewUserService(repo)

	mockUser := model.User{Username: "John Pork"}

	t.Run("Signup should not return an error", func(t *testing.T) {
		err := services.Signup(&mockUser)
		assert.Nil(t, err, "Signup should not return an error")
	})

	t.Run("User should be exist in database", func(t *testing.T) {
		var createdUser model.User
		result := db.First(&createdUser, "username = ?", mockUser.Username)
		assert.Nil(t, result.Error, "User should exist in the database")

		// Ensure stored user has expected values
		assert.Equal(t, mockUser.Username, createdUser.Username)
	})
	// Verify user exists in DB

}

func TestGetUserByID(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})

	db.AutoMigrate(&model.User{})
	repo := store.NewUserStore(db)
	services := services.NewUserService(repo)

	mockUser := model.User{ID: 1, Username: "John Pork"}

	db.Create(&mockUser)

	user, err := services.GetUserById(1)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "John Pork", user.Username)
}

func TestGetAllUser(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}

	db.AutoMigrate(&model.User{})

	repo := store.NewUserStore(db)
	services := services.NewUserService(repo)

	// Insert multiple users into the database
	mockUsers := []model.User{
		{ID: 1, Username: "John Pork", Email: "john@email.com"},
		{ID: 2, Username: "Chopped Chin", Email: "chin@email.com"},
	}

	if err := db.Create(&mockUsers).Error; err != nil {
		t.Fatalf("failed to insert mock users: %v", err)
	}

	// Call the service to fetch all users
	users, err := services.GetAllUsers()

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, users)

	for i, expectedUser := range mockUsers {
		assert.Equal(t, expectedUser.ID, (*users)[i].ID)             // Dereference the pointer to access the slice
		assert.Equal(t, expectedUser.Username, (*users)[i].Username) // Assuming Username is not a pointer
		assert.Equal(t, expectedUser.Email, (*users)[i].Email)
	}
}

func TestEditUser(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})

	if err != nil {
		t.Fatal("Failed connect to Database")
	}

	db.AutoMigrate(&model.User{})

	repo := store.NewUserStore(db)
	service := services.NewUserService(repo)

	mockUser := model.User{
		Username: "John pork",
		Email:    "john@email.com",
	}

	db.Create(&mockUser)

	var savedUser model.User

	db.First(&savedUser, "email = ?", mockUser.Email)

	newUsername := "lebron james"

	updatedUser, err := service.UpdateUser(savedUser, newUsername)

	if err != nil {
		t.Fatal("Failed to update user")
	}

	assert.NoError(t, err)
	assert.NotNil(t, updatedUser)

	assert.Equal(t, newUsername, updatedUser.Username)

}

func TestDeleteUser(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatal("Failed to connect to Database")
	}

	db.AutoMigrate(&model.User{})

	repo := store.NewUserStore(db)
	service := services.NewUserService(repo)

	// Step 1: Create and save the user
	mockUser := model.User{
		ID:       1,
		Username: "John pork",
		Email:    "john@email.com",
	}
	db.Create(&mockUser)

	// Step 2: Delete the user by ID
	if err := service.DeleteUserById(1); err != nil {
		t.Fatalf("Failed to delete user")
	}

	// Step 3: Attempt to find the user in the database
	var deletedUser model.User
	err = db.First(&deletedUser, "id = ?", 1).Error

	// Step 4: Assert that the user does not exist anymore
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err) // Ensures the user was not found
}
