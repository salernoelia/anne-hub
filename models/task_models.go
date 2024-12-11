package models

import (
	"anne-hub/pkg/uuid"
	"time"
)

// Task represents a task in the database
type Task struct {
    ID          int64 `db:"id"`
    UserID      uuid.UUID `db:"user_id"`
    Title       string    `db:"title"`
    Description string    `db:"description"`
    DueDate     time.Time `db:"due_date"`
    Completed   bool      `db:"completed"`
    CreatedAt   time.Time `db:"created_at"`
}
