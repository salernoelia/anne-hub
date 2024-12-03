// models/models.go

package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Message represents a single message in the conversation history.
type Message struct {
	Sender    string `json:"sender"`
	Content   string `json:"content"`
	Timestamp string `json:"timestamp"`
}

// ConversationData holds the conversation history as a list of messages.
type ConversationData struct {
	Messages []Message `json:"messages"`
}

// ConversationRequest represents the incoming request for a conversation.
type ConversationRequest struct {
	UserID      uuid.UUID `json:"user_id"`
	DeviceID    int       `json:"device_id"`
	RequestPCM  []byte    `json:"request_pcm"`
	Language    string    `json:"language,omitempty"`
	RequestTime string    `json:"request_time"`
}

// Conversation represents a conversation record in the database.
type Conversation struct {
	ID                  int64           `db:"id" json:"id"`
	UserID              uuid.UUID       `db:"user_id" json:"user_id"`
	CreatedAt           time.Time       `db:"created_at" json:"created_at"`
	ConversationHistory json.RawMessage `db:"conversation_history" json:"conversation_history"`
	UpdatedAt           time.Time       `db:"updated_at" json:"updated_at"`
}
