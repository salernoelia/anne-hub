// handlers/websocket_handler.go
package handlers

import (
	"anne-hub/models"
	"anne-hub/pkg/groq"
	"anne-hub/pkg/systemprompt"
	"anne-hub/services"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

// Upgrader configures the WebSocket upgrade parameters.
var wsUpgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool { return true }, // Allow all origins; adjust as needed for security
}


// WebSocketConversationHandler handles WebSocket connections for PCM data collection.
func WebSocketConversationHandler(c echo.Context) error {
    // Upgrade the HTTP connection to a WebSocket connection
    conn, err := wsUpgrader.Upgrade(c.Response(), c.Request(), nil)
    if err != nil {
        log.Println("WebSocket upgrade error:", err)
        return err
    }
    defer conn.Close()

    // Initialize variables to store headers and PCM data
    var headers models.WSRequestHeaders
    var pcmData []byte // Buffer to accumulate PCM data

    // Set a flag to check if headers have been received
    headersReceived := false

    for {
        // Read incoming messages
        messageType, message, err := conn.ReadMessage()
        if err != nil {
            if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
                log.Printf("Unexpected WebSocket error: %v", err)
            } else {
                log.Printf("WebSocket closed: %v", err)
            }
            break
        }

        switch messageType {
        case websocket.TextMessage:
            msg := string(message)
            log.Printf("Received text message: %s", msg)

            if !headersReceived {
                // Attempt to parse the headers JSON
                err := json.Unmarshal(message, &headers)
                if err != nil {
                    log.Printf("Error parsing headers JSON: %v", err)
                    conn.WriteMessage(websocket.TextMessage, []byte("Invalid headers format."))
                    continue
                }

                log.Printf("Custom Headers Received - User ID: %s, Device ID: %s, Language: %s",
                    headers.XUserID, headers.XDeviceID, headers.XLanguage)

                headersReceived = true
                conn.WriteMessage(websocket.TextMessage, []byte("Headers received successfully."))
                continue
            }

            // Handle control messages like "EOS"
            if msg == "EOS" { // Define "EOS" as the end-of-stream message
              
                req, err := services.HandleProcessConversationInput(pcmData, headers)
                if err != nil {
                    log.Printf("Error processing conversation: %v", err)
                    conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Processing error: %s", err.Error())))
                    break
                }


                // Handle audio conversion
                wavData, err := processPCMData(req.RequestPCM)
                if err != nil {
                    return err
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

                // Fetch previous conversation  
                lastConversation, conversationHistory, err := services.GetPreviousConversation(req.UserID, 15)
                if err != nil {
                    log.Printf("Failed to query conversation: %v\n", err)
                    break;
                }

                log.Printf("Last conversation: %v\n", lastConversation)
                log.Printf("Conversation history: %v\n", conversationHistory)

                systemPrompt := systemprompt.DynamicGeneration(req.UserID)

                // Append user message to conversation history
	            services.AppendMessageToConversationHistory(&conversationHistory, "user", transcription)

                // Generate LLM response
                llmResponse, err := groq.GenerateLLMResponseFromConversationData(conversationHistory, systemPrompt, req.Language)
                if err != nil {
                    log.Printf("Error generating LLM response: %v\n", err)
                    break;
                }

                if len(llmResponse.Choices) == 0 {
                    log.Println("No choices returned from LLM response")
                    break;
                }

                assistantResponse := llmResponse.Choices[0].Message.Content
                log.Printf("Assistant response extracted: %s\n", assistantResponse)

                // Append assistant message to conversation history
                services.AppendMessageToConversationHistory(&conversationHistory, "assistant", assistantResponse)

                // Marshal conversation history
                convoJSON, err := json.Marshal(conversationHistory)
                if err != nil {
                    log.Printf("Error marshalling ConversationHistory: %v\n", err)
                    break;
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




                fmt.Println("Sending WAV file to client")
                //send ws message to client with pcm data
                // err = conn.WriteMessage(websocket.BinaryMessage, pcmData)
                // if err != nil {
                //     log.Printf("Error sending PCM data: %v\n", err)
                // }
                err = conn.WriteMessage(websocket.TextMessage, []byte("Sending WAV file to client"))
                if err != nil {
                    log.Printf("Error sending PCM data: %v\n", err)
                }



                // conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("WAV file saved")))

                // Optionally, reset the PCM data buffer to allow for new recordings
                pcmData = nil
            } else {
                // Handle other text messages if necessary
                log.Printf("Unknown text message received: %s", msg)
                conn.WriteMessage(websocket.TextMessage, []byte("Unknown command."))
            }

        case websocket.BinaryMessage:
            if !headersReceived {
                log.Println("Received binary data before headers. Ignoring.")
                conn.WriteMessage(websocket.TextMessage, []byte("Headers must be sent before PCM data."))
                continue
            }

            // Append binary PCM data to the buffer
            log.Printf("Received %d bytes of PCM data", len(message))
            pcmData = append(pcmData, message...)

        default:
            log.Printf("Unsupported message type: %d", messageType)
            conn.WriteMessage(websocket.TextMessage, []byte("Unsupported message type."))
        }
    }

    return nil
}