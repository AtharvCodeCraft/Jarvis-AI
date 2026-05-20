package system

import (
	"fmt"
	"log"
	"net/url"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
)

// OpenWebsite opens a URL in the default browser
func OpenWebsite(urlStr string) string {
	if !strings.HasPrefix(urlStr, "http://") && !strings.HasPrefix(urlStr, "https://") {
		urlStr = "https://" + urlStr
	}

	cmd := exec.Command("cmd", "/c", "start", urlStr)
	if runtime.GOOS == "windows" {
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	}
	err := cmd.Start()
	if err != nil {
		log.Printf("Failed to open website %s: %v", urlStr, err)
		return "Failed to open website"
	}
	return "Opening " + urlStr
}

// GoogleSearch performs a Google search
func GoogleSearch(query string) string {
	searchURL := "https://www.google.com/search?q=" + url.QueryEscape(query)
	cmd := exec.Command("cmd", "/c", "start", searchURL)
	if runtime.GOOS == "windows" {
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	}
	err := cmd.Start()
	if err != nil {
		log.Printf("Failed to search Google: %v", err)
		return "Failed to perform Google search"
	}
	return "Searching Google for: " + query
}

// OpenWhatsApp opens WhatsApp Web
func OpenWhatsApp() string {
	return OpenWebsite("https://web.whatsapp.com")
}

// WhatsAppCall initiates a WhatsApp call (opens WhatsApp Web with call intent)
func WhatsAppCall(contactName string) string {
	// Open WhatsApp Web - user will need to manually initiate call
	// WhatsApp Web doesn't support direct deep linking to calls
	OpenWebsite("https://web.whatsapp.com")
	return fmt.Sprintf("Opening WhatsApp. Please search for %s and start the call", contactName)
}

// WhatsAppChat opens WhatsApp chat with a contact
func WhatsAppChat(contactName string) string {
	OpenWebsite("https://web.whatsapp.com")
	return fmt.Sprintf("Opening WhatsApp. Please search for %s to start chatting", contactName)
}

// OpenYouTube opens YouTube
func OpenYouTube() string {
	return OpenWebsite("https://www.youtube.com")
}

// OpenGmail opens Gmail
func OpenGmail() string {
	return OpenWebsite("https://mail.google.com")
}

// YouTubeSearch searches YouTube
func YouTubeSearch(query string) string {
	searchURL := "https://www.youtube.com/results?search_query=" + url.QueryEscape(query)
	cmd := exec.Command("cmd", "/c", "start", searchURL)
	if runtime.GOOS == "windows" {
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	}
	err := cmd.Start()
	if err != nil {
		log.Printf("Failed to search YouTube: %v", err)
		return "Failed to search YouTube"
	}
	return "Searching YouTube for: " + query
}
