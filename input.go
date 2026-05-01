package summarizer

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Input struct {
}

type InputItem interface {
	Name() string
	Open(ctx context.Context) (io.ReadCloser, error)
}

// STRING INPUT

type StringItem struct {
	name string
	text string
}

func NewStringItem(name, text string) *StringItem {
	return &StringItem{name: name, text: text}
}

func (i *StringItem) Name() string {
	return i.name
}

func (i *StringItem) Open(ctx context.Context) (io.ReadCloser, error) {
	return io.NopCloser(strings.NewReader(i.text)), nil
}

// FILE INPUT

type FileItem struct {
	path string
}

func NewFileItem(path string) *FileItem {
	return &FileItem{path: path}
}

func (i *FileItem) Name() string {
	return filepath.Base(i.path)
}

func (i *FileItem) Open(ctx context.Context) (io.ReadCloser, error) {
	return os.Open(i.path)
}
