package tts

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/haguro/elevenlabs-go"
)

func ElevenLabsTextToSpeech(text string) ([]byte, error) {

	env := os.Getenv("ELEVENLABS_API_KEY")

	client := elevenlabs.NewClient(context.Background(), env, 30*time.Second)

	ttsReq := elevenlabs.TextToSpeechRequest{
		Text:    text,
		ModelID: "eleven_monolingual_v1",
	}

	audio, err := client.TextToSpeech("cgSgspJ2msm6clMCkdW9", ttsReq, elevenlabs.OutputFormat("pcm_16000"))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	log.Println("Successfully generated PCM audio data")
	return audio, err
}
