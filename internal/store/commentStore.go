package store

import (
	"docker-test/model"

	"gorm.io/gorm"
)

type CommentStore struct {
	DB gorm.DB
}

func NewCommentStore(db *gorm.DB) *CommentStore {
	return &CommentStore{DB: *db}
}

func (s *CommentStore) CreateComment(userID uint, postID uint, content string) (*model.Comment, error) {
	comment := model.Comment{
		UserID:  userID,
		PostID:  postID,
		Content: content,
	}

	if err := s.DB.Create(&comment).Error; err != nil {
		return nil, err
	}

	return &comment, nil
}
