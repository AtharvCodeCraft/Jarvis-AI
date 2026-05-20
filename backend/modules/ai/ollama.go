package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type OllamaResponse struct {
	Response string `json:"response"`
}

func QueryOllama(prompt string) (string, error) {
	promptContext := GetContext() + "\nUser: " + prompt
	log.Println("Querying local LLM via Ollama...")

	reqBody := OllamaRequest{
		Model:  "qwen2.5:1.5b", // Running lightweight Qwen 1.5b instead of Llama 3 for 8GB RAM limit
		Prompt: promptContext,
		Stream: false,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	// Make HTTP POST request to default local Ollama instance
	resp, err := http.Post("http://127.0.0.1:11434/api/generate", "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("ollama API returned status: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	var ollamaResp OllamaResponse
	if err := json.Unmarshal(body, &ollamaResp); err != nil {
		return "", err
	}

	SaveInteraction(prompt, ollamaResp.Response)

	return ollamaResp.Response, nil
}
