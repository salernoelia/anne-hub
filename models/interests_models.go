package models

import "anne-hub/pkg/uuid"

type Interest struct {
	ID             int       `json:"id" db:"id"`
	UserID         uuid.UUID       `json:"user_id" db:"user_id"`
	CreatedAt      string    `json:"created_at" db:"created_at"`
	UpdatedAt      string    `json:"updated_at" db:"updated_at"`
	Name           string    `json:"name" db:"name"`
	Description    string    `json:"description" db:"description"`
	Level          int       `json:"level" db:"level"`
	LevelAccuracy  int       `json:"level_accuracy" db:"level_accuracy"`
	Users          []User    `json:"users" db:"users"`
}

type Interests struct {
	Interests []Interest `json:"interests"`
}