package dto

import "time"

type MediaItem struct {
	PublicID string `json:"public_id" binding:"required"`
	URL      string `json:"url" binding:"required"`
	Type     string `json:"type" binding:"required"`
	Duration int    `json:"duration"`
}

type Media struct {
	ID        uint       `json:"id"`
	PostID    uint       `json:"post_id"`
	PublicID  string     `json:"public_id"`
	URL       string     `json:"url"`
	Type      string     `json:"type"`
	Duration  *int        `json:"duration"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

type User struct {
	ID          uint       `json:"id"`
	Email       string     `json:"email"`
	Username    string     `json:"username"`
	DisplayName string     `json:"display_name"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

type Like struct {
	ID        uint      `json:"id"`
	PostID    uint      `json:"post_id"`
	AuthorID  uint      `json:"author_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Comment struct {
	ID        uint      `json:"id"`
	PostID    uint      `json:"post_id"`
	AuthorID  uint      `json:"author_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PostList struct {
	ID            uint       `json:"id"`
	Content       string     `json:"content"`
	LikesCount    int        `json:"likes_count"`
	CommentsCount int        `json:"comments_count"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at"`
	User          User       `json:"user"`
	Media     []Media    `json:"media"`
	// Likes     []Like      `json:"likes"`
	// Comments  []Comment   `json:"comments"`
}

type PostListResponse struct {
	Posts      []PostList         `json:"posts"`
	Pagination PaginationResponse `json:"pagination"`
}

type CommentList struct {
	ID        uint      `json:"id"`
	PostID    uint      `json:"post_id"`
	AuthorID  uint      `json:"author_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CommentListResponse struct {
	Comments   []CommentList
	Pagination PaginationResponse
}

type CreatePostRequest struct {
	Content string      `json:"content" binding:"required"`
	Media   []MediaItem `json:"media"`
}

type UpdatePostRequest struct {
	Content string      `json:"content" binding:"required"`
	Media   []MediaItem `json:"media"`
}

type CreateCommentRequest struct {
	Content string `json:"content" binding:"required"`
}

type UpdateCommentRequest struct {
	Content string `json:"content" binding:"required"`
}

type PaginationRequest struct {
	Page  int64 `json:"page"`
	Limit int64 `json:"limit"`
}

type PaginationResponse struct {
	TotalCount int64 `json:"total_count"`
	Page       int64 `json:"page"`
	Limit      int64 `json:"limit"`
	TotalPages int64 `json:"total_pages"`
}
