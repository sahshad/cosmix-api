package dto

import (
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
	CreatedAt string
	UpdatedAt string
}

type User struct {
	ID          uuid.UUID
	Email       string
	Username    string
	DisplayName string
	AvatarURL   string
	CreatedAt   string
	UpdatedAt   string
}

type Like struct {
	ID        uuid.UUID
	PostID    uuid.UUID
	AuthorID  uuid.UUID
	CreatedAt string
	UpdatedAt string
}

type Comment struct {
	ID        uuid.UUID
	PostID    uuid.UUID
	AuthorID  uuid.UUID
	Content   string
	CreatedAt string
	UpdatedAt string
}

type PostList struct {
	ID            uuid.UUID
	Content       string
	LikesCount    int
	CommentsCount int
	SharesCount   int
	IsLiked       bool
	IsOwner       bool
	CreatedAt     string
	UpdatedAt     string
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
	ID              uuid.UUID
	PostID          uuid.UUID
	AuthorID        uuid.UUID
	ParentCommentID *uuid.UUID
	Content         string
	LikesCount      int
	RepliesCount    int
	IsLiked         bool
	IsOwner         bool
	CreatedAt       string
	UpdatedAt       string
	User            User
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
	Content         string
	ParentCommentID *uuid.UUID
}

type UpdateCommentRequest struct {
	Content string
}

type PaginationRequest struct {
	Page  int32
	Limit int32
}

type PaginationResponse struct {
	TotalCount int32
	Page       int32
	Limit      int32
	TotalPages int32
}
