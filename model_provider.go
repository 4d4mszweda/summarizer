package summarizer

import "context"

type Provider interface {
	Name() string
	Summarize(context.Context, Request) (Response, error)
}

// TODO tutaj trzeba gotowych providerów lub innych construktorów dla danych modeli

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

// llama.cpp

type LlamacppProvider struct{}

func NewLlamacppProvider(model, apiKey string) LlamacppProvider {
	return LlamacppProvider{}
}

// Cloude

//type CloudeProvider struct {
//	model  string
//	apiKey string
//}
//
//func NewCloudeProvider(model, apiKey string) CloudeProvider {
//	return CloudeProvider{
//		model:  model,
//		apiKey: apiKey,
//	}
//}
