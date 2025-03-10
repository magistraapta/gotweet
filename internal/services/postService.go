package services

import (
	"docker-test/dto"
	"docker-test/internal/store"
	"docker-test/model"
)

type PostService struct {
	store *store.PostStore
}

func NewPostService(store *store.PostStore) *PostService {
	return &PostService{store: store}
}

func (s *PostService) CreatePost(content string, userID uint) (*model.Post, error) {
	return s.store.CreatePost(content, userID)
}

func (s *PostService) GetAllPost() (*[]dto.PostResponse, error) {
	return s.store.GetAllPost()
}

func (s *PostService) GetPostById(id int) (*dto.PostResponse, error) {
	return s.store.GetPostById(id)
}
