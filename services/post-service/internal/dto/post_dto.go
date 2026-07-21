package dto

import (
	"time"

	"github.com/google/uuid"
)

type MediaItem struct {
	PublicID string
	URL      string
	Type     string
	Duration int
}

type Media struct {
	ID        uuid.UUID
	PostID    uuid.UUID
	PublicID  string
	URL       string
	Type      string
	Duration  *int
	CreatedAt time.Time
	UpdatedAt *time.Time
}

type User struct {
	ID          uuid.UUID
	Email       string
	Username    string
	DisplayName string
	CreatedAt   time.Time
	UpdatedAt   *time.Time
}

type Like struct {
	ID        uuid.UUID
	PostID    uuid.UUID
	AuthorID  uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Comment struct {
	ID        uuid.UUID
	PostID    uuid.UUID
	AuthorID  uuid.UUID
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type PostList struct {
	ID            uuid.UUID
	Content       string
	LikesCount    int
	CommentsCount int
	CreatedAt     time.Time
	UpdatedAt     *time.Time
	User          User
	Media         []Media
	// Likes     []Like
	// Comments  []Comment
}

type PostListResponse struct {
	Posts      []PostList
	Pagination PaginationResponse
}

type CommentList struct {
	ID        uuid.UUID
	PostID    uuid.UUID
	AuthorID  uuid.UUID
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CommentListResponse struct {
	Comments   []CommentList
	Pagination PaginationResponse
}

type CreatePostRequest struct {
	Content string
	Media   []MediaItem
}

type UpdatePostRequest struct {
	Content string
	Media   []MediaItem
}

type CreateCommentRequest struct {
	Content string
}

type UpdateCommentRequest struct {
	Content string
}

type PaginationRequest struct {
	Page  int64
	Limit int64
}

type PaginationResponse struct {
	TotalCount int64
	Page       int64
	Limit      int64
	TotalPages int64
}
