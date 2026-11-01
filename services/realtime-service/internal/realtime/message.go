package realtime

import (
	"time"

	"github.com/google/uuid"
)

const (
	MessageTypeConnected    = "connected"
	MessageTypeNotification = "notification"
)

type Message struct {
	Type string `json:"type"`
	Data any    `json:"data,omitempty"`
}

type ConnectedData struct {
	UserID string `json:"userId"`
}

type Actor struct {
	ID          string `json:"id"`
	Username    string `json:"username,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
}

type Notification struct {
	Type       string    `json:"type"`
	Title      string    `json:"title"`
	Body       string    `json:"body"`
	EntityID   string    `json:"entityId,omitempty"`
	EntityType string    `json:"entityType,omitempty"`
	ActionURL  string    `json:"actionUrl,omitempty"`
	Actor      *Actor    `json:"actor,omitempty"`
	CreatedAt  time.Time `json:"createdAt"`
}

type Sender interface {
	SendToUser(userID uuid.UUID, message Message)
}
