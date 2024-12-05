// handlers/conversation_handler.go

package handlers

import (
	"anne-hub/models"
	"anne-hub/pkg/db"
	"anne-hub/pkg/groq"
	"anne-hub/pkg/pcm"
	"anne-hub/pkg/systemprompt"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

var conversationResetMinutes = 15

func ConversationHandler(c echo.Context) error {
	log.Println("Entered ConversationHandler")
	c.Logger().Info("Starting ConversationHandler")

	var req models.ConversationRequest
	contentType := c.Request().Header.Get("Content-Type")
	log.Printf("Content-Type: %s", contentType)

	if contentType == "application/json" {
		log.Println("Processing JSON payload")
		if err := c.Bind(&req); err != nil {
			c.Logger().Warnf("Bind error: %v", err)
			log.Printf("Bind error: %v\n", err)
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Invalid request payload.",
			})
		}
		log.Println("Successfully bound request to ConversationRequest model")
		c.Logger().Infof("Request bound: %+v", req)
	} else if contentType == "application/octet-stream" {
		log.Println("Processing raw PCM data")

		pcmData, err := io.ReadAll(c.Request().Body)
		if err != nil {
			c.Logger().Errorf("Error reading PCM data: %v", err)
			log.Printf("Error reading PCM data: %v\n", err)
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Failed to read PCM data.",
			})
		}
		log.Printf("Received PCM data of size: %d bytes", len(pcmData))

		log.Println("Logging all received headers:")
		for name, values := range c.Request().Header {
			for _, value := range values {
				log.Printf("Header: %s=%s", name, value)
			}
		}

		userIDStr := c.Request().Header.Get("X-User-ID")
		deviceIDStr := c.Request().Header.Get("X-Device-ID")
		language := c.Request().Header.Get("X-Language")

		if userIDStr == "" || deviceIDStr == "" {
			c.Logger().Warn("Missing required headers for raw PCM data")
			log.Println("Missing required headers: X-User-ID, X-Device-ID")
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Missing required headers: X-User-ID, X-Device-ID",
			})
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			c.Logger().Warnf("Invalid UserID format: %v", err)
			log.Printf("Invalid UserID format: %v\n", err)
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Invalid UserID format.",
			})
		}

		deviceID, err := strconv.Atoi(deviceIDStr)
		if err != nil {
			c.Logger().Warnf("Invalid DeviceID format: %v", err)
			log.Printf("Invalid DeviceID format: %v\n", err)
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Invalid DeviceID format.",
			})
		}

		req = models.ConversationRequest{
			UserID:      userID,
			DeviceID:    deviceID,
			RequestPCM:  pcmData,
			Language:    language,
		}

		log.Println("Constructed ConversationRequest from raw PCM data and headers")
		c.Logger().Infof("Request constructed: %+v", req)
	} else {
		c.Logger().Warnf("Unsupported Content-Type: %s", contentType)
		log.Printf("Unsupported Content-Type: %s\n", contentType)
		return c.JSON(http.StatusUnsupportedMediaType, map[string]string{
			"error": "Unsupported Media Type.",
		})
	}

	log.Printf("Validating RequestPCM length: %d", len(req.RequestPCM))
	if len(req.RequestPCM) < 16000 {
		c.Logger().Warn("RequestPCM too short")
		log.Println("RequestPCM is too short")
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "The request is too short.",
		})
	}
	log.Println("RequestPCM length is valid")


	var (
		previousConversation models.Conversation
		conversationData     models.ConversationData
		conversation         = new(models.Conversation)
		systemPrompt         string
	)
	conversation.UserID = req.UserID

	log.Println("Querying for previous conversation")
	query := `
        SELECT id, user_id, conversation_history, created_at
        FROM conversations
        WHERE user_id = $1
          AND created_at >= NOW() - $2 * INTERVAL '1 minute'
        ORDER BY created_at DESC
        LIMIT 1;
    `
	c.Logger().Info("Executing query to fetch previous conversation")
	log.Printf("Executing SQL Query with parameters: UserID=%s, ResetMinutes=%d", req.UserID, conversationResetMinutes)
	err := db.DB.QueryRow(query, req.UserID, conversationResetMinutes).Scan(
		&previousConversation.ID,
		&previousConversation.UserID,
		&previousConversation.ConversationHistory,
		&previousConversation.CreatedAt,
	)

	if err != nil && err != sql.ErrNoRows {
		c.Logger().Errorf("Error querying conversation: %v", err)
		log.Printf("Error querying conversation: %v\n", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to query conversation.",
		})
	}

	systemPrompt = systemprompt.DynamicGeneration(req.UserID)
	if err == sql.ErrNoRows {
		c.Logger().Info("No previous conversation found")
		log.Println("No previous conversation found within the reset time")
		
	} else {
		c.Logger().Infof("Found previous conversation ID: %d", previousConversation.ID)
		log.Printf("Previous conversation found: %+v\n", previousConversation)
		conversation.ID = previousConversation.ID
		conversation.CreatedAt = previousConversation.CreatedAt
		conversation.UserID = previousConversation.UserID

		if len(previousConversation.ConversationHistory) > 0 {
			err = json.Unmarshal(previousConversation.ConversationHistory, &conversationData)
			if err != nil {
				log.Printf("Error unmarshalling ConversationHistory: %v\n", err)
				c.Logger().Errorf("Error unmarshalling ConversationHistory: %v", err)
				return c.JSON(http.StatusInternalServerError, map[string]string{
					"error": "Failed to process conversation history.",
				})
			}
		} else {
			conversationData = models.ConversationData{}
		}
	}

	pcmData := req.RequestPCM
	log.Println("Received audio data")
	log.Printf("Audio data size: %d bytes", len(pcmData))

	log.Println("Converting PCM data to WAV format")
	wavData, err := pcm.ToWAV(pcmData)
	if err != nil {
		log.Printf("Failed to convert PCM to WAV: %v\n", err)
		c.Logger().Errorf("Failed to convert PCM to WAV: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to convert PCM to WAV.",
		})
	}
	log.Println("PCM data successfully converted to WAV format")

	log.Println("Generating transcription using Groq API")
	transcription, err := groq.GenerateWhisperTranscription(wavData, req.Language)
	if err != nil {
		log.Printf("Failed to get transcription: %v\n", err)
		c.Logger().Errorf("Failed to get transcription: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get transcription: " + err.Error(),
		})
	}
	log.Printf("Transcription received: %s\n", transcription)

	newMessage := models.Message{
		Sender:    "user",
		Content:   transcription,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
	conversationData.Messages = append(conversationData.Messages, newMessage)
	log.Printf("Appended user message to conversation history: %+v\n", newMessage)

	log.Println("Generating LLM response using Groq API")
	llmResponse, err := groq.GenerateLLMResponseFromConversationData(conversationData, systemPrompt, req.Language)
	if err != nil {
		c.Logger().Errorf("Error generating LLM response: %v", err)
		log.Printf("Error generating LLM response: %v\n", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to generate LLM response.",
		})
	}
	log.Printf("LLM response received: %+v\n", llmResponse)

	if len(llmResponse.Choices) == 0 {
		c.Logger().Warn("No choices returned from LLM response")
		log.Println("No choices returned from LLM response")
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "No choices returned from LLM response.",
		})
	}

	assistantResponse := llmResponse.Choices[0].Message.Content
	log.Printf("Assistant response extracted: %s\n", assistantResponse)

	assistantMessage := models.Message{
		Sender:    "assistant",
		Content:   assistantResponse,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
	conversationData.Messages = append(conversationData.Messages, assistantMessage)
	log.Printf("Appended assistant message to conversation history: %+v\n", assistantMessage)

	convoJSON, err := json.Marshal(conversationData)
	if err != nil {
		log.Printf("Error marshalling ConversationHistory: %v\n", err)
		c.Logger().Errorf("Error marshalling ConversationHistory: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to process conversation history.",
		})
	}

	if conversation.ID == 0 {
		log.Println("Inserting new conversation into the database")
		insertQuery := `
			INSERT INTO conversations (user_id, conversation_history)
			VALUES ($1, $2)
			RETURNING id, created_at;
		`

		log.Printf("Executing INSERT Query with UserID: %s and ConversationHistory: %s", conversation.UserID, string(convoJSON))

		err = db.DB.QueryRow(insertQuery, conversation.UserID, convoJSON).
			Scan(&conversation.ID, &conversation.CreatedAt)
		if err != nil {
			if pgErr, ok := err.(*pq.Error); ok && pgErr.Code.Name() == "foreign_key_violation" {
				c.Logger().Warn("Foreign key violation on inserting conversation")
				log.Println("Foreign key violation: Invalid user_id")
				return c.JSON(http.StatusBadRequest, map[string]string{
					"error": "Invalid user_id. User does not exist.",
				})
			}
			c.Logger().Errorf("Error inserting conversation: %v", err)
			log.Printf("Error inserting conversation: %v\n", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to store conversation.",
			})
		}
		log.Printf("New conversation inserted with ID: %d at %v\n", conversation.ID, conversation.CreatedAt)
	} else {
		log.Printf("Updating existing conversation ID: %d\n", conversation.ID)
		updateQuery := `
			UPDATE conversations
			SET conversation_history = $1, updated_at = NOW()
			WHERE id = $2
			RETURNING updated_at;
		`

		log.Printf("Executing UPDATE Query with ConversationHistory: %s and ID: %d", string(convoJSON), conversation.ID)

		err = db.DB.QueryRow(updateQuery, convoJSON, conversation.ID).
			Scan(&conversation.UpdatedAt)
		if err != nil {
			c.Logger().Errorf("Error updating conversation: %v", err)
			log.Printf("Error updating conversation: %v\n", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to update conversation.",
			})
		}
		log.Printf("Conversation ID: %d updated at %v\n", conversation.ID, conversation.UpdatedAt)
	}

	log.Printf("Current Conversation History: %s\n", string(convoJSON))
	log.Printf("Final assistant response to send: %s\n", assistantResponse)
	c.Logger().Info("Returning response to user")

	return c.JSON(http.StatusOK, map[string]string{
		"transcription": assistantResponse,
	})
}
