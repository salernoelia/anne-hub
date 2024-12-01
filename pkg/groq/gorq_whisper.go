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

/*
curl https://api.groq.com/openai/v1/audio/transcriptions \
  -H "Authorization: bearer ${GROQ_API_KEY}" \
  -F "file=@./sample_audio.m4a" \
  -F model=whisper-large-v3-turbo \
  -F temperature=0 \
  -F response_format=json \
  -F language=en
*/
func GenerateGroqWhisperTranscription(wavData []byte, language string) (string, error) {
    apiKey := os.Getenv("GROQ_API_KEY")
    if apiKey == "" {
        return "", fmt.Errorf("GROQ_API_KEY environment variable is not set")
    }

    url := "https://api.groq.com/openai/v1/audio/transcriptions"

    var b bytes.Buffer
    w := multipart.NewWriter(&b)

    // Add the audio file to the request
    fw, err := w.CreateFormFile("file", "audio.wav")
    if err != nil {
        return "", err
    }
    _, err = fw.Write(wavData)
    if err != nil {
        return "", err
    }

    // Add other form fields
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

    // Parse the response
    var apiResp models.GroqWhisperResponse
    if err := json.Unmarshal(body, &apiResp); err != nil {
        return "", fmt.Errorf("error decoding response from Groq API: %v", err)
    }

    return apiResp.Text, nil
}
