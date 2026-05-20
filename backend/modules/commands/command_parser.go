package commands

import (
	"jarvis-ai/modules/ai"
	"jarvis-ai/modules/system"
	"log"
	"strings"
)

// normalizeInput lowercases the input and strips the wake-word prefix
// in all three supported languages (English, Hindi, Marathi).
func normalizeInput(original string) string {
	input := strings.ToLower(strings.TrimSpace(original))

	// Strip wake-word prefix — all language variants
	wakeWords := []string{
		"jarvis ", "hey jarvis ", "ok jarvis ", "अरे जार्विस ", "जार्विस ", "hey rudra ", "rudra ",
	}
	for _, ww := range wakeWords {
		input = strings.TrimPrefix(input, ww)
	}
	return strings.TrimSpace(input)
}

// ── Hindi keyword helpers ────────────────────────────────────────────────────

func containsAny(input string, keywords []string) bool {
	for _, k := range keywords {
		if strings.Contains(input, k) {
			return true
		}
	}
	return false
}

// Hindi/Marathi open synonyms
var openKW = []string{"खोलो", "खोल", "उघड", "open"}

// Hindi/Marathi close synonyms
var closeKW = []string{"बंद करो", "बंद कर", "बंद", "close", "shut"}

// Hindi/Marathi search synonyms
var searchKW = []string{"खोजो", "ढूंढो", "खोज", "शोध", "search"}

// Hindi/Marathi play synonyms
var playKW = []string{"चलाओ", "चालू कर", "play"}

// Hindi/Marathi pause synonyms
var pauseKW = []string{"रोको", "थांब", "pause"}

// Hindi/Marathi next synonyms
var nextKW = []string{"अगला", "next"}

// Hindi/Marathi previous synonyms
var prevKW = []string{"पिछला", "previous"}

// Hindi/Marathi volume up synonyms
var volUpKW = []string{"आवाज़ बढ़ाओ", "आवाज बढ़ाओ", "आवाज वाढव", "volume up"}

// Hindi/Marathi volume down synonyms
var volDownKW = []string{"आवाज़ कम करो", "आवाज कम करो", "आवाज कमी कर", "volume down"}

// Hindi/Marathi mute synonyms
var muteKW = []string{"म्यूट करो", "म्यूट कर", "mute"}

// Hindi/Marathi shutdown synonyms
var shutdownKW = []string{"बंद करो सब", "shutdown", "बंद करना"}

// Hindi/Marathi restart synonyms
var restartKW = []string{"restart", "पुनः शुरू करो"}

// Hindi/Marathi lock synonyms
var lockKW = []string{"लॉक करो", "लॉक कर", "lock"}

// Hindi/Marathi sleep synonyms
var sleepKW = []string{"स्लीप", "sleep", "सो जाओ"}

// Hindi/Marathi screenshot synonyms
var screenshotKW = []string{"screenshot", "स्क्रीनशॉट", "फोटो"}

// ── Main parser ──────────────────────────────────────────────────────────────

func ParseAndExecute(input string) string {
	originalInput := input
	input = normalizeInput(input)
	log.Printf("Parsing command: %s", input)

	// Split multiple commands separated by " and " / " और " / " आणि "
	separators := []string{" and ", " और ", " आणि "}
	for _, sep := range separators {
		if strings.Contains(input, sep) {
			parts := strings.Split(input, sep)
			validSplit := true
			for _, part := range parts {
				if !containsActionMultilang(strings.TrimSpace(part)) {
					validSplit = false
					break
				}
			}
			if validSplit && len(parts) > 1 {
				var results []string
				for _, part := range parts {
					res := ParseAndExecute(strings.TrimSpace(part))
					results = append(results, res)
				}
				return strings.Join(results, ". then, ")
			}
		}
	}

	// WhatsApp Commands (EN + HI + MR)
	whatsappKW := []string{"whatsapp", "व्हाट्सएप", "व्हाट्सअप"}
	callKW := []string{"call", "कॉल", "फोन"}
	if containsAny(input, whatsappKW) {
		if containsAny(input, callKW) {
			contactName := extractContactName(input)
			return system.WhatsAppCall(contactName)
		}
		return system.OpenWhatsApp()
	}

	// Google Search
	googleKW := []string{"google", "गूगल"}
	if containsAny(input, searchKW) && containsAny(input, googleKW) {
		query := extractSearchQuery(input)
		return system.GoogleSearch(query)
	}

	// YouTube Commands
	youtubeKW := []string{"youtube", "यूट्यूब"}
	if containsAny(input, youtubeKW) {
		if containsAny(input, searchKW) {
			query := extractSearchQuery(input)
			return system.YouTubeSearch(query)
		}
		return system.OpenYouTube()
	}

	// Gmail
	gmailKW := []string{"gmail", "email", "ईमेल", "मेल"}
	if containsAny(input, gmailKW) {
		return system.OpenGmail()
	}

	// Open Website
	if strings.HasPrefix(input, "open") && (strings.Contains(input, ".com") || strings.Contains(input, ".org") || strings.Contains(input, ".net")) {
		website := strings.TrimSpace(strings.TrimPrefix(input, "open"))
		return system.OpenWebsite(website)
	}

	// Open Folder (EN + HI + MR)
	folderKW := []string{"folder", "go to", "downloads", "documents", "pictures", "desktop",
		"फ़ोल्डर", "फोल्डर", "डाउनलोड", "दस्तावेज़", "चित्र", "डेस्कटॉप"}
	if containsAny(input, folderKW) {
		folderName := extractFolderName(input)
		return system.OpenFolder(folderName)
	}

	// Volume Up
	if containsAny(input, volUpKW) {
		return system.ControlVolume("volume up")
	}

	// Volume Down
	if containsAny(input, volDownKW) {
		return system.ControlVolume("volume down")
	}

	// Mute
	if containsAny(input, muteKW) {
		return system.ControlVolume("mute")
	}

	// General Volume keyword
	if strings.Contains(input, "volume") {
		return system.ControlVolume(input)
	}

	// Brightness Control
	brightnessKW := []string{"brightness", "चमक", "प्रकाश"}
	if containsAny(input, brightnessKW) {
		return system.ControlBrightness(input)
	}

	// Power Commands
	if containsAny(input, shutdownKW) {
		return system.PowerControl("shutdown")
	}
	if containsAny(input, restartKW) {
		return system.PowerControl("restart")
	}
	if containsAny(input, lockKW) {
		return system.PowerControl("lock")
	}
	if containsAny(input, sleepKW) {
		return system.PowerControl("sleep")
	}

	// Screenshot
	if containsAny(input, screenshotKW) {
		return system.PowerControl("screenshot")
	}

	// Media — Pause / Next / Previous
	if containsAny(input, pauseKW) {
		return system.MediaControl("pause")
	}
	if containsAny(input, nextKW) {
		return system.MediaControl("next")
	}
	if containsAny(input, prevKW) {
		return system.MediaControl("previous")
	}

	// Media — Play (with optional song name)
	musicKW := []string{"music", "song", "गाना", "संगीत", "गीत"}
	if containsAny(input, playKW) {
		if !containsAny(input, pauseKW) && !containsAny(input, nextKW) && !containsAny(input, prevKW) {
			stripped := input
			for _, kw := range append(playKW, "the song ", "song ", "गाना ") {
				stripped = strings.TrimPrefix(stripped, kw)
			}
			stripped = strings.TrimSpace(stripped)
			if stripped != "" && !containsAny(stripped, musicKW) {
				return system.PlaySpecificSong(stripped)
			}
		}
		return system.MediaControl(input)
	}
	if containsAny(input, musicKW) {
		if containsAny(input, openKW) && !containsAny(input, playKW) {
			return system.OpenMusicPlayer()
		}
		return system.MediaControl(input)
	}

	// Open Application (EN + HI + MR open keywords)
	if containsAny(input, openKW) {
		// Strip all open synonyms to get the app name
		appName := input
		for _, kw := range openKW {
			appName = strings.TrimPrefix(appName, kw+" ")
			appName = strings.TrimSuffix(appName, " "+kw)
		}
		appName = strings.TrimSpace(appName)
		return system.OpenApp(appName)
	}

	// Write Code
	writeCodeKW := []string{
		"write the code for", "write code for", "write a program to",
		"कोड लिखो", "कोड लिखा", "प्रोग्राम लिखो",
	}
	for _, kw := range writeCodeKW {
		if strings.HasPrefix(input, kw) {
			prompt := strings.TrimSpace(strings.TrimPrefix(input, kw))
			return system.WriteCode(prompt)
		}
	}

	// Default fallback to AI processing
	response, err := ai.QueryOllama(originalInput)
	if err != nil {
		log.Println("AI Query Error:", err)
		return "I encountered an error connecting to my AI core."
	}
	return response
}

// Helper function to extract contact name from WhatsApp commands
func extractContactName(input string) string {
	callKWs := []string{"call", "कॉल", "फोन"}
	for _, kw := range callKWs {
		if strings.Contains(input, kw) {
			parts := strings.SplitN(input, kw, 2)
			if len(parts) > 1 {
				name := strings.TrimSpace(parts[1])
				name = strings.TrimSuffix(name, "on whatsapp")
				name = strings.TrimSuffix(name, "व्हाट्सएप पर")
				name = strings.TrimSpace(name)
				if name != "" {
					return name
				}
			}
		}
	}
	return "contact"
}

// Helper function to extract search query
func extractSearchQuery(input string) string {
	query := input
	prefixes := []string{"search", "for", "खोजो", "ढूंढो", "खोज", "शोध"}
	suffixes := []string{"on google", "on youtube", "गूगल पर", "यूट्यूब पर"}
	for _, p := range prefixes {
		query = strings.TrimPrefix(query, p+" ")
	}
	for _, s := range suffixes {
		query = strings.TrimSuffix(query, " "+s)
	}
	return strings.TrimSpace(query)
}

// Helper function to extract folder name
func extractFolderName(input string) string {
	folderMap := map[string]string{
		"downloads": "downloads",
		"डाउनलोड":   "downloads",
		"documents": "documents",
		"दस्तावेज़": "documents",
		"pictures":  "pictures",
		"photos":    "pictures",
		"चित्र":     "pictures",
		"desktop":   "desktop",
		"डेस्कटॉप":  "desktop",
		"music":     "music",
		"संगीत":     "music",
		"videos":    "videos",
		"वीडियो":    "videos",
	}
	for keyword, folderName := range folderMap {
		if strings.Contains(input, keyword) {
			return folderName
		}
	}

	// Try to extract folder name after "open", "folder", or "go to"
	parts := strings.Fields(input)
	for i := 0; i < len(parts); i++ {
		if parts[i] == "open" || parts[i] == "folder" || parts[i] == "to" || parts[i] == "खोलो" || parts[i] == "उघड" {
			if i+1 < len(parts) && parts[i+1] != "the" {
				return parts[i+1]
			} else if i+2 < len(parts) && parts[i+1] == "the" {
				return parts[i+2]
			}
		}
	}
	return "downloads"
}

// Helper function to check if a command string contains a recognizable action verb
// (used for multi-command splitting decision)
func containsActionMultilang(cmd string) bool {
	keywords := []string{
		// English
		"open", "go to", "play", "search", "shutdown", "restart", "lock", "sleep", "call", "volume", "mute", "write",
		// Hindi
		"खोलो", "खोल", "चलाओ", "खोजो", "ढूंढो", "बंद", "लॉक", "कॉल", "आवाज़", "म्यूट", "कोड",
		// Marathi
		"उघड", "चालू", "शोध", "लॉक कर", "म्यूट कर",
	}
	cmd = strings.ToLower(cmd)
	return containsAny(cmd, keywords)
}
