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

	// Create a new client
	client := elevenlabs.NewClient(context.Background(), env, 30*time.Second)

	// Create a TextToSpeechRequest
	ttsReq := elevenlabs.TextToSpeechRequest{
	Text:    text,
	ModelID: "eleven_monolingual_v1",
	}

	// Call the TextToSpeech method on the client, using the "Adam"'s voice ID.
	audio, err := client.TextToSpeech("cgSgspJ2msm6clMCkdW9", ttsReq, elevenlabs.OutputFormat("pcm_16000"))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	// // Create a TextToSpeechRequest
	// ttsReq := elevenlabs.TextToSpeechRequest{
	// Text:    text,
	// ModelID: "eleven_monolingual_v1",
	// }

	// // Call the TextToSpeech method on the client, using the "Adam"'s voice ID.
	// audio, err := client.TextToSpeech("cgSgspJ2msm6clMCkdW9", ttsReq, elevenlabs.OutputFormat("pcm_16000"))
	// if err != nil {
	// 	log.Fatal(err)
	// 	return nil, err
	// }

	// // Write the audio file bytes to disk
	// if err := os.WriteFile("adam.mp3", audio, 0644); err != nil {
	// log.Fatal(err)
	// }

	log.Println("Successfully generated PCM audio data")
	return audio, err
}