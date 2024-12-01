package models

import (
	"anne-hub/pkg/uuid"
	"time"
)

// Users table
type User struct {
    ID          uuid.UUID     `json:"id" db:"id"`
    Username    string    `json:"username" db:"username"`
    Email       string    `json:"email" db:"email"`
    PasswordHash string    `json:"password_hash" db:"password_hash"`
    CreatedAt   time.Time `json:"created_at" db:"created_at"`
    Age         *int      `json:"age,omitempty" db:"age"` // Nullable column
    Interests   []string  `json:"interests,omitempty" db:"interests"` // Array of text
}

// Devices table
type Device struct {
    ID             int64      `json:"id" db:"id"`
    UserID         int64      `json:"user_id" db:"user_id"`
    DeviceName     string     `json:"device_name" db:"device_name"`
    LastSynced     *time.Time `json:"last_synced,omitempty" db:"last_synced"` // Nullable timestamp
    CompanionAppID *int64     `json:"companion_app_id,omitempty" db:"companion_app_id"` // Nullable foreign key
    CreatedAt      time.Time  `json:"created_at" db:"created_at"`
}

// Tasks table
type Task struct {
    ID            int64     `json:"id" db:"id"`
    UserID        int64     `json:"user_id" db:"user_id"`
    Title         string    `json:"title" db:"title"`
    Description   *string   `json:"description,omitempty" db:"description"` // Nullable column
    DueDate       *time.Time `json:"due_date,omitempty" db:"due_date"` // Nullable timestamp
    Completed     bool      `json:"completed" db:"completed"`
    CreatedAt     time.Time `json:"created_at" db:"created_at"`
    InterestLinks []string  `json:"interest_links,omitempty" db:"interest_links"` // Array of text
}

// CompanionApps table
type CompanionApp struct {
    ID        int64           `json:"id" db:"id"`
    Settings  map[string]any  `json:"settings" db:"settings"` // JSONB column
    CreatedAt time.Time       `json:"created_at" db:"created_at"`
    UserID    int64           `json:"user_id" db:"user_id"`
}
