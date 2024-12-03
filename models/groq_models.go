package models

type GroqLLMResponse struct {
    ID     string `json:"id"`
    Object string `json:"object"`
    Created int `json:"created"`
    Model  string `json:"model"`
    Choices []struct {
        Index   int `json:"index"`
        Message struct {
            Role    string `json:"role"`
            Content string `json:"content"`
        } `json:"message"`
        Logprobs     interface{} `json:"logprobs"`
        FinishReason string      `json:"finish_reason"`
    } `json:"choices"`
    Usage struct {
        QueueTime       float64 `json:"queue_time"`
        PromptTokens    int     `json:"prompt_tokens"`
        PromptTime      float64 `json:"prompt_time"`
        CompletionTokens int    `json:"completion_tokens"`
        CompletionTime  float64 `json:"completion_time"`
        TotalTokens     int     `json:"total_tokens"`
        TotalTime       float64 `json:"total_time"`
    } `json:"usage"`
    SystemFingerprint string `json:"system_fingerprint"`
    XGroq             struct {
        ID string `json:"id"`
    } `json:"x_groq"`
}

type GroqWhisperResponse struct {
    ID               int     `json:"id"`
    Seek             int     `json:"seek"`
    Start            float64 `json:"start"`
    End              float64 `json:"end"`
    Text             string  `json:"text"`
    Tokens           []int   `json:"tokens"`
    Temperature      int     `json:"temperature"`
    AvgLogprob       float64 `json:"avg_logprob"`
    CompressionRatio float64 `json:"compression_ratio"`
    NoSpeechProb     float64 `json:"no_speech_prob"`
}

