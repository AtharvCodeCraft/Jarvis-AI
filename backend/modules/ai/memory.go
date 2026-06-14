package ai

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"jarvis-ai/modules/database"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Message matches the database and Ollama role structures
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// SaveMessage stores a message in the SQLite database
func SaveMessage(convID, sender, content string, agentAssigned string) error {
	id := uuid.New().String()
	query := `INSERT INTO messages (id, conversation_id, sender, content, agent_assigned, created_at) 
	          VALUES (?, ?, ?, ?, ?, ?)`
	_, err := database.DB.Exec(query, id, convID, sender, content, agentAssigned, time.Now())
	if err != nil {
		log.Printf("Error saving message: %v", err)
	}
	return err
}

// GetConversationHistory retrieves the last N messages for a conversation
func GetConversationHistory(convID string, limit int) ([]Message, error) {
	query := `SELECT sender, content FROM messages WHERE conversation_id = ? ORDER BY created_at ASC LIMIT ?`
	rows, err := database.DB.Query(query, convID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []Message
	for rows.Next() {
		var sender, content string
		if err := rows.Scan(&sender, &content); err != nil {
			return nil, err
		}
		// Map sender DB value ('user', 'assistant') to Ollama standard roles ('user', 'assistant')
		role := sender
		if role == "system" {
			role = "system"
		}
		history = append(history, Message{Role: role, Content: content})
	}
	return history, nil
}

// GetOrCreateConversation fetches the latest active conversation or creates one
func GetOrCreateConversation() string {
	var convID string
	query := `SELECT id FROM conversations ORDER BY updated_at DESC LIMIT 1`
	err := database.DB.QueryRow(query).Scan(&convID)
	if err == sql.ErrNoRows {
		convID = uuid.New().String()
		insertQuery := `INSERT INTO conversations (id, title, created_at, updated_at) VALUES (?, ?, ?, ?)`
		_, err = database.DB.Exec(insertQuery, convID, "Jarvis Interaction Session", time.Now(), time.Now())
		if err != nil {
			log.Printf("Error creating initial conversation: %v", err)
		}
	}
	return convID
}

// SaveSemanticMemory saves a preference or general factual node to SQLite memory_nodes
func SaveSemanticMemory(content string, memoryType string) error {
	id := uuid.New().String()
	query := `INSERT INTO memory_nodes (id, user_id, memory_type, content, last_accessed, created_at)
	          VALUES (?, ?, ?, ?, ?, ?)`
	_, err := database.DB.Exec(query, id, "default_user", memoryType, content, time.Now(), time.Now())
	if err != nil {
		log.Printf("Error saving memory node: %v", err)
	}
	return err
}

// GetRelevantMemories queries the SQLite database for memory context matching the keywords
func GetRelevantMemories(prompt string) string {
	// Simple keyword extraction for matching database contents (local-first fallback)
	words := strings.Fields(strings.ToLower(prompt))
	if len(words) == 0 {
		return ""
	}

	// Dynamic SQLite LIKE lookup
	queryParts := make([]string, len(words))
	args := make([]interface{}, len(words))
	for i, word := range words {
		queryParts[i] = "content LIKE ?"
		args[i] = "%" + word + "%"
	}

	query := fmt.Sprintf("SELECT content FROM memory_nodes WHERE %s LIMIT 5", strings.Join(queryParts, " OR "))
	rows, err := database.DB.Query(query, args...)
	if err != nil {
		log.Printf("Error querying memories: %v", err)
		return ""
	}
	defer rows.Close()

	var memories []string
	for rows.Next() {
		var content string
		if err := rows.Scan(&content); err == nil {
			memories = append(memories, content)
		}
	}

	if len(memories) > 0 {
		return "\n[Retrieved User Memories]:\n- " + strings.Join(memories, "\n- ")
	}
	return ""
}

// SaveInteraction satisfies legacy interface requirements
func SaveInteraction(userMsg, aiMsg string) {
	convID := GetOrCreateConversation()
	SaveMessage(convID, "user", userMsg, "")
	SaveMessage(convID, "assistant", aiMsg, "Orchestrator")
	
	// Auto-extract preference logic in background if keyword matched
	if strings.Contains(strings.ToLower(userMsg), "i prefer") || strings.Contains(strings.ToLower(userMsg), "my name is") || strings.Contains(strings.ToLower(userMsg), "remember that") {
		go func() {
			SaveSemanticMemory(userMsg, "preference")
			log.Println("New preference saved to semantic memory node.")
		}()
	}
}

// GetContext returns system parameters, localization rules and retrieved long-term memories
func GetContext(prompt string) string {
	memContext := GetRelevantMemories(prompt)

	baseContext := `You are Jarvis (also called Rudra), a highly advanced local AI assistant managing this system. Keep your responses concise.

IMPORTANT — Language Rules:
- Detect the language of the user's message automatically.
- If the user writes or speaks in Hindi (हिंदी), you MUST reply in Hindi.
- If the user writes or speaks in Marathi (मराठी), you MUST reply in Marathi.
- If the user writes or speaks in English, reply in English.
- Never switch languages on your own. Always match the user's language exactly.
- You are fluent in English, Hindi, and Marathi.`

	if memContext != "" {
		baseContext += "\n" + memContext
	}

	return baseContext
}

// GenerateEmbeddings returns local vector representations using Ollama's embedding engine
func GenerateEmbeddings(text string) ([]float32, error) {
	type EmbedRequest struct {
		Model  string `json:"model"`
		Prompt string `json:"prompt"`
	}
	type EmbedResponse struct {
		Embedding []float32 `json:"embedding"`
	}

	reqBody := EmbedRequest{
		Model:  "nomic-embed-text",
		Prompt: text,
	}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", "http://127.0.0.1:11434/api/embeddings", strings.NewReader(string(jsonBody)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("embedding API returned status: %d", resp.StatusCode)
	}

	var embedResp EmbedResponse
	if err := json.NewDecoder(resp.Body).Decode(&embedResp); err != nil {
		return nil, err
	}

	return embedResp.Embedding, nil
}
