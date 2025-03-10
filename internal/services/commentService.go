package services

import (
	"docker-test/internal/store"
	"docker-test/model"
)

type CommentService struct {
	store *store.CommentStore
}

func NewCommentService(store *store.CommentStore) *CommentService {
	return &CommentService{store: store}
}

func (s *CommentService) CreateComment(userID uint, postID uint, content string) (*model.Comment, error) {
	return s.store.CreateComment(userID, postID, content)
}
