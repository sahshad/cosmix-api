package post

import (
	"time"

	"github.com/google/uuid"
)

const (
	EventVersionOne = "1.0.0"
)

type PostLiked struct {
	EventVersion    string
	PostID          uuid.UUID
	PostOwnerID     uuid.UUID
	ActorID         uuid.UUID
	ActorUsername   string
	ActorDisplayName string
	PostTitle       string
	LikedAt         time.Time
}

type PostCommented struct {
	EventVersion      string
	PostID            uuid.UUID
	PostOwnerID       uuid.UUID
	CommentID         uuid.UUID
	ActorID           uuid.UUID
	ActorUsername     string
	ActorDisplayName  string
	Content           string
	PostTitle         string
	CreatedAt         time.Time
}

type CommentReplied struct {
	EventVersion          string
	PostID                uuid.UUID
	ParentCommentID       uuid.UUID
	ParentCommentOwnerID  uuid.UUID
	ReplyCommentID        uuid.UUID
	ActorID               uuid.UUID
	ActorUsername         string
	ActorDisplayName      string
	Content               string
	PostTitle             string
	CreatedAt             time.Time
}

type CommentLiked struct {
	EventVersion       string
	CommentID          uuid.UUID
	CommentOwnerID     uuid.UUID
	PostID             uuid.UUID
	ActorID            uuid.UUID
	ActorUsername      string
	ActorDisplayName   string
	CommentContent     string
	PostTitle          string
	LikedAt            time.Time
}
