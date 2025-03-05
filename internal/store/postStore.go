package store

import (
	"docker-test/model"
	"log"

	"gorm.io/gorm"
)

type PostStore struct {
	DB gorm.DB
}

func NewPostStore(db *gorm.DB) *PostStore {
	return &PostStore{DB: *db}
}

func (s *PostStore) CreatePost(content string, userID uint) (*model.Post, error) {

	post := model.Post{UserID: userID, Content: content}

	if err := s.DB.Create(&post).Error; err != nil {
		return nil, err
	}

	return &post, nil
}

func (s *PostStore) GetAllPost() (*[]model.Post, error) {
	var posts []model.Post

	result := s.DB.Preload("User").Find(&posts)

	if result.Error != nil {
		log.Print("Failed to execute query: Get all posts")
		return nil, result.Error
	}

	return &posts, nil
}
