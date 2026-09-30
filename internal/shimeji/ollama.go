package shimeji

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	config "shimeji-pet-companion/internal/config"
)

type OllamaResponse struct {
	SubtitleEn string `json:"subtitle_en"`
	Animation  string `json:"animation"`
}

var httpClient = &http.Client{
	Timeout: 60 * time.Second, // Allow local LLMs time to complete generation
}

func (s *ShimejiService) AskOllama(prompt string) (*OllamaResponse, error) {
	builtPrompt, err := s.promptGenerator.BuildPrompt(config.PromptContext{
		PetName:   s.cfg.Pet.Name,
		Character: s.cfg.Pet.Character,
		UserEvent: prompt,
	})
	if err != nil {
		return nil, fmt.Errorf("prompt build error: %w", err)
	}

	payload := map[string]any{
		"model":  s.cfg.Ollama.Model,
		"prompt": builtPrompt,
		"format": "json",
		"stream": s.cfg.Ollama.Stream,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("payload marshal error: %w", err)
	}

	// 1. Ensure the URL points to Ollama's /api/generate
	baseURL := strings.TrimRight(s.cfg.Ollama.Endpoint, "/")
	url := fmt.Sprintf("%s/api/generate", baseURL)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// 2. Actually execute the HTTP request over the network
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ollama connection failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respErrBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ollama returned error status %d: %s", resp.StatusCode, string(respErrBytes))
	}

	// 3. Read the actual Ollama response body
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var outer struct {
		Response string `json:"response"`
	}
	if err := json.Unmarshal(respBytes, &outer); err != nil {
		return nil, fmt.Errorf("failed to parse outer json: %w", err)
	}

	// 4. Unmarshal the generated JSON string into OllamaResponse
	var ollamaResp OllamaResponse
	if err := json.Unmarshal([]byte(outer.Response), &ollamaResp); err != nil {
		return nil, fmt.Errorf("failed to parse pet response JSON (%s): %w", outer.Response, err)
	}

	return &ollamaResp, nil
}