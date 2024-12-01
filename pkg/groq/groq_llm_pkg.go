// File: anne-hub/pkg/groq/groq_llm_pkg.go
package groq

import (
	"anne-hub/models"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// GenerateGroqLLMResponse generates a response from the Groq LLM API
func GenerateGroqLLMResponse(requestContent string, language string) (string, error) {
    apiKey := os.Getenv("GROQ_API_KEY")
    if apiKey == "" {
        return "", fmt.Errorf("GROQ_API_KEY environment variable is not set")
    }

    if language == "german" {
        requestContent += " Answer in German please"
    } else if language == "english" {
        requestContent += " Answer in English please"
    }

    url := "https://api.groq.com/openai/v1/chat/completions"
    payload := map[string]interface{}{
        "messages": []map[string]string{
            {"role": "user", "content": requestContent},
        },
        "model": "llama-3.1-70b-versatile",
    }
    jsonData, err := json.Marshal(payload)
    if err != nil {
        return "", fmt.Errorf("error encoding request content: %v", err)
    }

    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
    if err != nil {
        return "", fmt.Errorf("error creating request: %v", err)
    }

    req.Header.Set("Authorization", "Bearer "+apiKey)
    req.Header.Set("Content-Type", "application/json")

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


    var apiResp models.GroqLLMResponse
    if err := json.Unmarshal(body, &apiResp); err != nil {
        return "", fmt.Errorf("error decoding response from Groq API: %v", err)
    }

    if len(apiResp.Choices) > 0 && len(apiResp.Choices[0].Message.Content) > 0 {
        return apiResp.Choices[0].Message.Content, nil
    }
    return "No sentence generated.", nil
}
