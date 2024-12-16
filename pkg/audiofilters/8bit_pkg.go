package audiofilters

import (
	"fmt"
	"os/exec"
)

func ConvertTo8Bit(inputPath, outputPath string) error {
    // Construct the ffmpeg command with appropriate flags
    cmd := exec.Command("ffmpeg", "-i", inputPath, "-af", "[0:a]acrusher=samples=20:bits=8,atrim=start=0.5,apulsator=mode=sine:hz=3:width=0.1:offset_r=0[out]", outputPath)

    // Capture combined standard output and standard error
    output, err := cmd.CombinedOutput()
    if err != nil {
        return fmt.Errorf("FFmpeg conversion failed: %v, Output: %s", err, string(output))
    }
    return nil
}
