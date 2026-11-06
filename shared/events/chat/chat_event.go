package chat

import (
	"time"

	"github.com/google/uuid"
)

const EventVersionOne = "1.0.0"

type ChatMessageSent struct {
	EventVersion      string      `json:"event_version"`
	MessageID         uuid.UUID   `json:"message_id"`
	ConversationID    uuid.UUID   `json:"conversation_id"`
	SenderID          uuid.UUID   `json:"sender_id"`
	SenderUsername    string      `json:"sender_username"`
	SenderDisplayName string      `json:"sender_display_name"`
	SenderAvatarURL   *string     `json:"sender_avatar_url,omitempty"`
	RecipientIDs      []uuid.UUID `json:"recipient_ids"`
	Content           string      `json:"content"`
	CreatedAt         time.Time   `json:"created_at"`
}
