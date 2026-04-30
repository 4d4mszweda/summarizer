package summarizer

type RequestOption func(*Request)

func WithInstructions(instructions string) RequestOption {
	return func(r *Request) {
		r.Instructions = instructions
	}
}

func WithModel(model string) RequestOption {
	return func(r *Request) {
		r.Model = model
	}
}

func WithTemperature(temperature float64) RequestOption {
	return func(r *Request) {
		r.Temperature = temperature
	}
}

func WithMaxTokens(maxTokens int) RequestOption {
	return func(r *Request) {
		r.MaxTokens = maxTokens
	}
}
