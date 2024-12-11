package models

import (
	"anne-hub/pkg/uuid"
	"time"
)

// Tasks table
type Task struct {
    ID            int64     `json:"id" db:"id"`
    UserID        uuid.UUID     `json:"user_id" db:"user_id"`
    Title         string    `json:"title" db:"title"`
    Description   *string   `json:"description,omitempty" db:"description"` // Nullable column
    DueDate       *time.Time `json:"due_date,omitempty" db:"due_date"` // Nullable timestamp
    Completed     bool      `json:"completed" db:"completed"`
    CreatedAt     time.Time `json:"created_at" db:"created_at"`
    InterestLinks []string  `json:"interest_links,omitempty" db:"interest_links"` // Array of text
}