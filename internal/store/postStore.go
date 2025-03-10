package store

import (
	"docker-test/dto"
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

func (s *PostStore) GetAllPost() (*[]dto.PostResponse, error) {
	var posts []model.Post

	result := s.DB.Preload("User").Preload("Comments.User").Find(&posts)

	if result.Error != nil {
		log.Print("Failed to execute query: Get all posts")
		return nil, result.Error
	}

	postResponses := make([]dto.PostResponse, len(posts))
	for i, post := range posts {
		postResponses[i] = dto.PostResponse{
			ID:       post.ID,
			Username: post.User.Username, // Ensure the post author's username is included
			Content:  post.Content,
			Comments: make([]dto.CommentResponse, len(post.Comments)),
		}

		// Map comments
		for j, comment := range post.Comments {
			postResponses[i].Comments[j] = dto.CommentResponse{
				Username: comment.User.Username, // Ensure comment author's username is included
				Content:  comment.Content,
			}
		}
	}

	return &postResponses, nil
}

func (s *PostStore) GetPostById(id int) (*dto.PostResponse, error) {
	var post model.Post

	err := s.DB.Preload("User").Preload("Comments.User").Find(&post, id).Error

	if err != nil {
		return nil, err
	}

	postResponse := dto.PostResponse{
		ID:       post.ID,
		Username: post.User.Username, // Only include username
		Content:  post.Content,
		Comments: make([]dto.CommentResponse, len(post.Comments)),
	}

	// Map comments content only
	for i, comment := range post.Comments {
		postResponse.Comments[i] = dto.CommentResponse{
			Username: comment.User.Username,
			Content:  comment.Content,
		}
	}

	return &postResponse, nil
}
