package services

import (
	"anne-hub/models"
	"anne-hub/pkg/db"
	"anne-hub/pkg/uuid"
	"database/sql"
	"fmt"
	"log"
)

func FetchUserData(userID uuid.UUID) (models.UserData, error) {
    var user models.User
    query := "SELECT id, username, email, password_hash, created_at, age, country, city, first_name, last_name FROM users WHERE id = $1"

    // Fetch user details
    err := db.DB.QueryRow(query, userID).Scan(
        &user.ID,
        &user.Username,
        &user.Email,
        &user.PasswordHash,
        &user.CreatedAt,
        &user.Age,
        &user.Country,
        &user.City,
        &user.FirstName,
        &user.LastName,
    )
    if err != nil {
        if err == sql.ErrNoRows {
            log.Printf("No user found with ID: %s", userID)
            return models.UserData{}, fmt.Errorf("user with ID %s not found", userID)
        }
        log.Printf("Error fetching user data: %v", err)
        return models.UserData{}, fmt.Errorf("error fetching user data: %v", err)
    }

    // Fetch user interests
    interestsQuery := `
        SELECT id, user_id, created_at, updated_at, name, description, level, level_accuracy
        FROM interests
        WHERE user_id = $1
    `
    rows, err := db.DB.Query(interestsQuery, userID)
    if err != nil {
        log.Printf("Error fetching user interests: %v", err)
        return models.UserData{}, fmt.Errorf("error fetching user interests: %v", err)
    }
    defer rows.Close()
	
    var interests []models.Interest


    for rows.Next() {
        var interest models.Interest
        err := rows.Scan(
            &interest.ID,
            &interest.UserID,
            &interest.CreatedAt,
            &interest.UpdatedAt,
            &interest.Name,
            &interest.Description,
            &interest.Level,
            &interest.LevelAccuracy,
        )
        if err != nil {
            log.Printf("Error scanning interest row: %v", err)
            continue
        }
        interests = append(interests, interest)
    }

	fmt.Println("interests",interests)

    if err = rows.Err(); err != nil {
        log.Printf("Error iterating over interest rows: %v", err)
        return models.UserData{}, fmt.Errorf("error iterating over interests data: %v", err)
    }

    // Convert models.User to models.UserDetails
    userDetails := models.UserDetails{
        ID:        user.ID,
        FirstName: user.FirstName,
        LastName:  user.LastName,
        CreatedAt: user.CreatedAt,
        Age:       user.Age,
        Email:     user.Email,
        City:      user.City,
        Country:   user.Country,
    }

    // Combine user details and interests
    return models.UserData{
        User:      userDetails,
        Interests: interests,
    }, nil
}