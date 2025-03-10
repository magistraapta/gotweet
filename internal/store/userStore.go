package store

import (
	"docker-test/dto"
	"docker-test/model"
	"log"

	"gorm.io/gorm"
)

type UserStore struct {
	DB *gorm.DB
}

func NewUserStore(db *gorm.DB) *UserStore {
	return &UserStore{DB: db}
}

func (s *UserStore) GetUserById(id int) (*dto.UserResponse, error) {
	var user model.User

	result := s.DB.Find(&user, id)

	if result.Error != nil {
		return nil, result.Error
	}

	userResponse := dto.UserResponse{
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
	}

	return &userResponse, nil
}

func (s *UserStore) GetAllUsers() (*[]model.User, error) {
	var users []model.User

	result := s.DB.Find(&users)

	if result.Error != nil {
		log.Println("Failed to run query to get all users")
		return nil, result.Error
	}

	return &users, nil
}

func (s *UserStore) DeleteUserbyId(id int) error {
	var user model.User

	result := s.DB.Delete(&user, id)

	if result.Error != nil {
		log.Println("Failed to run query to delete user by id")
		return result.Error

	}

	log.Print("User deleted successfully")
	return nil
}

func (s *UserStore) UpdateUserById(id int, updatedUser model.User) error {
	// find the existing user
	var user model.User
	result := s.DB.Find(&user, id)

	if result.Error != nil {
		log.Print("Failed to find user")
		return result.Error
	}

	// update user
	result = s.DB.Model(&user).Updates(updatedUser)

	if result.Error != nil {
		log.Print("Failed to update user")
		return result.Error
	}

	log.Print("User updated successfully")
	return nil
}

func (s *UserStore) Signup(userRequest model.User) error {

	result := s.DB.Create(&userRequest)

	if result.Error != nil {
		log.Print("Failed to create user")
		return result.Error
	}

	log.Print("User created successfully")
	return nil
}

func (s *UserStore) Login(email string) (*model.User, error) {
	// find user by email
	var user model.User
	if err := s.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *UserStore) UpdateUser(updatedUser model.User, username string) (*model.User, error) {

	result := s.DB.Model(&updatedUser).Updates(map[string]interface{}{
		"username": username,
	})

	if result.Error != nil {
		log.Print("Failed to run query: Update user")
		return nil, result.Error
	}

	return &updatedUser, nil
}
