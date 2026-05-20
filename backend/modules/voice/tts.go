package voice

import (
	"log"
	"os/exec"
)

func Speak(text string) {
	log.Println("Speaking out loud:", text)
	// Fallback to built-in Windows speech synthesis (SAPI) with PowerShell
	// In production, integrate with Coqui TTS or Piper locally
	
	cmd := exec.Command("powershell", "-Command", "Add-Type -AssemblyName System.Speech; (New-Object System.Speech.Synthesis.SpeechSynthesizer).Speak('"+text+"')")
	err := cmd.Start()
	if err != nil {
		log.Println("TTS Error:", err)
	}
}
