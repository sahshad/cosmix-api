package dto

import "github.com/google/uuid"

type User struct {
	ID          uuid.UUID
	Username    string
	DisplayName string
	AvatarURL   string
	IsVerified  bool
}

type Message struct {
	ID             uuid.UUID
	ConversationID uuid.UUID
	SenderID       uuid.UUID
	Content        string
	CreatedAt      string
	UpdatedAt      string
	Sender         User
}

type Conversation struct {
	ID           uuid.UUID
	Type         string
	CreatedAt    string
	UpdatedAt    string
	Participants []User
	LastMessage  *Message
	UnreadCount  uint
}

type ConversationListResponse struct {
	Conversations []Conversation
	Pagination    PaginationResponse
}

type MessageListResponse struct {
	Messages   []Message
	Pagination PaginationResponse
}

type PaginationRequest struct {
	Page  uint
	Limit uint
}

type PaginationResponse struct {
	TotalCount uint
	Page       uint
	Limit      uint
	TotalPages uint
}
