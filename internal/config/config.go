package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Ollama struct {
		Endpoint 	string 	`yaml:"endpoint"`
		Model    	string 	`yaml:"model"`
		Stream   	bool   	`yaml:"stream"`
	} `yaml:"ollama"`

	Pet struct {
		Name             	string `yaml:"name"`
		Character         	string `yaml:"character"`
	} `yaml:"pet"`

	Window struct {
		StartClickThrough bool `yaml:"start_click_through"`
	} `yaml:"window"`
}

// LoadConfig reads config.yaml, or generates sensible defaults if not found
func LoadConfig(path string) *Config {
	cfg := &Config{}

	// Defaults in case the file doesn't exist
	cfg.Ollama.Endpoint = "http://127.0.0.1:11434"
	cfg.Ollama.Model = "hf.co/DreamFast/Qwen3-VL-4b-Heretic-GGUF:Q4_K_M"
	cfg.Pet.Character = "cute anime girl"
	cfg.Window.StartClickThrough = true

	data, err := os.ReadFile(path)
	if err != nil {
		log.Printf("config.yaml not found or unreadable, using default settings: %v", err)
		return cfg
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		log.Printf("Error parsing config.yaml, using defaults: %v", err)
		return cfg
	}

	return cfg
}