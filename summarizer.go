// Package summarizer provides communication layer with AI models to summarize texts
package summarizer

import (
	"context"
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

func (s *Service) AddContext(in InputItem) {}

func (s *Service) AddItem(in InputItem) {}

func (s *Service) Summarize(ctx context.Context, text string, opts ...RequestOption) (Response, error) {
	return Response{}, nil
}
