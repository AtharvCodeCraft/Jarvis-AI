package system

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
)

// MediaControl handles media playback commands
func MediaControl(command string) string {
	command = strings.ToLower(command)

	var cmd *exec.Cmd

	if strings.Contains(command, "play") {
		// Send Play/Pause media key
		cmd = exec.Command("powershell", "-Command", "(New-Object -ComObject WScript.Shell).SendKeys([char]179)")
		runHidden(cmd)
		return "Playing music"
	} else if strings.Contains(command, "pause") || strings.Contains(command, "stop") {
		// Send Play/Pause media key (toggles)
		cmd = exec.Command("powershell", "-Command", "(New-Object -ComObject WScript.Shell).SendKeys([char]179)")
		runHidden(cmd)
		return "Pausing music"
	} else if strings.Contains(command, "next") {
		// Send Next Track media key
		cmd = exec.Command("powershell", "-Command", "(New-Object -ComObject WScript.Shell).SendKeys([char]176)")
		runHidden(cmd)
		return "Playing next song"
	} else if strings.Contains(command, "previous") || strings.Contains(command, "prev") || strings.Contains(command, "back") {
		// Send Previous Track media key
		cmd = exec.Command("powershell", "-Command", "(New-Object -ComObject WScript.Shell).SendKeys([char]177)")
		runHidden(cmd)
		return "Playing previous song"
	}

	log.Printf("Media command not recognized: %s", command)
	return "Media command not recognized"
}

// OpenMusicPlayer opens Windows Media Player or Groove Music
func OpenMusicPlayer() string {
	// Try to open Groove Music (Windows 10/11 default)
	cmd := exec.Command("cmd", "/c", "start", "mswindowsmusic:")
	if runtime.GOOS == "windows" {
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	}
	err := cmd.Start()
	if err != nil {
		log.Printf("Failed to open music player: %v", err)
		return "Failed to open music player"
	}
	return "Opening music player"
}

// PlaySpecificSong searches for a song locally and fallback to YouTube
func PlaySpecificSong(songName string) string {
	songName = strings.TrimSpace(strings.ToLower(songName))

	// Strip common prefixes
	prefixes := []string{"play the song ", "play song ", "play "}
	for _, p := range prefixes {
		if strings.HasPrefix(songName, p) {
			songName = strings.TrimPrefix(songName, p)
			break
		}
	}

	psScript := fmt.Sprintf(`
$songName = "%s"
$paths = @(
    "$env:USERPROFILE\Music",
    "$env:USERPROFILE\Downloads",
    "$env:USERPROFILE\Desktop"
)
$found = $false
foreach ($path in $paths) {
    if (Test-Path $path) {
        $file = Get-ChildItem -Path $path -Filter "*$songName*.mp3" -Recurse -ErrorAction SilentlyContinue | Select-Object -First 1
        if (-not $file) {
            $file = Get-ChildItem -Path $path -Filter "*$songName*.wav" -Recurse -ErrorAction SilentlyContinue | Select-Object -First 1
        }
        if (-not $file) {
            $file = Get-ChildItem -Path $path -Filter "*$songName*.m4a" -Recurse -ErrorAction SilentlyContinue | Select-Object -First 1
        }
        if ($file) {
            Start-Process $file.FullName
            $found = $true
            exit 0
        }
    }
}
if (-not $found) {
    Start-Process "https://www.youtube.com/results?search_query=$([uri]::EscapeDataString($songName))"
}
`, songName)

	cmd := exec.Command("powershell", "-WindowStyle", "Hidden", "-Command", psScript)
	if runtime.GOOS == "windows" {
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	}
	cmd.Start()

	return "Playing " + songName
}
