package tool

import (
	"errors"
	"jarvis-ai/modules/system"
)

func init() {
	// Register: Open Application
	GlobalRegistry.Register(
		"open_application",
		"Launch a Windows desktop application (e.g. chrome, notepad, calculator, settings, etc.)",
		map[string]ToolParam{
			"app_name": {
				Type:        "string",
				Description: "The name of the application to launch (e.g. 'chrome', 'notepad', 'calc')",
				Required:    true,
			},
		},
		false, // Does not require HITL confirmation (safe launcher)
		func(args map[string]interface{}) (string, error) {
			appName, ok := args["app_name"].(string)
			if !ok {
				return "", errors.New("app_name argument is missing or not a string")
			}
			return system.OpenApp(appName), nil
		},
	)

	// Register: Web Search (Google)
	GlobalRegistry.Register(
		"google_search",
		"Search the web using Google",
		map[string]ToolParam{
			"query": {
				Type:        "string",
				Description: "The search terms or query to run",
				Required:    true,
			},
		},
		false,
		func(args map[string]interface{}) (string, error) {
			query, ok := args["query"].(string)
			if !ok {
				return "", errors.New("query argument is missing or not a string")
			}
			return system.GoogleSearch(query), nil
		},
	)

	// Register: YouTube Search
	GlobalRegistry.Register(
		"youtube_search",
		"Search for a video or music track on YouTube",
		map[string]ToolParam{
			"query": {
				Type:        "string",
				Description: "The video name, channel name, or search query",
				Required:    true,
			},
		},
		false,
		func(args map[string]interface{}) (string, error) {
			query, ok := args["query"].(string)
			if !ok {
				return "", errors.New("query argument is missing or not a string")
			}
			return system.YouTubeSearch(query), nil
		},
	)

	// Register: Open Website
	GlobalRegistry.Register(
		"open_website",
		"Navigate to a website URL in the default browser",
		map[string]ToolParam{
			"url": {
				Type:        "string",
				Description: "The URL of the website to open (e.g. 'github.com', 'gmail.com')",
				Required:    true,
			},
		},
		false,
		func(args map[string]interface{}) (string, error) {
			urlStr, ok := args["url"].(string)
			if !ok {
				return "", errors.New("url argument is missing or not a string")
			}
			return system.OpenWebsite(urlStr), nil
		},
	)

	// Register: Navigate Folders
	GlobalRegistry.Register(
		"open_folder",
		"Open a common Windows user folder path (downloads, documents, pictures, desktop, music, videos)",
		map[string]ToolParam{
			"folder_name": {
				Type:        "string",
				Description: "The target folder name (downloads, documents, pictures, desktop, music, videos)",
				Required:    true,
				Enum:        []string{"downloads", "documents", "pictures", "desktop", "music", "videos"},
			},
		},
		false,
		func(args map[string]interface{}) (string, error) {
			folderName, ok := args["folder_name"].(string)
			if !ok {
				return "", errors.New("folder_name argument is missing or not a string")
			}
			return system.OpenFolder(folderName), nil
		},
	)

	// Register: Control Volume
	GlobalRegistry.Register(
		"control_volume",
		"Adjust or mute system audio volume",
		map[string]ToolParam{
			"action": {
				Type:        "string",
				Description: "The action to execute: 'up' (increase), 'down' (decrease), or 'mute' (toggle mute)",
				Required:    true,
				Enum:        []string{"up", "down", "mute"},
			},
		},
		false,
		func(args map[string]interface{}) (string, error) {
			action, ok := args["action"].(string)
			if !ok {
				return "", errors.New("action argument is missing or not a string")
			}
			return system.ControlVolume(action), nil
		},
	)

	// Register: Power Management
	GlobalRegistry.Register(
		"power_control",
		"Trigger OS power states (lock, sleep, restart, shutdown). Warning: Destructive actions require user prompt confirmation.",
		map[string]ToolParam{
			"state": {
				Type:        "string",
				Description: "The desired system power state: 'lock', 'sleep', 'restart', or 'shutdown'",
				Required:    true,
				Enum:        []string{"lock", "sleep", "restart", "shutdown"},
			},
		},
		true, // Requires Human-in-the-loop validation
		func(args map[string]interface{}) (string, error) {
			state, ok := args["state"].(string)
			if !ok {
				return "", errors.New("state argument is missing or not a string")
			}
			return system.PowerControl(state), nil
		},
	)
	
	// Register: WhatsApp Chat / Call Launcher
	GlobalRegistry.Register(
		"whatsapp_action",
		"Initiate a chat or call with a contact on WhatsApp Web",
		map[string]ToolParam{
			"contact_name": {
				Type:        "string",
				Description: "The name of the contact/group to call or chat with",
				Required:    true,
			},
			"action": {
				Type:        "string",
				Description: "The WhatsApp action: 'chat' or 'call'",
				Required:    true,
				Enum:        []string{"chat", "call"},
			},
		},
		false,
		func(args map[string]interface{}) (string, error) {
			contactName, ok := args["contact_name"].(string)
			if !ok {
				return "", errors.New("contact_name argument is missing or not a string")
			}
			action, ok := args["action"].(string)
			if !ok {
				return "", errors.New("action argument is missing or not a string")
			}
			
			if action == "call" {
				return system.WhatsAppCall(contactName), nil
			}
			return system.WhatsAppChat(contactName), nil
		},
	)
}
