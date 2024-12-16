package pcm

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
)

func TTStoWav(tts []byte, filePath string) error {
    if len(tts) == 0 {
        return fmt.Errorf("input TTS data is empty")
    }

    // Define WAV parameters matching the PCM data
    sampleRate := 16000 // 16 kHz
    bitsPerSample := 16
    numChannels := 1 // Mono

    // Create a new WAV file
    outFile, err := os.Create(filePath)
    if err != nil {
        return fmt.Errorf("failed to create WAV file: %w", err)
    }
    defer outFile.Close()

    // Initialize WAV encoder
    encoder := wav.NewEncoder(outFile, sampleRate, bitsPerSample, numChannels, 1)
    defer encoder.Close()

    // Create PCM buffer
    pcmBuffer := &audio.IntBuffer{
        Format: &audio.Format{
            SampleRate:  sampleRate,
            NumChannels: numChannels,
        },
        Data: make([]int, len(tts)/2), // 2 bytes per sample for 16-bit audio
    }

    // Convert PCM bytes to samples
    for i := 0; i < len(tts)-1; i += 2 {
        // Little-endian conversion
        sample := int(int16(tts[i]) | int16(tts[i+1])<<8)
        pcmBuffer.Data[i/2] = sample
    }

    // Write PCM data to WAV
    if err := encoder.Write(pcmBuffer); err != nil {
        return fmt.Errorf("failed to write PCM data to WAV: %w", err)
    }

    // Finalize the WAV file
    if err := encoder.Close(); err != nil {
        return fmt.Errorf("failed to finalize WAV file: %w", err)
    }

    return nil
}

func ToWAV(pcmData []byte) ([]byte, error) {
    var wavData bytes.Buffer

    // WAV file parameters
    numChannels := uint16(1)       // Mono
    sampleRate := uint32(16000)    // 16 kHz
    bitsPerSample := uint16(16)    // 16-bit samples
    byteRate := sampleRate * uint32(numChannels) * uint32(bitsPerSample/8)
    blockAlign := numChannels * bitsPerSample / 8
    dataSize := uint32(len(pcmData))

    // Write RIFF header
    wavData.WriteString("RIFF")
    binary.Write(&wavData, binary.LittleEndian, uint32(36+dataSize))
    wavData.WriteString("WAVE")

    // Write fmt subchunk
    wavData.WriteString("fmt ")
    binary.Write(&wavData, binary.LittleEndian, uint32(16))          // Subchunk1Size
    binary.Write(&wavData, binary.LittleEndian, uint16(1))           // AudioFormat (PCM)
    binary.Write(&wavData, binary.LittleEndian, numChannels)
    binary.Write(&wavData, binary.LittleEndian, sampleRate)
    binary.Write(&wavData, binary.LittleEndian, byteRate)
    binary.Write(&wavData, binary.LittleEndian, blockAlign)
    binary.Write(&wavData, binary.LittleEndian, bitsPerSample)

    // Write data subchunk
    wavData.WriteString("data")
    binary.Write(&wavData, binary.LittleEndian, dataSize)
    wavData.Write(pcmData)

    return wavData.Bytes(), nil
}


// ToWAV converts PCM data to WAV format.
// It assumes PCM data is 16-bit mono at 16 kHz. Adjust parameters as needed.
func M5ToWAVFile(pcmData []byte) ([]byte, error) {
    var wavData bytes.Buffer

    // WAV file parameters
    numChannels := uint16(1)       // Mono
    sampleRate := uint32(16000)    // 16 kHz
    bitsPerSample := uint16(16)    // 16-bit samples
    byteRate := sampleRate * uint32(numChannels) * uint32(bitsPerSample/8)
    blockAlign := numChannels * bitsPerSample / 8
    dataSize := uint32(len(pcmData))

    // Write RIFF header
    wavData.WriteString("RIFF")
    binary.Write(&wavData, binary.LittleEndian, uint32(36+dataSize))
    wavData.WriteString("WAVE")

    // Write fmt subchunk
    wavData.WriteString("fmt ")
    binary.Write(&wavData, binary.LittleEndian, uint32(16))          // Subchunk1Size
    binary.Write(&wavData, binary.LittleEndian, uint16(1))           // AudioFormat (PCM)
    binary.Write(&wavData, binary.LittleEndian, numChannels)
    binary.Write(&wavData, binary.LittleEndian, sampleRate)
    binary.Write(&wavData, binary.LittleEndian, byteRate)
    binary.Write(&wavData, binary.LittleEndian, blockAlign)
    binary.Write(&wavData, binary.LittleEndian, bitsPerSample)

    // Write data subchunk
    wavData.WriteString("data")
    binary.Write(&wavData, binary.LittleEndian, dataSize)
    wavData.Write(pcmData)

    return wavData.Bytes(), nil
}


