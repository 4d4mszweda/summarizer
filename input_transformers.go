package summarizer

import (
	"context"
	"strings"
)

// TODO lua script jako funkcja mutująca

type Transformer func(ctx context.Context, item InputItem, text string) (string, error)

func TrimSpace() Transformer {
	return func(ctx context.Context, item InputItem, text string) (string, error) {
		return strings.TrimSpace(text), nil
	}
}
