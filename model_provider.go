package summarizer

import "context"

type Provider interface {
	Name() string
	Summarize(context.Context, Request) (Response, error)
	// TODO tokenize
	// TODO detokenize
	// TODO ping
}

// Request
type Request struct {
}

// Response
type Response struct {
}

// openai

type OpenaiProvider struct {
	model  string
	apiKey string
}

func NewOpenaiProvider(model, apiKey string) OpenaiProvider {
	return OpenaiProvider{
		model:  model,
		apiKey: apiKey,
	}
}

func (p *OpenaiProvider) Name() string {
	return ""
}

func (p *OpenaiProvider) Summarize(ctx context.Context, req Request) (Response, error) {
	return Response{}, nil
}

// llama.cpp

type LlamacppProvider struct {
	url string
}

func NewLlamacppProvider(url string) LlamacppProvider {
	return LlamacppProvider{url: url}
}

func (p *LlamacppProvider) Name() string {
	return ""
}

func (p *LlamacppProvider) Summarize(ctx context.Context, req Request) (Response, error) {
	return Response{}, nil
}
