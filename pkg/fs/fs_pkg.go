package fs

import (
	"fmt"
	"os"
)

func WritePCMDataToFile(filepath string, pcmData []byte) error {
    f, err := os.Create(filepath)
    if err != nil {
        return fmt.Errorf("failed to create file %s: %w", filepath, err)
    }
    defer f.Close()

    n, err := f.Write(pcmData)
    if err != nil {
        return fmt.Errorf("failed to write PCM data: %w", err)
    }

    if n != len(pcmData) {
        return fmt.Errorf("incomplete write: wrote %d bytes out of %d", n, len(pcmData))
    }

    // Ensure data is flushed to disk
    if err := f.Sync(); err != nil {
        return fmt.Errorf("failed to sync file: %w", err)
    }

    return nil
}

// saveWAVFile saves the WAV data to the filesystem with the given filename.
func WriteWAVDataToFile(filename string, data []byte) error {
    // Create or truncate the file
    file, err := os.Create(filename)
    if err != nil {
        return fmt.Errorf("failed to create file %s: %w", filename, err)
    }
    defer file.Close()

    _, err = file.Write(data)
    if err != nil {
        return fmt.Errorf("failed to write data to file %s: %w", filename, err)
    }

    return nil
}
