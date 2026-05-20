package voice

import (
	"log"
)

func TranscribeAudio(audioPath string) string {
	log.Println("Transcribing audio using local Whisper model:", audioPath)
	// Example stub. Real implementation would run whisper.cpp or similar binaries
	return "Transcribed text from user command."
}
