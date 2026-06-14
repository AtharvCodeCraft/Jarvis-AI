package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"jarvis-ai/modules/ai"
	"jarvis-ai/modules/tool"
	"log"
	"net/http"
	"strings"
	"time"
)

type Agent struct {
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	SystemPrompt string   `json:"system_prompt"`
	AllowedTools []string `json:"allowed_tools"`
}

type Router struct {
	agents map[string]*Agent
}

var GlobalRouter = &Router{
	agents: make(map[string]*Agent),
}

func init() {
	// Register: System Agent
	GlobalRouter.Register(&Agent{
		Name:        "SystemAgent",
		Description: "Handles OS tasks, starting desktop applications, folder navigation, volume adjustments, and system power commands.",
		SystemPrompt: `You are the System Agent of the Jarvis platform. Your primary purpose is to help the user manage their computer system. 
You have direct access to tools for launching apps, navigating local user directories, adjusting sound volume, and managing power status.
Be concise and focus on executing the exact system request.`,
		AllowedTools: []string{"open_application", "open_folder", "control_volume", "power_control"},
	})

	// Register: Browser Agent
	GlobalRouter.Register(&Agent{
		Name:        "BrowserAgent",
		Description: "Handles web navigation, opening websites, performing Google/YouTube searches, and WhatsApp automation.",
		SystemPrompt: `You are the Browser Agent. Your role is to assist with web navigation, online search queries, YouTube video streaming lookups, and opening WhatsApp.
Formulate clear search intents. Use Google Search and Web navigation tools dynamically.`,
		AllowedTools: []string{"google_search", "youtube_search", "open_website", "whatsapp_action"},
	})

	// Register: Coding Agent
	GlobalRouter.Register(&Agent{
		Name:        "CodingAgent",
		Description: "Specialized in code generation, code review, debugging, project scaffolding, and architecture suggestions.",
		SystemPrompt: `You are the Coding Agent. Your task is to review, write, edit, and explain software source code.
Maintain high quality, use best practices, comment complex logic, and help scaffold modular structures.`,
		AllowedTools: []string{}, // Coding toolings will be added dynamically in sandbox
	})
}

// Register adds an agent to the router
func (r *Router) Register(a *Agent) {
	r.agents[a.Name] = a
	log.Printf("Agent registered: %s", a.Name)
}

// SelectAgent matches the user prompt to the correct agent using simple semantic rules or keyword routes
func (r *Router) Route(prompt string) (*Agent, error) {
	promptLower := strings.ToLower(prompt)

	// Routing logic for Coding
	codingKeywords := []string{"code", "program", "function", "debug", "write a script", "compile", "refactor", "bug", "scaffold", "golang", "javascript", "react", "html"}
	for _, kw := range codingKeywords {
		if strings.Contains(promptLower, kw) {
			return r.agents["CodingAgent"], nil
		}
	}

	// Routing logic for Browser
	browserKeywords := []string{"search", "google", "youtube", "website", "whatsapp", "call rahul", "chat with", "link", "browse", "internet", "open chrome"}
	for _, kw := range browserKeywords {
		if strings.Contains(promptLower, kw) {
			return r.agents["BrowserAgent"], nil
		}
	}

	// Default fallback to general SystemAgent
	return r.agents["SystemAgent"], nil
}

// ExecuteAgent query passes context and prompt execution to the selected agent
func ExecuteAgent(prompt string) (string, error) {
	agent, err := GlobalRouter.Route(prompt)
	if err != nil {
		return "", err
	}

	log.Printf("[AgentOrchestrator] Selected Specialized Agent: %s", agent.Name)

	// Prepare dynamic system instructions
	memContext := ai.GetRelevantMemories(prompt)
	fullSystemPrompt := agent.SystemPrompt
	if memContext != "" {
		fullSystemPrompt += "\n" + memContext
	}

	convID := ai.GetOrCreateConversation()
	history, err := ai.GetConversationHistory(convID, 8)
	if err != nil {
		log.Printf("Failed loading conversation history: %v", err)
	}

	// Assemble message payloads
	var messages []ai.Message
	messages = append(messages, ai.Message{Role: "system", Content: fullSystemPrompt})
	messages = append(messages, history...)
	messages = append(messages, ai.Message{Role: "user", Content: prompt})

	// Filter and wrap only allowed tools for this specific agent
	var ollamaTools []ai.OllamaToolWrapper
	for _, toolName := range agent.AllowedTools {
		rt, exists := tool.GlobalRegistry.Get(toolName)
		if !exists {
			continue
		}

		props := make(map[string]interface{})
		var required []string

		for pName, pDef := range rt.Definition.Parameters {
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

		ollamaTools = append(ollamaTools, ai.OllamaToolWrapper{
			Type: "function",
			Function: ai.OllamaToolFunction{
				Name:        rt.Definition.Name,
				Description: rt.Definition.Description,
				Parameters: ai.OllamaParamsDefinition{
					Type:       "object",
					Properties: props,
					Required:   required,
				},
			},
		})
	}

	reqBody := ai.OllamaChatRequest{
		Model:    ai.GetAvailableModel(),
		Messages: messages,
		Stream:   false,
		Tools:    ollamaTools,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	// Query Ollama Chat Endpoint
	respBytes, err := queryOllamaChat(jsonBody)
	if err != nil {
		return "", err
	}

	var chatResp ai.OllamaChatResponse
	if err := json.Unmarshal(respBytes, &chatResp); err != nil {
		return "", err
	}

	// If agent requested a tool call, invoke it
	if len(chatResp.Message.ToolCalls) > 0 {
		toolCall := chatResp.Message.ToolCalls[0]
		toolName := toolCall.Function.Name
		toolArgsMap := toolCall.Function.Arguments
		toolArgsJSON, _ := json.Marshal(toolArgsMap)

		log.Printf("[%s] Invoking tool: %s with args: %s", agent.Name, toolName, string(toolArgsJSON))

		// Check tool registration and HITL policy
		tObj, exists := tool.GlobalRegistry.Get(toolName)
		if exists && tObj.Definition.RequiresHITL {
			ai.LogToolExecution(agent.Name, toolName, toolArgsMap, "pending", nil, 0)
			return fmt.Sprintf("CONFIRMATION_REQUIRED:%s:%s", toolName, string(toolArgsJSON)), nil
		}

		startTime := time.Now()
		result, runErr := tool.GlobalRegistry.Invoke(toolName, string(toolArgsJSON))
		duration := time.Since(startTime).Milliseconds()

		status := "success"
		if runErr != nil {
			status = "failed"
		}
		
		ai.LogToolExecution(agent.Name, toolName, toolArgsMap, status, runErr, duration)

		if runErr != nil {
			return fmt.Sprintf("Error executing tool %s: %v", toolName, runErr), nil
		}

		ai.SaveInteraction(prompt, result)
		return result, nil
	}

	reply := chatResp.Message.Content
	ai.SaveInteraction(prompt, reply)
	return reply, nil
}

// Low-level helper to isolate HTTP query client logic
func queryOllamaChat(body []byte) ([]byte, error) {
	resp, err := http.Post("http://127.0.0.1:11434/api/chat", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("ollama API status %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}
