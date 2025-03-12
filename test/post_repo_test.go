package test

import (
	"docker-test/internal/services"
	"docker-test/internal/store"
	"docker-test/model"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestCreatePost(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&model.Post{})

	repo := store.NewPostStore(db)
	services := services.NewPostService(repo)

	mockPost := model.Post{Content: "this is new post"}

	t.Run("Create post should not return error", func(t *testing.T) {
		_, err := services.CreatePost(mockPost.Content, 1)

		assert.Nil(t, err, "CreatePost should not return an error")
	})

	t.Run("New Post should be exist in database", func(t *testing.T) {
		var post model.Post
		result := db.First(&post, "content = ?", mockPost.Content)

		assert.Nil(t, result.Error, "Post should exist in the database")
		assert.Equal(t, mockPost.Content, post.Content, "Stored post content should match")
	})

}

func TestGetPostById(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&model.Post{}, &model.Comment{})

	repo := store.NewPostStore(db)
	services := services.NewPostService(repo)

	// create mock post
	mockPost := model.Post{
		ID:      1,
		Content: "This is mock post",
	}

	db.Create(&mockPost)

	post, err := services.GetPostById(1)

	if err != nil {
		t.Fatal("failed to get post")
	}

	assert.NoError(t, err)
	assert.NotNil(t, post)
	assert.Equal(t, mockPost.Content, post.Content)
}

func TestGetAllPosts(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}

	db.AutoMigrate(&model.Post{}, &model.Comment{})

	repo := store.NewPostStore(db)
	services := services.NewPostService(repo)

	// Insert multiple users into the database

	mockPosts := []model.Post{
		{
			ID:      1,
			Content: "this is new post",
		},
		{
			ID:      2,
			Content: "this is mock post",
		},
	}

	if err := db.Create(&mockPosts).Error; err != nil {
		t.Fatalf("failed to insert mock Posts: %v", err)
	}

	// Call the service to fetch all users
	posts, err := services.GetAllPost()

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, posts)

	for i, expectedUser := range mockPosts {
		assert.Equal(t, expectedUser.ID, (*posts)[i].ID) // Dereference the pointer to access the slice
		assert.Equal(t, expectedUser.Content, (*posts)[i].Content)
	}
}
