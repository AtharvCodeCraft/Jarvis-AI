package system

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
)

// OpenFolder opens a specific Windows folder
func OpenFolder(folderName string) string {
	folderName = strings.ToLower(folderName)

	var folderPath string
	userProfile := os.Getenv("USERPROFILE")

	switch folderName {
	case "downloads", "download":
		folderPath = filepath.Join(userProfile, "Downloads")
	case "documents", "document", "docs":
		folderPath = filepath.Join(userProfile, "Documents")
	case "pictures", "picture", "photos", "photo", "images":
		folderPath = filepath.Join(userProfile, "Pictures")
	case "desktop":
		folderPath = filepath.Join(userProfile, "Desktop")
	case "music":
		folderPath = filepath.Join(userProfile, "Music")
	case "videos", "video":
		folderPath = filepath.Join(userProfile, "Videos")
	default:
		return "Folder not recognized. Try: downloads, documents, pictures, desktop, music, or videos"
	}

	cmd := exec.Command("cmd", "/c", "start", folderPath)
	if runtime.GOOS == "windows" {
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	}
	err := cmd.Start()
	if err != nil {
		log.Printf("Failed to open folder %s: %v", folderPath, err)
		return "Failed to open " + folderName + " folder"
	}

	return "Opening " + folderName + " folder"
}
