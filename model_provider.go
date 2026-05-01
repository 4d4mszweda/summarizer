package summarizer

import "context"

type Provider interface {
	Name() string
	Summarize(context.Context, Request) (Response, error)
}

// TODO tutaj trzeba gotowych providerów lub innych construktorów dla danych modeli

// openai

type OpenaiProvider struct{}

// Cloude

type CloudeProvider struct{}

// llama.cpp

type LlmacppProvider struct{}
