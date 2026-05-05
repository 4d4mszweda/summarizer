package summarizer

import (
	"context"
	"io"
)

type Input struct {
	transformers []Transformer
	inputQueue   []InputItem
}

func (i *Input) addItem(item InputItem) {
	i.inputQueue = append(i.inputQueue, item)
}

func (i *Input) clearItems() {
	i.inputQueue = i.inputQueue[:0] // without moving cap
}

func (i *Input) use(transformer Transformer) {
	i.transformers = append(i.transformers, transformer)
}

type Chunk struct {
	ItemName string
	Index    int
	Text     string
}

func (i *Input) process(ctx context.Context, chunkSize int) ([]Chunk, error) {
	var chunks []Chunk

	for _, item := range i.inputQueue {

		text, err := readInputItem(ctx, item)
		if err != nil {
			return nil, err
		}

		for _, transformer := range i.transformers {
			text, err = transformer(ctx, item, text)
			if err != nil {
				return nil, err
			}
		}

		itemChunks := chunkText(item.Name(), text, chunkSize)
		chunks = append(chunks, itemChunks...)
	}

	return chunks, nil
}

// TODO SPORO KOIPIOWANIA DANYCH

// TODO ReadAll nie jest raczej dobre dla dużych tekstów
func readInputItem(ctx context.Context, item InputItem) (string, error) {
	r, err := item.Open(ctx)
	if err != nil {
		return "", err
	}
	defer r.Close()

	b, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func chunkText(itemName string, text string, size int) []Chunk {
	if size <= 0 {
		return []Chunk{
			{
				ItemName: itemName,
				Index:    0,
				Text:     text,
			},
		}
	}

	runes := []rune(text)

	var chunks []Chunk
	for start, index := 0, 0; start < len(runes); start, index = start+size, index+1 {
		end := start + size
		if end > len(runes) {
			end = len(runes)
		}

		chunks = append(chunks, Chunk{
			ItemName: itemName,
			Index:    index,
			Text:     string(runes[start:end]),
		})
	}

	return chunks
}
