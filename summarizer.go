// Package summarizer provides communication layer with AI models to summarize texts
package summarizer

import (
	"context"
)

// Service
type Service struct {
	provider   Provider
	prompt     Prompt
	input      Input
	prevResult Response
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
/////////////////////////////////////// SERVICE OPTIONS ///////////////////////////////////////////////////
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
func (s *Service) AddItem(item InputItem) *Service { s.input.addItem(item); return s }

// ClearItems - clean items queue
func (s *Service) ClearItems() *Service { s.input.clearItems(); return s }

// AddPrevResult - add previous summarization as a context for new
func (s *Service) AddPrevResult(in InputItem) *Service { return s }

// UseTransformer
func (s *Service) UseTransformer(trans Transformer) *Service { s.input.use(trans); return s }

// RunSummarize - run summarize proccess
func (s *Service) RunSummarize(ctx context.Context, opts ...RequestOption) (Response, error) {
	req := Request{}
	for _, opt := range opts {
		opt(&req)
	}

	// TODO jak narazie rozpatruje:
	// * tylko stringi
	// * tylko openai + llama.cpp

	// 1. process input && batch input
	chunks, err := s.input.process(ctx, 1000)
	if err != nil {
		return Response{}, err
	}

	lenChunks := len(chunks)

	// 2. build prompt
	for index, chunk := range chunks {

		// 2.1. add prevResult to prompt
		// 2.2. add chunk to prompt

		// 3. Apply options

		// 4. Send Request
	}

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

func WithTemperature(temperature float64) RequestOption {
	return func(r *Request) {
	}
}

func WithMaxTokens(maxTokens int) RequestOption {
	return func(r *Request) {
	}
}
