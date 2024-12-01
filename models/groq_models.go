package models

type GroqLLMResponse struct {
    Choices []struct {
        Message struct {
            Role    string `json:"role"`
            Content string `json:"content"`
        } `json:"message"`
    } `json:"choices"`
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

