// handlers/websocket_handler.go
package handlers

import (
	"anne-hub/models"
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type WavHeader struct {
    RIFF      [4]byte
    FileSize  uint32
    WAVE      [4]byte
}

type FormatChunk struct {
    FmtID          [4]byte
    FmtSize        uint32
    AudioFormat    uint16
    NumChannels    uint16
    SampleRate     uint32
    ByteRate       uint32
    BlockAlign     uint16
    BitsPerSample  uint16
    // If fmtSize > 16, there may be extra bytes here...
}

type DataChunkHeader struct {
    DataID   [4]byte
    DataSize uint32
}

// WebSocketConversationHandler handles WebSocket connections for PCM data collection.
func WebSocketTestHandler(c echo.Context) error {
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
              
            //    ttsByteArray, err := os.ReadFile("test-8bit-6000.wav")
			// 	if err != nil {
			// 		log.Printf("Error reading WAV file: %v", err)
			// 		conn.WriteMessage(websocket.TextMessage, []byte("Error reading WAV file"))
			// 		break
			// 	}

			// 	// Typically, for a standard WAV file, the first 44 bytes are the header.
			// 	// You should confirm this by inspecting the WAV header if you're unsure.
			// 	audioStart := 44 // standard for PCM WAV
			// 	if len(ttsByteArray) <= audioStart {
			// 		log.Println("WAV file is too small or corrupted.")
			// 		break
			// 	}

				// rawAudioData := ttsByteArray[audioStart:]



                fmt.Println("Sending WAV file link for Streaming to client")
                err = conn.WriteMessage(websocket.TextMessage, []byte("/files/test_linear16.wav"))
                if err != nil {
                    log.Printf("Error sending PCM data: %v\n", err)
                }

            //    chunkSize := 400 
			// 	for start := 0; start < len(rawAudioData); start += chunkSize {
			// 		end := start + chunkSize
			// 		if end > len(rawAudioData) {
			// 			end = len(rawAudioData)
			// 		}
			// 		chunk := rawAudioData[start:end]

			// 		err := conn.WriteMessage(websocket.BinaryMessage, chunk)
			// 		if err != nil {
			// 			log.Printf("Error sending chunk: %v\n", err)
			// 			break
			// 		}
			// 		time.Sleep(10 * time.Millisecond) // optional throttle
			// 	}
			// 	log.Println("All chunks sent.")






                    // outFile, err := os.Create("test_linear16.wav")
                // if err != nil {
                //     return fmt.Errorf("failed to create output file: %w", err)
                // }
                // defer outFile.Close()
                // _, err = outFile.Write(ttsByteArray)
                // if err != nil {
                //     return fmt.Errorf("failed to write audio content to file: %w", err)
                // }
                
           

                // conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("WAV file saved")))

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