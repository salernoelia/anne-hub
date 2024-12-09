// handlers/conversation_handler.go

package handlers

import (
	"anne-hub/models"
	"anne-hub/pkg/fs"
	"anne-hub/pkg/groq"
	"anne-hub/pkg/pcm"
	"anne-hub/pkg/systemprompt"
	"anne-hub/services"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

var conversationResetMinutes = 15

// ConversationHandler handles incoming conversation requests.
func ConversationHandler(c echo.Context) error {
	log.Println("Entered ConversationHandler")

	var req models.AnneWearConversationRequest
	contentType := c.Request().Header.Get("Content-Type")
	log.Printf("Content-Type: %s", contentType)

	// Parse request based on Content-Type
	if err := parseRequest(c, contentType, &req); err != nil {
		return err
	}

	// Validate RequestPCM
	if len(req.RequestPCM) < 16000 {
		log.Println("RequestPCM is too short")
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "The request is too short.",
		})
	}
	log.Println("RequestPCM length is valid")

	// Fetch previous conversation
	lastConversation, conversationHistory, err := services.GetPreviousConversation(req.UserID, conversationResetMinutes)
	if err != nil {
		log.Printf("Failed to query conversation: %v\n", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to query conversation.",
		})
	}

	systemPrompt := systemprompt.DynamicGeneration(req.UserID)

	// Handle audio conversion
	wavData, err := processPCMData(req.RequestPCM)
	if err != nil {
		return err
	}

	// Save WAV file
	if err := fs.WriteWAVDataToFile("recording.wav", wavData); err != nil {
		log.Println("Error saving WAV file")
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to save WAV file.",
		})
	}

	// Generate transcription
	transcription, err := groq.GenerateWhisperTranscription(wavData, req.Language)
	if err != nil {
		log.Printf("Failed to get transcription: %v\n", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get transcription: " + err.Error(),
		})
	}
	log.Printf("Transcription received: %s\n", transcription)

	// Append user message to conversation history
	services.AppendMessageToConversationHistory(&conversationHistory, "user", transcription)

	// Generate LLM response
	llmResponse, err := groq.GenerateLLMResponseFromConversationData(conversationHistory, systemPrompt, req.Language)
	if err != nil {
		log.Printf("Error generating LLM response: %v\n", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to generate LLM response.",
		})
	}

	if len(llmResponse.Choices) == 0 {
		log.Println("No choices returned from LLM response")
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "No choices returned from LLM response.",
		})
	}

	assistantResponse := llmResponse.Choices[0].Message.Content
	log.Printf("Assistant response extracted: %s\n", assistantResponse)

	// Append assistant message to conversation history
	services.AppendMessageToConversationHistory(&conversationHistory, "assistant", assistantResponse)

	// Marshal conversation history
	convoJSON, err := json.Marshal(conversationHistory)
	if err != nil {
		log.Printf("Error marshalling ConversationHistory: %v\n", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to process conversation history.",
		})
	}

	// Insert or update conversation in the database
	if lastConversation == nil {
		if err := services.InsertNewConversation(req.UserID, convoJSON); err != nil {
			return err
		}
	} else {
		if err := services.UpdateExistingConversation(lastConversation.ID, convoJSON); err != nil {
			return err
		}
	}

	log.Printf("Final assistant response to send: %s\n", assistantResponse)
	c.Logger().Info("Returning response to user")

	return c.JSON(http.StatusOK, map[string]string{
		"transcription": assistantResponse,
	})
}

// parseRequest parses the incoming request based on Content-Type.
func parseRequest(c echo.Context, contentType string, req *models.AnneWearConversationRequest) error {
	if contentType == "application/json" {
		log.Println("Processing JSON payload")
		if err := c.Bind(req); err != nil {
			log.Printf("Bind error: %v\n", err)
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Invalid request payload.",
			})
		}
		log.Println("Successfully bound request to AnneWearConversationRequest model")
	} else if contentType == "application/octet-stream" {
		log.Println("Processing raw PCM data")
		if err := parsePCMRequest(c, req); err != nil {
			return err
		}
	} else {
		log.Printf("Unsupported Content-Type: %s\n", contentType)
		return c.JSON(http.StatusUnsupportedMediaType, map[string]string{
			"error": "Unsupported Media Type.",
		})
	}
	return nil
}

// parsePCMRequest handles the parsing of PCM data requests.
func parsePCMRequest(c echo.Context, req *models.AnneWearConversationRequest) error {
	pcmData, err := io.ReadAll(c.Request().Body)
	if err != nil {
		log.Printf("Error reading PCM data: %v\n", err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Failed to read PCM data.",
		})
	}
	log.Printf("Received PCM data of size: %d bytes", len(pcmData))

	if len(pcmData) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Empty PCM data.",
		})
	}

	userIDStr := c.Request().Header.Get("X-User-ID")
	deviceIDStr := c.Request().Header.Get("X-Device-ID")
	language := c.Request().Header.Get("X-Language")

	if userIDStr == "" || deviceIDStr == "" || language == "" {
		log.Println("Missing required headers: X-User-ID, X-Device-ID, X-Language")
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Missing required headers: X-User-ID, X-Device-ID, X-Language",
		})
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		log.Printf("Invalid UserID format: %v\n", err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid UserID format.",
		})
	}

	deviceID, err := strconv.Atoi(deviceIDStr)
	if err != nil {
		log.Printf("Invalid DeviceID format: %v\n", err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid DeviceID format.",
		})
	}

	*req = models.AnneWearConversationRequest{
		UserID:     userID,
		DeviceID:   deviceID,
		RequestPCM: pcmData,
		Language:   language,
	}

	log.Println("Constructed AnneWearConversationRequest from raw PCM data and headers")
	return nil
}



// processPCMData converts PCM data to WAV format.
func processPCMData(pcmData []byte) ([]byte, error) {
	log.Println("Converting PCM data to WAV format")
	wavData, err := pcm.ToWAV(pcmData)
	if err != nil {
		log.Printf("Failed to convert PCM to WAV: %v\n", err)
		return nil, &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  "Failed to convert PCM to WAV.",
			Internal: err,
		}
	}
	log.Println("PCM data successfully converted to WAV format")
	return wavData, nil
}



