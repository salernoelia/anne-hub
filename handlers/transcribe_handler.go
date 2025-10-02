package handlers

import (
	"anne-hub/pkg/groq"
	"anne-hub/pkg/pcm"
	"io"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func TranscribeAudio(c echo.Context) error {
	log.Println("TranscribeAudio handler called")

	pcmData, err := io.ReadAll(c.Request().Body)
	if err != nil {
		log.Println("Failed to read request body:", err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Failed to read request body",
		})
	}
	defer c.Request().Body.Close()
	log.Println("Received audio data of size:", len(pcmData))

	wavData, err := pcm.ToWAV(pcmData)
	if err != nil {
		log.Println("Failed to convert PCM to WAV:", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to convert PCM to WAV",
		})
	}
	log.Println("PCM data converted to WAV format")

	transcription, err := groq.GenerateWhisperTranscription(wavData, "en")
	if err != nil {
		log.Println("Failed to get transcription:", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get transcription: " + err.Error(),
		})
	}
	log.Println("Received transcription from Groq API")

	log.Println("Transcription:", transcription)

	log.Println("Transcription sent to client")
	return c.JSON(http.StatusOK, map[string]string{
		"transcription": transcription,
	})
}
