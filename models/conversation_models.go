package models

import (
	"time"
)

type CreateConversationRequest struct {
    UserID   string `json:"user_id" validate:"required,uuid"`
    Request  string `json:"request" validate:"required"`
    Language string `json:"language,omitempty"`
}


type Conversation struct {
	ID             int64     `db:"id" json:"id"`
	UserID         string    `db:"user_id" json:"user_id"`
	ConversationID string    `db:"conversation_id" json:"conversation_id"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
	Request        string    `db:"request" json:"request"`
	Response       string    `db:"response" json:"response"`
	ModelUsed      string    `db:"model_used,omitempty" json:"model_used,omitempty"`
	Role           string    `db:"role,omitempty" json:"role,omitempty"`
}
