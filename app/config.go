package app

import (
	"fmt"
	"os"
	"strings"

	"summary"
	"summary/providers/mock"
	"summary/providers/openai"
)

const (
	ProviderMock       = "mock"
	ProviderOpenAI     = "openai"
	ProviderOpenRouter = "openrouter"
	ProviderGroq       = "groq"
	ProviderTogether   = "together"
	ProviderFireworks  = "fireworks"
	ProviderOllama     = "ollama"
)

type Config struct {
	Provider string
	BaseURL  string
	APIKey   string
	Model    string
}

func ConfigFromEnv() Config {
	cfg := Config{
		Provider: strings.TrimSpace(os.Getenv("SUMMARY_PROVIDER")),
		BaseURL:  strings.TrimSpace(os.Getenv("SUMMARY_BASE_URL")),
		APIKey:   strings.TrimSpace(os.Getenv("SUMMARY_API_KEY")),
		Model:    strings.TrimSpace(os.Getenv("SUMMARY_MODEL")),
	}

	if cfg.Provider == "" {
		cfg.Provider = ProviderMock
	}
	if cfg.Model == "" {
		cfg.Model = "gpt-4.1-mini"
	}

	return cfg
}

func NewProvider(cfg Config) (summary.Provider, error) {
	name := strings.ToLower(strings.TrimSpace(cfg.Provider))
	baseURL := strings.TrimSpace(cfg.BaseURL)
	if baseURL == "" {
		baseURL = defaultBaseURLFor(name)
	}

	switch name {
	case ProviderMock:
		return mock.New(), nil
	case ProviderOpenAI, ProviderOpenRouter, ProviderGroq, ProviderTogether, ProviderFireworks, ProviderOllama:
		return openai.New(openai.Config{
			Name:    name,
			BaseURL: baseURL,
			APIKey:  cfg.APIKey,
			Model:   cfg.Model,
		})
	default:
		return nil, fmt.Errorf("unsupported provider %q", cfg.Provider)
	}
}

func defaultBaseURLFor(provider string) string {
	switch provider {
	case ProviderOpenRouter:
		return "https://openrouter.ai/api/v1"
	case ProviderGroq:
		return "https://api.groq.com/openai/v1"
	case ProviderTogether:
		return "https://api.together.xyz/v1"
	case ProviderFireworks:
		return "https://api.fireworks.ai/inference/v1"
	case ProviderOllama:
		return "http://localhost:11434/v1"
	default:
		return ""
	}
}
