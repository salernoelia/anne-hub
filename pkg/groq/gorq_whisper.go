package groq

import (
	"anne-hub/models"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

func GenerateWhisperTranscription(wavData []byte, language string) (string, error) {
	apiKey := os.Getenv("GROQ_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("GROQ_API_KEY environment variable is not set")
	}

	url := "https://api.groq.com/openai/v1/audio/transcriptions"

	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	fw, err := w.CreateFormFile("file", "audio.wav")
	if err != nil {
		return "", err
	}
	_, err = fw.Write(wavData)
	if err != nil {
		return "", err
	}

	w.WriteField("model", "whisper-large-v3-turbo")
	w.WriteField("temperature", "0")
	w.WriteField("response_format", "json")
	w.WriteField("language", language)

	w.Close()

	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", w.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request to Groq API: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	log.Println("Groq API response:", string(body))

	var apiResp models.GroqWhisperResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return "", fmt.Errorf("error decoding response from Groq API: %v", err)
	}

	return apiResp.Text, nil
}
