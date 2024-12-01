package models

import (
	"anne-hub/pkg/uuid"
	"time"
)

// CompanionApps table
type CompanionApp struct {
    ID        int64           `json:"id" db:"id"`
    Settings  map[string]any  `json:"settings" db:"settings"` // JSONB column
    CreatedAt time.Time       `json:"created_at" db:"created_at"`
    UserID    uuid.UUID           `json:"user_id" db:"user_id"`
}
