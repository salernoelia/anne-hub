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
