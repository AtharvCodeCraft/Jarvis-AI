package system

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Define local structs to avoid import cycle with ai module
type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type OllamaResponse struct {
	Response string `json:"response"`
}

// WriteCode generates code via Ollama, saves it to a file, and opens it in VS Code
func WriteCode(prompt string) string {
	log.Printf("Generating code for prompt: %s", prompt)

	// Ensure we are asking for ONLY code
	systemPrompt := `You are an expert programmer. Write the exact code the user requests. 
DO NOT include any markdown formatting like \` + "\x60" + `\x60\x60python\x60\x60` + "\x60" + `. 
DO NOT include any explanations or conversational text. 
ONLY output the raw, executable code.`

	fullPrompt := systemPrompt + "\nUser: " + prompt

	reqBody := OllamaRequest{
		Model:  "qwen2.5:1.5b",
		Prompt: fullPrompt,
		Stream: false,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		log.Printf("Error marshalling code prompt: %v", err)
		return "Failed to generate code."
	}

	resp, err := http.Post("http://127.0.0.1:11434/api/generate", "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Printf("Error requesting code from Ollama: %v", err)
		return "Failed to generate code. Is Ollama running?"
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "Ollama returned an error generating your code."
	}

	body, err := ioutil.ReadAll(resp.Body)
	var ollamaResp OllamaResponse
	if err := json.Unmarshal(body, &ollamaResp); err != nil {
		return "Failed to parse the generated code."
	}

	rawCode := strings.TrimSpace(ollamaResp.Response)
	// Strip markdown blocks if the model ignored our instructions
	if strings.HasPrefix(rawCode, "```") {
		lines := strings.Split(rawCode, "\n")
		if len(lines) > 2 {
			rawCode = strings.Join(lines[1:len(lines)-1], "\n")
		}
	}

	// Determine file extension
	ext := ".txt"
	promptLower := strings.ToLower(prompt)
	if strings.Contains(promptLower, "python") {
		ext = ".py"
	} else if strings.Contains(promptLower, "javascript") || strings.Contains(promptLower, " js ") {
		ext = ".js"
	} else if strings.Contains(promptLower, "html") {
		ext = ".html"
	} else if strings.Contains(promptLower, "go ") || strings.Contains(promptLower, "golang") {
		ext = ".go"
	} else if strings.Contains(promptLower, "java ") {
		ext = ".java"
	} else if strings.Contains(promptLower, "c++") || strings.Contains(promptLower, "cpp") {
		ext = ".cpp"
	} else if strings.Contains(promptLower, "c#") || strings.Contains(promptLower, "csharp") {
		ext = ".cs"
	} else if strings.Contains(promptLower, "php") {
		ext = ".php"
	} else if strings.Contains(promptLower, "css") {
		ext = ".css"
	}

	// Create JarvisCode folder on Desktop
	// Create JarvisCode folder on Desktop or in Temp as fallback
	userProfile := os.Getenv("USERPROFILE")
	folderPath := filepath.Join(userProfile, "Desktop", "JarvisCode")

	if userProfile == "" {
		folderPath = filepath.Join(os.TempDir(), "JarvisCode")
	}

	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		err := os.MkdirAll(folderPath, 0755)
		if err != nil {
			log.Printf("Failed to create code folder: %v", err)
			return "I couldn't create the folder to save your code."
		}
	}

	// Save the file
	filename := filepath.Join(folderPath, "generated_script"+ext)
	err = ioutil.WriteFile(filename, []byte(rawCode), 0644)
	if err != nil {
		log.Printf("Failed to write code to file: %v", err)
		return "I failed to save the generated code."
	}

	// Open the file in VS Code
	cmd := exec.Command("cmd", "/c", "code", filename)
	startHidden(cmd)

	return "I wrote the code and opened it in VS Code."
}
