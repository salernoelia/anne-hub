package groq

import (
	"anne-hub/models"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
)

/*
Currently needs chat imlementation trough storing some kind of chat id and state, which gets passed to this function.
*/

func GenerateGroqLLMResponse(data string, language string) string {
    apiKey := os.Getenv("GROQ_API_KEY")
    if apiKey == "" {
        log.Fatal("GROQ_API_KEY environment variable is not set.")
    }

    if language == "german" {
        data += " Answer in German please"
    } else if language == "english" {
        data += " Answer in English please"
    }

    url := "https://api.groq.com/openai/v1/chat/completions"
    payload := map[string]interface{}{
        "messages": []map[string]string{
            {"role": "user", "content": data},
        },
        "model": "llama-3.1-70b-versatile", 
    }
    jsonData, err := json.Marshal(payload)
    if err != nil {
        log.Fatalf("Error encoding request data: %v", err)
    }

    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
    if err != nil {
        log.Fatalf("Error creating request: %v", err)
    }

    req.Header.Set("Authorization", "Bearer "+apiKey)
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        log.Fatalf("Error sending request to Groq API: %v", err)
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        log.Fatalf("Error reading response body: %v", err)
    }

    var apiResp models.APIResponse
    if err := json.Unmarshal(body, &apiResp); err != nil {
        log.Fatalf("Error decoding response from Groq API: %v", err)
    }

    if len(apiResp.Choices) > 0 && len(apiResp.Choices[0].Message.Content) > 0 {
        return apiResp.Choices[0].Message.Content
    }
    return "No sentence generated."
}
