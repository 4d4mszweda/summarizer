// Package summarizer provides communication layer with AI to summarize texts
package summarizer

import (
	"context"
	"strings"
)

const defaultPromptTemplate = `You are a precise summarization assistant.
Summarize the following text in a concise, readable way.
Focus on the most important facts, claims, and outcomes.

Text:
{{text}}`

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
	provider       Provider
	promptTemplate string
}

type Provider interface {
	Name() string
	Summarize(context.Context, Request) (Response, error)
}

type Option func(*Service)

func WithPromptTemplate(prompt string) Option {
	return func(s *Service) {
		if strings.TrimSpace(prompt) != "" {
			s.promptTemplate = prompt
		}
	}
}

func New(provider Provider, opts ...Option) (*Service, error) {
	if provider == nil {
		return nil, ErrNoProvider
	}

	svc := &Service{
		provider:       provider,
		promptTemplate: defaultPromptTemplate,
	}

	for _, opt := range opts {
		opt(svc)
	}

	return svc, nil
}

func MustNew(provider Provider, opts ...Option) *Service {
	svc, err := New(provider, opts...)
	if err != nil {
		panic(err)
	}

	return svc
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
		Instructions: renderPrompt(s.promptTemplate, text),
	}

	for _, opt := range opts {
		opt(&req)
	}

	if strings.TrimSpace(req.Instructions) == "" {
		req.Instructions = renderPrompt(s.promptTemplate, text)
	}

	return s.provider.Summarize(ctx, req)
}

func DefaultPromptTemplate() string {
	return defaultPromptTemplate
}

func renderPrompt(tpl, text string) string {
	return strings.ReplaceAll(tpl, "{{text}}", text)
}
