package summarizer

import (
	"bytes"
	"text/template"
)

const defaultPromptTemplate = `You are a precise summarization assistant.

{{- if .Context }}
CONTEXT:
{{ .Context }}
{{- end }}

{{- if .Instructions }}
INSTRUCTIONS:
{{ .Instructions }}
{{- end }}

{{- if .Output }}
OUTPUT:
{{ .Output }}
{{- end }}

{{- if .Intput }}
INTPUT:
{{ .Intput }}
{{- end }}
`

type promptData struct {
	Context      string
	Instructions string
	Output       string
	Input        string
}

type Prompt struct {
	tmpl *template.Template
	data promptData
}

func DefaultPrompt() (*Prompt, error) {
	tmpl, err := template.New("prompt").Parse(defaultPromptTemplate)
	if err != nil {
		return nil, err
	}

	return &Prompt{tmpl: tmpl}, nil
}

func (p *Prompt) Build() (string, error) {
	var buf bytes.Buffer

	err := p.tmpl.Execute(&buf, p.data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
