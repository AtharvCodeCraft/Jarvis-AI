package commands

import (
	"jarvis-ai/modules/agent"
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

func ParseAndExecute(input string) string {
	normalized := normalizeInput(input)
	log.Printf("[CommandRouter] Normalized request: %s", normalized)

	// Dynamically delegate request execution directly to the specialized Agent Orchestrator
	response, err := agent.ExecuteAgent(normalized)
	if err != nil {
		log.Println("[CommandRouter] Error invoking Agent Orchestrator:", err)
		return "I encountered an error executing that request. Please verify that Ollama is running locally."
	}
	
	return response
}
