package pcm

import (
	"bytes"
	"encoding/binary"
)

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
