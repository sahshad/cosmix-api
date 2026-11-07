package dto

import "github.com/google/uuid"

type User struct {
	ID          uuid.UUID
	Username    string
	DisplayName string
	AvatarURL   string
	IsVerified  bool
}

type Participant struct {
	ID          uuid.UUID
	Username    string
	DisplayName string
	AvatarURL   string
	IsVerified  bool
	Role        string
	JoinedAt    string
}

type Message struct {
	ID             uuid.UUID
	ConversationID uuid.UUID
	SenderID       uuid.UUID
	Content        string
	CreatedAt      string
	UpdatedAt      string
	Sender         User
	ReadCount      uint
}

type MessageReadReceipt struct {
	User   User
	ReadAt string
}

type Conversation struct {
	ID              uuid.UUID
	Type            string
	Name            string
	AvatarURL       string
	CreatedBy       uuid.UUID
	CreatedAt       string
	UpdatedAt       string
	Participants    []Participant
	LastMessage     *Message
	UnreadCount     uint
	CurrentUserRole string
}

type ConversationListResponse struct {
	Conversations []Conversation
	Pagination    PaginationResponse
}

type MessageListResponse struct {
	Messages   []Message
	Pagination PaginationResponse
}

type MessageReadReceiptListResponse struct {
	Receipts []MessageReadReceipt
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
