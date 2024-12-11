package services

import (
	"anne-hub/models"
	"anne-hub/pkg/db"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

func HandleProcessConversationInput(
	pcmData []byte, headers models.WSRequestHeaders) (models.AnneWearConversationRequest, error) {

	if len(pcmData) == 0 {
		return models.AnneWearConversationRequest{}, errors.New("no PCM data received")
	}

	if headers.XUserID == "" || headers.XDeviceID == "" || headers.XLanguage == "" {
		log.Println("Missing required headers:", headers.XUserID, headers.XDeviceID, headers.XLanguage)
		return models.AnneWearConversationRequest{}, errors.New("missing required headers")
	}

	if headers.XLanguage != "en" && headers.XLanguage != "de" {
		log.Println("Invalid language:", headers.XLanguage)
		return models.AnneWearConversationRequest{}, errors.New("invalid language")
	}

	userID, err := uuid.Parse(headers.XUserID)
	if err != nil {
		log.Println("Invalid user ID:", headers.XUserID)
		return models.AnneWearConversationRequest{}, errors.New("invalid user ID")
	}

	deviceID, err := strconv.Atoi(headers.XDeviceID)
	if err != nil {
		log.Println("Invalid device ID:", headers.XDeviceID)
		return models.AnneWearConversationRequest{}, errors.New("invalid device ID")
	}

	req := models.AnneWearConversationRequest{
		UserID:      userID,
		DeviceID:    deviceID,
		RequestPCM:  pcmData,
		Language:    headers.XLanguage,
	}

	return req, nil
}

// checks if a previous conversation exists and retrieves it.
func GetPreviousConversation(userID uuid.UUID, resetMinutes int) (*models.Conversation, models.ConversationHistory, error) {
	var lastConversation models.Conversation
	var conversationHistory models.ConversationHistory

	query := `
        SELECT id, user_id, conversation_history, created_at
        FROM conversations
        WHERE user_id = $1
          AND created_at >= NOW() - $2 * INTERVAL '1 minute'
        ORDER BY created_at DESC
        LIMIT 1;
    `

	// log.Printf("Executing SQL Query with UserID=%s, ResetMinutes=%d", userID, resetMinutes)
	err := db.DB.QueryRow(query, userID, resetMinutes).Scan(
		&lastConversation.ID,
		&lastConversation.UserID,
		&lastConversation.ConversationHistory,
		&lastConversation.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("No previous conversation found within the reset time")
			return nil, models.ConversationHistory{}, nil
		}
		log.Printf("Error querying conversation: %v\n", err)
		return nil, models.ConversationHistory{}, err
	}

	// log.Printf("Previous conversation found: %+v\n", lastConversation)

	// Unmarshal conversation history if present.
	if len(lastConversation.ConversationHistory) > 0 {
		err = json.Unmarshal(lastConversation.ConversationHistory, &conversationHistory)
		if err != nil {
			log.Printf("Error unmarshalling ConversationHistory: %v\n", err)
			return nil, models.ConversationHistory{}, err
		}
	}

	return &lastConversation, conversationHistory, nil
}

// updates an existing conversation in the database.
func UpdateExistingConversation(convoID int64, convoJSON []byte) error {
	// log.Printf("Updating existing conversation ID: %d\n", convoID)
	updateQuery := `
		UPDATE conversations
		SET conversation_history = $1, updated_at = NOW()
		WHERE id = $2
		RETURNING updated_at;
	`

	// log.Printf("Executing UPDATE Query with ConversationHistory: %s and ID: %d", string(convoJSON), convoID)

	var updatedAt time.Time
	err := db.DB.QueryRow(updateQuery, convoJSON, convoID).Scan(&updatedAt)
	if err != nil {
		log.Printf("Error updating conversation: %v\n", err)
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  "Failed to update conversation.",
			Internal: err,
		}
	}
	// log.Printf("Conversation ID: %d updated at %v\n", convoID, updatedAt)
	return nil
}

//  inserts a new conversation into the database.
func InsertNewConversation(userID uuid.UUID, convoJSON []byte) error {
	log.Println("Inserting new conversation into the database")
	insertQuery := `
		INSERT INTO conversations (user_id, conversation_history)
		VALUES ($1, $2)
		RETURNING id, created_at;
	`

	// log.Printf("Executing INSERT Query with UserID: %s and ConversationHistory: %s", userID, string(convoJSON))

	var newID int64
	var createdAt time.Time
	err := db.DB.QueryRow(insertQuery, userID, convoJSON).Scan(&newID, &createdAt)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code.Name() == "foreign_key_violation" {
			log.Println("Foreign key violation: Invalid user_id")
			return &echo.HTTPError{
				Code:    http.StatusBadRequest,
				Message: "Invalid user_id. User does not exist.",
			}
		}
		log.Printf("Error inserting conversation: %v\n", err)
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  "Failed to store conversation.",
			Internal: err,
		}
	}
	// log.Printf("New conversation inserted with ID: %d at %v\n", newID, createdAt)
	return nil
}


// appends a new message to the conversation history.
func AppendMessageToConversationHistory(history *models.ConversationHistory, sender, content string) {
	message := models.Message{
		Sender:    sender,
		Content:   content,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
	history.Messages = append(history.Messages, message)
	// log.Printf("Appended %s message to conversation history: %+v\n", sender, message)
}

