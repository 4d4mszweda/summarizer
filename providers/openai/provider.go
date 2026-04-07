package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"summary"
)

const defaultBaseURL = "https://api.openai.com/v1"

var (
	ErrMissingModel = errors.New("openai provider: missing model")
)

type Config struct {
	Name       string
	BaseURL    string
	APIKey     string
	Model      string
	HTTPClient *http.Client
}

type Provider struct {
	name       string
	baseURL    string
	apiKey     string
	model      string
	httpClient *http.Client
}

func New(cfg Config) (*Provider, error) {
	if strings.TrimSpace(cfg.Model) == "" {
		return nil, ErrMissingModel
	}

	baseURL := strings.TrimRight(strings.TrimSpace(cfg.BaseURL), "/")
	if baseURL == "" {
		baseURL = defaultBaseURL
	}

	httpClient := cfg.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 60 * time.Second}
	}

	name := strings.TrimSpace(cfg.Name)
	if name == "" {
		name = "openai-compatible"
	}

	return &Provider{
		name:       name,
		baseURL:    baseURL,
		apiKey:     cfg.APIKey,
		model:      cfg.Model,
		httpClient: httpClient,
	}, nil
}

func (p *Provider) Name() string {
	return p.name
}

type chatCompletionRequest struct {
	Model       string        `json:"model"`
	Messages    []chatMessage `json:"messages"`
	Temperature *float64      `json:"temperature,omitempty"`
	MaxTokens   int           `json:"max_tokens,omitempty"`
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatCompletionResponse struct {
	Choices []struct {
		Message chatMessage `json:"message"`
	} `json:"choices"`
}

func (p *Provider) Summarize(ctx context.Context, req summary.Request) (summary.Response, error) {
	model := strings.TrimSpace(req.Model)
	if model == "" {
		model = p.model
	}

	payload := chatCompletionRequest{
		Model: model,
		Messages: []chatMessage{
			{
				Role:    "system",
				Content: "You summarize text precisely and without fluff.",
			},
			{
				Role:    "user",
				Content: req.Instructions,
			},
		},
		MaxTokens: req.MaxTokens,
	}

	if req.Temperature != 0 {
		payload.Temperature = &req.Temperature
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return summary.Response{}, fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		p.baseURL+"/chat/completions",
		bytes.NewReader(body),
	)
	if err != nil {
		return summary.Response{}, fmt.Errorf("build request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	if strings.TrimSpace(p.apiKey) != "" {
		httpReq.Header.Set("Authorization", "Bearer "+p.apiKey)
	}

	httpRes, err := p.httpClient.Do(httpReq)
	if err != nil {
		return summary.Response{}, fmt.Errorf("send request: %w", err)
	}
	defer httpRes.Body.Close()

	rawBody, err := io.ReadAll(httpRes.Body)
	if err != nil {
		return summary.Response{}, fmt.Errorf("read response: %w", err)
	}

	if httpRes.StatusCode >= 400 {
		return summary.Response{}, fmt.Errorf("provider returned %s: %s", httpRes.Status, strings.TrimSpace(string(rawBody)))
	}

	var payloadRes chatCompletionResponse
	if err := json.Unmarshal(rawBody, &payloadRes); err != nil {
		return summary.Response{}, fmt.Errorf("decode response: %w", err)
	}

	if len(payloadRes.Choices) == 0 {
		return summary.Response{}, errors.New("provider returned no choices")
	}

	content := strings.TrimSpace(payloadRes.Choices[0].Message.Content)

	return summary.Response{
		Summary:  content,
		Provider: p.Name(),
		Model:    model,
		Raw:      string(rawBody),
	}, nil
}
