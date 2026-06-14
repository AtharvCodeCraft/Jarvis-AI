package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"jarvis-ai/modules/database"
	"jarvis-ai/modules/tool"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

// OllamaChatRequest schema for /api/chat
type OllamaChatRequest struct {
	Model    string              `json:"model"`
	Messages []Message           `json:"messages"`
	Stream   bool                `json:"stream"`
	Tools    []OllamaToolWrapper `json:"tools,omitempty"`
}

type OllamaToolWrapper struct {
	Type     string             `json:"type"`
	Function OllamaToolFunction `json:"function"`
}

type OllamaToolFunction struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  OllamaParamsDefinition `json:"parameters"`
}

type OllamaParamsDefinition struct {
	Type       string                 `json:"type"`
	Properties map[string]interface{} `json:"properties"`
	Required   []string               `json:"required,omitempty"`
}

// OllamaChatResponse schema for /api/chat
type OllamaChatResponse struct {
	Model     string      `json:"model"`
	Message   RespMessage `json:"message"`
	Done      bool        `json:"done"`
}

type RespMessage struct {
	Role      string     `json:"role"`
	Content   string     `json:"content"`
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
}

type ToolCall struct {
	Function FunctionCallDetail `json:"function"`
}

type FunctionCallDetail struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
}

// LogToolExecution records tool audits into the SQLite database
func LogToolExecution(agentName, toolName string, args map[string]interface{}, status string, err error, durationMs int64) {
	id := uuid.New().String()
	argsJSON, _ := json.Marshal(args)
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}

	query := `INSERT INTO tool_logs (id, agent_name, tool_name, arguments, execution_status, error_message, duration_ms, created_at)
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	_, dbErr := database.DB.Exec(query, id, agentName, toolName, string(argsJSON), status, errMsg, durationMs, time.Now())
	if dbErr != nil {
		log.Printf("Failed to write tool log to DB: %v", dbErr)
	}
}

// QueryOllama performs functional routing using local Ollama model (with system prompt, memory & tools)
func QueryOllama(prompt string) (string, error) {
	log.Println("Agent Orchestrator: processing user query...")

	convID := GetOrCreateConversation()
	
	// Fetch previous messages for short-term history
	history, err := GetConversationHistory(convID, 10)
	if err != nil {
		log.Printf("Failed to load conversation history: %v", err)
	}

	// Prepare current prompt context with user preferences retrieved from semantic memory
	systemPrompt := GetContext(prompt)

	// Build all messages list
	var messages []Message
	messages = append(messages, Message{Role: "system", Content: systemPrompt})
	messages = append(messages, history...)
	messages = append(messages, Message{Role: "user", Content: prompt})

	// Wrap registered tools into Ollama schema structure
	var ollamaTools []OllamaToolWrapper
	registeredTools := tool.GlobalRegistry.ListDefinitions()
	for _, rt := range registeredTools {
		props := make(map[string]interface{})
		var required []string

		for pName, pDef := range rt.Parameters {
			paramMap := map[string]interface{}{
				"type":        pDef.Type,
				"description": pDef.Description,
			}
			if len(pDef.Enum) > 0 {
				paramMap["enum"] = pDef.Enum
			}
			props[pName] = paramMap
			if pDef.Required {
				required = append(required, pName)
			}
		}

		ollamaTools = append(ollamaTools, OllamaToolWrapper{
			Type: "function",
			Function: OllamaToolFunction{
				Name:        rt.Name,
				Description: rt.Description,
				Parameters: OllamaParamsDefinition{
					Type:       "object",
					Properties: props,
					Required:   required,
				},
			},
		})
	}

	reqBody := OllamaChatRequest{
		Model:    GetAvailableModel(), // Dynamically select local available model
		Messages: messages,
		Stream:   false,
		Tools:    ollamaTools,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", "http://127.0.0.1:11434/api/chat", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("ollama API status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var chatResp OllamaChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return "", err
	}

	// 1. Check if LLM requested a Tool Call
	if len(chatResp.Message.ToolCalls) > 0 {
		toolCall := chatResp.Message.ToolCalls[0]
		toolName := toolCall.Function.Name
		toolArgsMap := toolCall.Function.Arguments
		toolArgsJSON, _ := json.Marshal(toolArgsMap)

		log.Printf("[Orchestrator] LLM routed intent to tool: %s with args: %s", toolName, string(toolArgsJSON))

		// Check if safety policy permits immediate execution or requires UI confirmation (HITL)
		tObj, exists := tool.GlobalRegistry.Get(toolName)
		if exists && tObj.Definition.RequiresHITL {
			// Save log as pending
			LogToolExecution("Orchestrator", toolName, toolArgsMap, "pending", nil, 0)
			return fmt.Sprintf("CONFIRMATION_REQUIRED:%s:%s", toolName, string(toolArgsJSON)), nil
		}

		// Execute tool synchronously
		startTime := time.Now()
		result, runErr := tool.GlobalRegistry.Invoke(toolName, string(toolArgsJSON))
		duration := time.Since(startTime).Milliseconds()

		status := "success"
		if runErr != nil {
			status = "failed"
		}
		
		// Audit log to DB
		LogToolExecution("Orchestrator", toolName, toolArgsMap, status, runErr, duration)

		if runErr != nil {
			return fmt.Sprintf("Error executing tool %s: %v", toolName, runErr), nil
		}

		// Save the transaction history to memory
		SaveInteraction(prompt, result)
		return result, nil
	}

	// 2. Otherwise handle default conversational response
	reply := chatResp.Message.Content
	SaveInteraction(prompt, reply)

	return reply, nil
}

// GetAvailableModel returns the best available model, fallback to a default
func GetAvailableModel() string {
	type ListResponse struct {
		Models []struct {
			Name string `json:"name"`
		} `json:"models"`
	}

	resp, err := http.Get("http://127.0.0.1:11434/api/tags")
	if err != nil {
		return "qwen2.5:1.5b"
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		var list ListResponse
		if err := json.NewDecoder(resp.Body).Decode(&list); err == nil && len(list.Models) > 0 {
			// Search preference order
			order := []string{"qwen2.5:1.5b", "qwen2.5:0.5b", "qwen2.5-coder:7b", "llama3.2:1b"}
			for _, preferred := range order {
				for _, m := range list.Models {
					if m.Name == preferred || strings.HasPrefix(m.Name, preferred) {
						return m.Name
					}
				}
			}
			// If none of preferred found, use the first available model
			return list.Models[0].Name
		}
	}
	return "qwen2.5:1.5b"
}
