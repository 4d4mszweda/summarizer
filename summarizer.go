// Package summarizer provides communication layer with AI models to summarize texts
package summarizer

import (
	"context"
)

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

///////////////////////////////////////////////////////////////////////////////////////////////////////////
/////////////////////////////////////// CONSTRUCTOR OPTIONS ///////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////

// WithDefaultPromptTemplate -
func WithDefaultPromptTemplate() Option {
	return func(s *Service) {
	}
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////
/////////////////////////////////////// SERVICE METHODS ///////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////

// AddItem - add item to be summarized
func (s *Service) AddItem(in InputItem) {}

// RunSummarize - run summarize proccess
func (s *Service) RunSummarize(ctx context.Context, opts ...RequestOption) (Response, error) {
	return Response{}, nil
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////// REQUEST OPTIONS ///////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////

type RequestOption func(*Request)

func WithInstructions(instructions string) RequestOption {
	return func(r *Request) {
	}
}

func WithModel(model string) RequestOption {
	return func(r *Request) {
	}
}

func WithTemperature(temperature float64) RequestOption {
	return func(r *Request) {
	}
}

func WithMaxTokens(maxTokens int) RequestOption {
	return func(r *Request) {
	}
}
