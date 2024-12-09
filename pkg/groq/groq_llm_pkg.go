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
func GenerateLLMResponse(userPrompt string, systemPrompt string, language string) (models.GroqLLMResponse, error) {
	apiKey := os.Getenv("GROQ_API_KEY")
	if apiKey == "" {
		return models.GroqLLMResponse{}, fmt.Errorf("GROQ_API_KEY environment variable not set")
	}

	if language == "german" {
		userPrompt += " Answer in German please"
	} else if language == "english" {
		userPrompt += " Answer in English please"
	}

	url := "https://api.groq.com/openai/v1/chat/completions"
	payload := map[string]interface{}{
		"messages": []map[string]string{
			{"role": "system", "content": systemPrompt},
			{"role": "user", "content": userPrompt},
		},
		"model": "llama-3.1-70b-versatile",
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return models.GroqLLMResponse{}, fmt.Errorf("error encoding request content: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return models.GroqLLMResponse{}, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return models.GroqLLMResponse{}, fmt.Errorf("error sending request to Groq API: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.GroqLLMResponse{}, fmt.Errorf("error reading response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return models.GroqLLMResponse{}, fmt.Errorf("groq API returned status %d: %s", resp.StatusCode, string(body))
	}

	var apiResp models.GroqLLMResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return models.GroqLLMResponse{}, fmt.Errorf("error decoding response from Groq API: %w", err)
	}

	if len(apiResp.Choices) == 0 || len(apiResp.Choices[0].Message.Content) == 0 {
		return models.GroqLLMResponse{}, fmt.Errorf("no valid response received from Groq API")
	}

	return apiResp, nil
}


// GenerateGroqLLMResponse generates a response from the Groq LLM API using structured conversation data

// GenerateGroqLLMResponseFromConversationData generates a response from the Groq LLM API using structured conversation data
func GenerateLLMResponseFromConversationData(conversation models.ConversationHistory, systemPrompt string, language string) (models.GroqLLMResponse, error) {
	apiKey := os.Getenv("GROQ_API_KEY")
	if apiKey == "" {
		return models.GroqLLMResponse{}, fmt.Errorf("GROQ_API_KEY environment variable not set")
	}

	if language == "german" {
		systemPrompt += " Bitte antworte auf Deutsch."
	} else if language == "english" {
		systemPrompt += " Please respond in English."
	}

	url := "https://api.groq.com/openai/v1/chat/completions"

	messages := []map[string]string{
		{"role": "system", "content": systemPrompt},
	}

	for _, msg := range conversation.Messages {
		role := "user"
		if msg.Sender == "assistant" {
			role = "assistant"
		}
		messages = append(messages, map[string]string{
			"role":    role,
			"content": msg.Content,
		})
	}

	payload := map[string]interface{}{
		"messages": messages,
		"model":    "llama-3.1-70b-versatile",
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return models.GroqLLMResponse{}, fmt.Errorf("error encoding request content: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return models.GroqLLMResponse{}, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return models.GroqLLMResponse{}, fmt.Errorf("error sending request to Groq API: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.GroqLLMResponse{}, fmt.Errorf("error reading response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return models.GroqLLMResponse{}, fmt.Errorf("groq API returned status %d: %s", resp.StatusCode, string(body))
	}

	var apiResp models.GroqLLMResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return models.GroqLLMResponse{}, fmt.Errorf("error decoding response from Groq API: %w", err)
	}

	if len(apiResp.Choices) == 0 || apiResp.Choices[0].Message.Content == "" {
		return models.GroqLLMResponse{}, fmt.Errorf("no valid response received from Groq API")
	}

	return apiResp, nil
}