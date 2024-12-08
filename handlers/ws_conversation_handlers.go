// handlers/websocket_handler.go
package handlers

import (
	"anne-hub/pkg/pcm"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

// Define a struct to represent the custom headers
type CustomHeaders struct {
    XUserID    string `json:"X-User-ID"`
    XDeviceID  string `json:"X-Device-ID"`
    XLanguage  string `json:"X-Language"`
}

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
    var headers CustomHeaders
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
                if len(pcmData) == 0 {
                    log.Println("No PCM data received before EOS.")
                    conn.WriteMessage(websocket.TextMessage, []byte("No PCM data received."))
                    continue
                }

                log.Printf("Total PCM data size: %d bytes", len(pcmData))

                // Convert the accumulated PCM data to WAV format
                wavBytes, err := pcm.ToWAV(pcmData)
                if err != nil {
                    log.Printf("Error converting PCM to WAV: %v", err)
                    conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Conversion error: %v", err)))
                    break
                }

                // Generate a unique filename using the current timestamp
                filename := fmt.Sprintf("recording_%d", time.Now().Unix())
                err = saveWAVFile(filename+".wav", wavBytes)
                if err != nil {
                    log.Printf("Error saving WAV file: %v", err)
                    conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("File saving error: %v", err)))
                    break
                }

                err = saveRawPCMFile(filename+".pcm", pcmData)
                if err != nil {
                    log.Printf("Error saving raw PCM file: %v", err)
                }

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

                log.Printf("WAV file saved successfully as %s", filename)
                conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("WAV file saved as %s", filename)))

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