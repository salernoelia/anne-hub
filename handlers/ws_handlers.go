// handlers/websocket_handler.go
package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"anne-hub/pkg/pcm" // Replace with your actual module path

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

// Upgrader configures the WebSocket upgrade parameters.
var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool { return true }, // Allow all origins; adjust as needed for security
}

// WebSocketHandler handles WebSocket connections for PCM data collection.
func WebSocketHandler(c echo.Context) error {
    // Upgrade the HTTP connection to a WebSocket connection
    conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
    if err != nil {
        log.Println("WebSocket upgrade error:", err)
        return err
    }
    defer conn.Close()

    var pcmData []byte // Buffer to accumulate PCM data

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
        case websocket.BinaryMessage:
            // Append binary PCM data to the buffer
            log.Printf("Received %d bytes of PCM data", len(message))
            pcmData = append(pcmData, message...)

        case websocket.TextMessage:
            msg := string(message)
            log.Printf("Received text message: %s", msg)

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
                err = saveWAVFile(filename + ".wav", wavBytes)
                if err != nil {
                    log.Printf("Error saving WAV file: %v", err)
                    conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("File saving error: %v", err)))
                    break
                }

				err = saveRawPCMFile(filename + ".pcm", pcmData)
				if err != nil {
					log.Printf("Error saving raw PCM file: %v", err)
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

        default:
            log.Printf("Unsupported message type: %d", messageType)
            conn.WriteMessage(websocket.TextMessage, []byte("Unsupported message type."))
        }
    }

    return nil
}

// saveWAVFile saves the WAV data to the filesystem with the given filename.
func saveWAVFile(filename string, data []byte) error {
    // Create or truncate the file
    file, err := os.Create(filename)
    if err != nil {
        return fmt.Errorf("failed to create file %s: %w", filename, err)
    }
    defer file.Close()

    // Write WAV data to the file
    _, err = file.Write(data)
    if err != nil {
        return fmt.Errorf("failed to write data to file %s: %w", filename, err)
    }

    return nil
}

func saveRawPCMFile(filename string, data []byte) error {
    file, err := os.Create(filename)
    if err != nil {
        return fmt.Errorf("failed to create file %s: %w", filename, err)
    }
    defer file.Close()

    _, err = file.Write(data)
    if err != nil {
        return fmt.Errorf("failed to write PCM data to file %s: %w", filename, err)
    }

    return nil
}