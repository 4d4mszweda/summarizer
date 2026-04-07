# summary

`summary` is a Go library for AI text summarization with pluggable LLM providers and a Bubble Tea TUI.

## Goals

- import the library into your own Go projects
- switch between LLM distributors behind a shared provider interface
- run a local terminal app for manual summarization

## Packages

- `summary` - public API for summarization services
- `providers/openai` - OpenAI-compatible provider that can target different distributors through `BaseURL`
- `providers/mock` - local fallback provider for development and tests
- `app` - provider wiring used by the TUI
- `cmd/summary` - Bubble Tea terminal app

## Environment

The TUI reads:

- `SUMMARY_PROVIDER` - `mock`, `openai`, `openrouter`, `groq`, `together`, `fireworks`, or `ollama`
- `SUMMARY_API_KEY` - API key for the selected distributor
- `SUMMARY_BASE_URL` - optional explicit OpenAI-compatible base URL override
- `SUMMARY_MODEL` - model name, for example `gpt-4.1-mini`

## Example library usage

```go
package main

import (
	"context"
	"fmt"

	"summary"
	"summary/providers/mock"
)

func main() {
	service := summary.MustNew(mock.New())

	res, err := service.Summarize(context.Background(), "Long text to summarize")
	if err != nil {
		panic(err)
	}

	fmt.Println(res.Summary)
}
```

## Run TUI

```bash
go run ./cmd/summary
```
