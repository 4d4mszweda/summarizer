// Package summarizer provides communication layer with AI models to summarize texts
package summarizer

import (
	"context"
	"strings"
)

type Request struct {
	Text         string
	Instructions string
	Model        string
	Temperature  float64
	MaxTokens    int
}

type Response struct {
	Summary  string
	Provider string
	Model    string
	Raw      string
}

type Service struct {
	provider Provider
	prompt   Prompt
	input    Input
}

type Option func(*Service)

func New(provider Provider, opts ...Option) (*Service, error) {
	if provider == nil {
		return nil, ErrNoProvider
	}

	svc := &Service{
		provider: provider,
	}

	for _, opt := range opts {
		opt(svc)
	}

	return svc, nil
}

func WithDefaultPromptTemplate() Option {
	return func(s *Service) {

	}
}

func (s *Service) Summarize(ctx context.Context, text string, opts ...RequestOption) (Response, error) {
	if s == nil || s.provider == nil {
		return Response{}, ErrNoProvider
	}

	text = strings.TrimSpace(text)
	if text == "" {
		return Response{}, ErrEmptyText
	}

	req := Request{
		Text:         text,
		Instructions: renderPrompt(s.prompt, text),
	}

	for _, opt := range opts {
		opt(&req)
	}

	if strings.TrimSpace(req.Instructions) == "" {
		req.Instructions = renderPrompt(s.prompt, text)
	}

	return s.provider.Summarize(ctx, req)
}
