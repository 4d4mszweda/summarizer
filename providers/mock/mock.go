package mock

import (
	"context"
	"strings"

	"summary"
)

type Provider struct{}

func New() *Provider {
	return &Provider{}
}

func (p *Provider) Name() string {
	return "mock"
}

func (p *Provider) Summarize(_ context.Context, req summary.Request) (summary.Response, error) {
	text := strings.TrimSpace(req.Text)
	words := strings.Fields(text)
	if len(words) > 32 {
		words = words[:32]
	}

	result := strings.Join(words, " ")
	if len(strings.Fields(text)) > len(words) {
		result += "..."
	}

	return summary.Response{
		Summary:  result,
		Provider: p.Name(),
		Model:    "mock-preview",
		Raw:      result,
	}, nil
}
