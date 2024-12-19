// handlers/websocket_handler.go
package handlers

import (
	"anne-hub/models"
	"anne-hub/pkg/groq"
	"anne-hub/pkg/pcm"
	"anne-hub/pkg/systemprompt"
	"anne-hub/pkg/tts"
	"anne-hub/services"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var lastEmotionSent string = "suspicious"
var emotionChanged bool = false

var allowedEmotions = map[string]bool{
	"celebration": true,
	"suspicious":  true,
	"cute_smile":  true,
	"curiosity":   true,
	"confused":    true,
	"sleep":       true,
	"lucky_smile": true,
	"surprised":   true,
}

var wsUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func playAudio(filePath string) error {
	cmd := exec.Command("afplay", filePath)
	return cmd.Run()
}

type TaskCompletion struct {
	Task      string `json:"task,omitempty"`
	Completed string `json:"completed,omitempty"`
}

type LLMResponseJSONfromPrompt struct {
	Message        string         `json:"message"`
	Emotion        string         `json:"emotion"`
	TaskCompletion TaskCompletion `json:"task_completion"`
}

var assistantResponseJSON string

func WebSocketConversationHandler(c echo.Context) error {
	conn, err := wsUpgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return err
	}
	defer conn.Close()

	var headers models.WSRequestHeaders
	var pcmData []byte
	headersReceived := false

	for {
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

			if msg == "PING" {
				conn.WriteMessage(websocket.TextMessage, []byte("PONG"))
				break
			}

			if !headersReceived {
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
				emotionChanged = false
				continue
			}

			if msg == "EOS" {
				currentConversation, err := services.HandleProcessConversationInput(pcmData, headers)
				if err != nil {
					log.Printf("Error processing conversation: %v", err)
					conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Processing error: %s", err.Error())))
					break
				}

				wavData, err := processPCMData(currentConversation.RequestPCM)
				if err != nil {
					return err
				}

				outFile, err := os.Create("m5audio.wav")
				if err != nil {
					return fmt.Errorf("failed to create output file: %w", err)
				}

				_, err = outFile.Write(wavData)
				if err != nil {
					return fmt.Errorf("failed to write audio content to file: %w", err)
				}

				defer outFile.Close()

				transcription, err := groq.GenerateWhisperTranscription(wavData, currentConversation.Language)
				if err != nil {
					log.Printf("Failed to get transcription: %v\n", err)
					return c.JSON(http.StatusInternalServerError, map[string]string{
						"error": "Failed to get transcription: " + err.Error(),
					})
				}

				transcription += "<for assistant: you must return as json as instructed in system prompt format: {\"message\": \"<your message>\", \"emotion\": \"<emotion>\", \"task_completion\": {\"task:\": \"<task_id>\", \"completed\": \"<value>\"}}"
				transcription += ", if there was no task mentioned, add an empty task_completion object>"

                log.Print("/----------------------------------------------------------------/")
				log.Printf("Transcription received: %s\n", transcription)
                log.Print("/----------------------------------------------------------------/")

				lastConversation, conversationHistory, err := services.GetPreviousConversation(currentConversation.UserID, 15)
				if err != nil {
					log.Printf("Failed to query conversation: %v\n", err)
					break
				}

				// log.Printf("Last conversation: %v\n", lastConversation)
				// log.Printf("Conversation history: %v\n", conversationHistory)

				systemPrompt := systemprompt.DynamicGeneration(currentConversation.UserID)
				services.AppendMessageToConversationHistory(&conversationHistory, "user", transcription)

				llmResponse, err := groq.GenerateLLMResponseFromConversationData(conversationHistory, systemPrompt, currentConversation.Language)
				if err != nil {
					log.Printf("Error generating LLM response: %v", err)
					break
				}

				if len(llmResponse.Choices) == 0 {
					log.Println("No choices returned from LLM response")
					break
				}

				DirtyAssistantResponseJSON := llmResponse.Choices[0].Message.Content
				log.Printf("/----------------------------------------------------------------/\n")
				log.Printf("Assistant response JSON: %s\n", DirtyAssistantResponseJSON)
				log.Printf("/----------------------------------------------------------------/\n")

				if !strings.Contains(DirtyAssistantResponseJSON, "{") || !strings.Contains(DirtyAssistantResponseJSON, "}") {
					log.Printf("\033[31mInvalid JSON response received: %s\033[0m\n", DirtyAssistantResponseJSON)
					assistantResponseJSON := `{
						"message": "I didn't quite understand that. Could you please try again?",
						"emotion": "confused",
						"task_completion": {}
					}`
					handleDefaultResponse(&conversationHistory, assistantResponseJSON, currentConversation, lastConversation)
					break
				}

				assistantResponseJSON = DirtyAssistantResponseJSON[strings.Index(DirtyAssistantResponseJSON, "{"):strings.LastIndex(DirtyAssistantResponseJSON, "}")+1]

				var assistantResponse LLMResponseJSONfromPrompt
				err = json.Unmarshal([]byte(assistantResponseJSON), &assistantResponse)
				if err != nil {
					log.Printf("\033[31mError unmarshalling assistant response JSON: %v\033[0m\n", err)
					break
				}

				if strings.TrimSpace(assistantResponse.Message) == "" {
					log.Println("\033[31mAssistant response message is empty\033[0m")
					break
				}

				if strings.TrimSpace(assistantResponse.Emotion) == "" {
					log.Println("\033[31mAssistant response emotion is empty\033[0m")
					assistantResponse.Emotion = "cute_smile"
				}

				// validTaskIDs := extractValidTaskIDs(conversationHistory)

				if !isValidFormat(assistantResponse) {
					log.Printf("\033[31mInvalid JSON response format received: %s\033[0m\n", assistantResponseJSON)
					assistantResponseJSON := `{
						"message": "I didn't quite understand that. Could you please try again?",
						"emotion": "confused",
						"task_completion": {}
					}`
					handleDefaultResponse(&conversationHistory, assistantResponseJSON, currentConversation, lastConversation)
					break
				}

				services.AppendMessageToConversationHistory(&conversationHistory, "assistant", assistantResponse.Message)

				convoJSON, err := json.Marshal(conversationHistory)
				if err != nil {
					log.Printf("\033[31mError marshalling ConversationHistory: %v\033[0m\n", err)
					break
				}

				if lastConversation == nil {
					if err := services.InsertNewConversation(currentConversation.UserID, convoJSON); err != nil {
						log.Print("Failed inserting new Conversation")
						break
					}
				} else {
					if err := services.UpdateExistingConversation(lastConversation.ID, convoJSON); err != nil {
						log.Print("Failed updating Conversation")
						break
					}
				}
				log.Printf("/----------------------------------------------------------------/\n")
				log.Printf("Final assistant response to send: %s\n", assistantResponse.Message)
				log.Printf("/----------------------------------------------------------------/\n")

				conn.WriteMessage(websocket.TextMessage, []byte(assistantResponse.Emotion))

                tts, err := tts.ElevenLabsTextToSpeech(assistantResponse.Message)
                if err != nil {
                    log.Print("Error converting text to speech:", err)
                    break;
                }


              
                audioDir := "./audio"
                err = os.MkdirAll(audioDir, os.ModePerm)
                if err != nil {
                    log.Fatalf("Failed to create audio directory: %v", err)
                }

                rand.Seed(time.Now().UnixNano())
                randomNumber := rand.Intn(1000000) // Random number up to 6 digits
                fileName := fmt.Sprintf("audio_%d.wav", randomNumber) // Assuming the output format is .wav
                filePath := filepath.Join(audioDir, fileName)


                err = pcm.TTStoWav(tts, filePath)
                if err != nil {
                    log.Printf("Failed to convert TTS to WAV: %v", err)
                    break
                }

                fmt.Printf("Audio file saved successfully at %s\n", filePath)

                err = playAudio(filePath)
                if err != nil {
                    log.Fatalf("Failed to play audio: %v", err)
                }

				emotionChanged = false

				assistantResponse = LLMResponseJSONfromPrompt{}

				pcmData = nil
			}

		case websocket.BinaryMessage:
			if !headersReceived {
				log.Println("Received binary data before headers. Ignoring.")
				conn.WriteMessage(websocket.TextMessage, []byte("Headers must be sent before PCM data."))
				continue
			}

			log.Printf("Received %d bytes of PCM data", len(message))
			pcmData = append(pcmData, message...)

		default:
			log.Printf("Unsupported message type: %d", messageType)
			conn.WriteMessage(websocket.TextMessage, []byte("Unsupported message type."))
		}
	}

	return nil
}

func handleDefaultResponse(conversationHistory *models.ConversationHistory, defaultJSON string, currentConversation models.AnneWearConversationRequest, lastConversation *models.Conversation) {
	var defaultResponse LLMResponseJSONfromPrompt
	err := json.Unmarshal([]byte(defaultJSON), &defaultResponse)
	if err != nil {
		log.Printf("\033[31mError unmarshalling default response JSON: %v\033[0m\n", err)
		return
	}

	services.AppendMessageToConversationHistory(conversationHistory, "assistant", defaultResponse.Message)

	convoJSON, err := json.Marshal(conversationHistory)
	if err != nil {
		log.Printf("\033[31mError marshalling ConversationHistory: %v\033[0m\n", err)
		return
	}

	if lastConversation == nil {
		if err := services.InsertNewConversation(currentConversation.UserID, convoJSON); err != nil {
			log.Printf("\033[31mFailed inserting new Conversation: %v\033[0m\n", err)
			return
		}
	} else {
		if err := services.UpdateExistingConversation(lastConversation.ID, convoJSON); err != nil {
			log.Printf("\033[31mFailed updating Conversation: %v\033[0m\n", err)
			return
		}
	}

	log.Printf("/----------------------------------------------------------------/\n")
	log.Printf("Final assistant response to send: %s\n", defaultResponse.Message)
	log.Printf("/----------------------------------------------------------------/\n")
}

func isValidFormat(response LLMResponseJSONfromPrompt) bool {
	if _, exists := allowedEmotions[response.Emotion]; !exists {
		log.Printf("\033[31mEmotion '%s' is not allowed.\033[0m\n", response.Emotion)
		return false
	}

	if response.TaskCompletion.Task == "" && response.TaskCompletion.Completed == "" {
		return true
	}

	if response.TaskCompletion.Task != "" || response.TaskCompletion.Completed != "" {
		if response.TaskCompletion.Task == "" || response.TaskCompletion.Completed == "" {
			log.Printf("\033[31mBoth 'task' and 'completed' fields must be present if one is present.\033[0m\n")
			return false
		}
		completedLower := strings.ToLower(response.TaskCompletion.Completed)
		if completedLower != "true" && completedLower != "false" {
			log.Printf("\033[31m'completed' field must be either 'true' or 'false'. Received: '%s'\033[0m\n", response.TaskCompletion.Completed)
			return false
		}

		return true
	}

	log.Printf("\033[31mInvalid task_completion format.\033[0m\n")
	return false
}

// func isValidTaskID(taskID string, validTaskIDs []string) bool {
// 	for _, id := range validTaskIDs {
// 		if id == taskID {
// 			return true
// 		}
// 	}
// 	return false
// }

// func extractValidTaskIDs(conversationHistory models.ConversationHistory) []string {
// 	var taskIDs []string
// 	for _, task := range conversationHistory.Tasks {
// 		taskIDs = append(taskIDs, task.ID)
// 	}
// 	return taskIDs
// }
