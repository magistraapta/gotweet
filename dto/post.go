package dto

type PostResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Content  string `json:"content"`
	Comments []CommentResponse
}

type CommentResponse struct {
	Username string `json:"username"`
	Content  string `json:"content"`
}
