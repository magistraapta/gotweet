package services

import (
	"docker-test/dto"
	"docker-test/internal/store"
	"docker-test/model"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	store *store.UserStore
}

func NewUserService(store *store.UserStore) *UserService {
	return &UserService{store: store}
}

func (s *UserService) GetUserById(id int) (*dto.UserResponse, error) {
	return s.store.GetUserById(id)
}

func (s *UserService) GetAllUsers() (*[]model.User, error) {
	return s.store.GetAllUsers()
}

func (s *UserService) DeleteUserById(id int) error {
	return s.store.DeleteUserbyId(id)
}

func (s *UserService) UpdateUserById(id int, updatedUser *model.User) error {
	return s.store.UpdateUserById(id, *updatedUser)
}

func (s *UserService) Signup(userRequest *model.User) error {
	return s.store.Signup(*userRequest)
}

func (s *UserService) Login(email string, passwrod string) (string, error) {
	user, err := s.store.Login(email)

	if err != nil {
		return "", errors.New("Invalid email")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(passwrod)); err != nil {
		return "", errors.New("Failed to compare password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(10 * time.Minute).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))

	if err != nil {
		return "", errors.New("Failed to create JWT Token")
	}

	return tokenString, nil
}

func (s *UserService) UpdateUser(updatedUser model.User, username string) (*model.User, error) {
	return s.store.UpdateUser(updatedUser, username)
}
