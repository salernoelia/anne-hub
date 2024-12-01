package models

import (
	"anne-hub/pkg/uuid"
	"time"
)

// Devices table
type Device struct {
    ID             int64      `json:"id" db:"id"`
    UserID         uuid.UUID      `json:"user_id" db:"user_id"`
    DeviceName     string     `json:"device_name" db:"device_name"`
    LastSynced     *time.Time `json:"last_synced,omitempty" db:"last_synced"` // Nullable timestamp
    CompanionAppID *int64     `json:"companion_app_id,omitempty" db:"companion_app_id"` // Nullable foreign key
    CreatedAt      time.Time  `json:"created_at" db:"created_at"`
}

