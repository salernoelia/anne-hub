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
	RIFF     [4]byte
	FileSize uint32
	WAVE     [4]byte
}

type FormatChunk struct {
	FmtID         [4]byte
	FmtSize       uint32
	AudioFormat   uint16
	NumChannels   uint16
	SampleRate    uint32
	ByteRate      uint32
	BlockAlign    uint16
	BitsPerSample uint16
}

type DataChunkHeader struct {
	DataID   [4]byte
	DataSize uint32
}

func WebSocketTestHandler(c echo.Context) error {
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
				continue
			}

			if msg == "EOS" { // Define "EOS" as the end-of-stream message

				fmt.Println("Sending WAV file link for Streaming to client")
				err = conn.WriteMessage(websocket.TextMessage, []byte("/files/test_linear16.wav"))
				if err != nil {
					log.Printf("Error sending PCM data: %v\n", err)
				}

				pcmData = nil
			} else {

				log.Printf("Unknown text message received: %s", msg)
				conn.WriteMessage(websocket.TextMessage, []byte("Unknown command."))
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
