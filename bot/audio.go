package bot

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os/exec"
	"time"
)

func convertOggToMp3(oggData []byte, duration time.Duration) ([]byte, error) {
	// Create a byte reader with the ogg data
	reader := bytes.NewReader(oggData)

	// Use FFmpeg to convert the OGG data to MP3 with the specified duration
	cmd := exec.Command(
		"ffmpeg",
		"-i", "pipe:0",
		"-f", "mp3",
		"-t", duration.String(),
		"-acodec", "libmp3lame",
		"pipe:1",
	)
	cmd.Stdin = reader
	var mp3Buffer bytes.Buffer
	cmd.Stdout = &mp3Buffer
	cmd.Stderr = ioutil.Discard

	// Execute the command
	err := cmd.Run()
	if err != nil {
		return nil, errors.New("Error converting OGG to MP3: " + err.Error())
	}

	return mp3Buffer.Bytes(), nil
}
