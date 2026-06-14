package tool

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
)

// ToolParam defines a schema parameter for a tool argument
type ToolParam struct {
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Required    bool     `json:"required,omitempty"`
	Enum        []string `json:"enum,omitempty"`
}

// ToolDefinition defines the structure passed to LLMs for tool selection
type ToolDefinition struct {
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Parameters  map[string]ToolParam `json:"parameters"`
	RequiresHITL bool                `json:"requires_hitl"` // Requires Human-In-The-Loop confirmation
}

// Handler defines the function to execute for a tool
type Handler func(args map[string]interface{}) (string, error)

// Tool links definition to its actual backend execution handler
type Tool struct {
	Definition ToolDefinition
	Execute    Handler
}

// Registry manages the set of available tools
type Registry struct {
	mu    sync.RWMutex
	tools map[string]*Tool
}

var GlobalRegistry = &Registry{
	tools: make(map[string]*Tool),
}

// Register adds a new tool to the registry
func (r *Registry) Register(name string, desc string, params map[string]ToolParam, requiresHitl bool, handler Handler) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.tools[name] = &Tool{
		Definition: ToolDefinition{
			Name:         name,
			Description:  desc,
			Parameters:   params,
			RequiresHITL: requiresHitl,
		},
		Execute: handler,
	}
	log.Printf("Tool registered successfully: %s", name)
}

// Get retrieves a tool by name
func (r *Registry) Get(name string) (*Tool, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	t, exists := r.tools[name]
	return t, exists
}

// ListDefinitions returns JSON-ready tool definitions for the LLM
func (r *Registry) ListDefinitions() []ToolDefinition {
	r.mu.RLock()
	defer r.mu.RUnlock()

	defs := make([]ToolDefinition, 0, len(r.tools))
	for _, t := range r.tools {
		defs = append(defs, t.Definition)
	}
	return defs
}

// Invoke executes a registered tool with the provided raw JSON arguments
func (r *Registry) Invoke(name string, jsonArgs string) (string, error) {
	t, exists := r.Get(name)
	if !exists {
		return "", fmt.Errorf("tool not found: %s", name)
	}

	var args map[string]interface{}
	if err := json.Unmarshal([]byte(jsonArgs), &args); err != nil {
		return "", fmt.Errorf("failed to parse arguments for tool %s: %w", name, err)
	}

	// Dynamic validation can be implemented here based on t.Definition.Parameters

	log.Printf("Invoking tool [%s] with args: %s", name, jsonArgs)
	return t.Execute(args)
}
