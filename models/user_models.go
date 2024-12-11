package models

import "anne-hub/pkg/uuid"

type User struct {
    ID           uuid.UUID `json:"id"`
    Username     string `json:"username"`
    FirstName    string `json:"first_name"`
    LastName     string `json:"last_name"`
    Email        string `json:"email"`
    PasswordHash string `json:"password_hash"`
    CreatedAt    string `json:"created_at"`
    Age          int    `json:"age"`
    Country      string `json:"country"`
    City         string `json:"city"`
}

type UserDetails struct {
    ID        uuid.UUID `json:"id"`
    FirstName string `json:"first_name"`
    LastName  string `json:"last_name"`
    CreatedAt string `json:"created_at"`
    Age       int    `json:"age"`
    Email     string `json:"email"`
    City      string `json:"city"`
    Country   string `json:"country"`
}

type UserData struct {
    User      UserDetails `json:"user"`
    Interests []Interest  `json:"interests"`
    Tasks     []Task      `json:"tasks"`
}
