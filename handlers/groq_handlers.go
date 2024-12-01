// anne-hub/handlers/groq.go

package handlers

import (
	"anne-hub/models"
	"anne-hub/pkg/db"
	"anne-hub/pkg/groq"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

// CreateConversationHandler handles the creation of a new conversation
func CreateConversationHandler(c echo.Context) error {
    // Bind the incoming JSON to the request struct
    req := new(models.CreateConversationRequest)
    if err := c.Bind(req); err != nil {
        c.Logger().Warnf("Bind error: %v", err)
        return c.JSON(http.StatusBadRequest, map[string]string{
            "error": "Invalid request payload.",
        })
    }

    // Validate UserID format
    if _, err := uuid.Parse(req.UserID); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{
            "error": "Invalid user_id format.",
        })
    }

    // Ensure the Request field is not empty
    if req.Request == "" {
        return c.JSON(http.StatusBadRequest, map[string]string{
            "error": "Request field is required.",
        })
    }


    // Generate LLM response using the groq package
    response, err := groq.GenerateGroqLLMResponse(req.Request, req.Language)
    if err != nil {
        c.Logger().Errorf("Error generating LLM response: %v", err)
        return c.JSON(http.StatusInternalServerError, map[string]string{
            "error": "Failed to generate LLM response.",
        })
    }


    // Create a new Conversation instance
    conv := models.Conversation{
        UserID:    req.UserID,
        Request:   req.Request,
        Response:  response,
        ModelUsed: "llama-3.1-70b-versatile",
        Role:      "user",
    }

    // Insert the conversation into the database
    insertQuery := `
        INSERT INTO conversations (user_id, request, response, model_used, role)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id, conversation_id, created_at
    `
    err = db.DB.QueryRow(insertQuery, conv.UserID, conv.Request, conv.Response, conv.ModelUsed, conv.Role).Scan(&conv.ID, &conv.ConversationID, &conv.CreatedAt)
    if err != nil {
        // Handle foreign key constraint (e.g., invalid user_id)
        if pgErr, ok := err.(*pq.Error); ok {
            switch pgErr.Code.Name() {
            case "foreign_key_violation":
                return c.JSON(http.StatusBadRequest, map[string]string{
                    "error": "Invalid user_id. User does not exist.",
                })
            }
        }

        c.Logger().Errorf("Error inserting conversation: %v", err)
        return c.JSON(http.StatusInternalServerError, map[string]string{
            "error": "Failed to store conversation.",
        })
    }

	// for now the role is hardcoded to assistant
    conv.Role = "assistant"

    return c.JSON(http.StatusCreated, conv)
}
