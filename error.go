package summarizer

import "errors"

var (
	ErrNoProvider        = errors.New("summary: provider is required")
	ErrEmptyText         = errors.New("summary: text is required")
	ErrProviderUnvaiable = errors.New("summary: provider is unvaiable")
)
