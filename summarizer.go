// Package summarizer provides communication layer with AI models to summarize texts
package summarizer

import (
	"context"
)

// Request
type Request struct {
}

// Response
type Response struct {
}

// Service
type Service struct {
	provider              Provider
	prompt                Prompt
	input                 Input
	inputProccessingQueue []string // TODO
}

// Option
type Option func(*Service)

// New - constructor
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

// WithDefaultPromptTemplate -
func WithDefaultPromptTemplate() Option {
	return func(s *Service) {
	}
}

// AddItem - add item to be summarized
func (s *Service) AddItem(in InputItem) {}

// RunSummarize - run summarize proccess
func (s *Service) RunSummarize(ctx context.Context, text string, opts ...RequestOption) (Response, error) {
	return Response{}, nil
}
