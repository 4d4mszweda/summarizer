package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"summary"
	"summary/app"
)

type summaryResultMsg struct {
	response summary.Response
	err      error
}

type model struct {
	cfg           app.Config
	service       *summary.Service
	providerInfo  string
	input         textarea.Model
	providerInput textinput.Model
	modelInput    textinput.Model
	baseURLInput  textinput.Model
	output        viewport.Model
	spinner       spinner.Model
	activeField   int
	loading       bool
	err           string
	width         int
	height        int
}

func main() {
	cfg := app.ConfigFromEnv()
	provider, err := app.NewProvider(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "provider setup failed: %v\n", err)
		os.Exit(1)
	}

	service, err := summary.New(provider)
	if err != nil {
		fmt.Fprintf(os.Stderr, "service setup failed: %v\n", err)
		os.Exit(1)
	}

	p := tea.NewProgram(newModel(cfg, service), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "tui failed: %v\n", err)
		os.Exit(1)
	}
}

func newModel(cfg app.Config, service *summary.Service) model {
	input := textarea.New()
	input.Placeholder = "Wklej tekst do podsumowania..."
	input.SetWidth(80)
	input.SetHeight(12)
	input.Focus()
	input.Prompt = "│ "
	input.CharLimit = 0

	modelInput := textinput.New()
	modelInput.Placeholder = "Model"
	modelInput.SetValue(cfg.Model)
	modelInput.Prompt = "Model: "

	providerInput := textinput.New()
	providerInput.Placeholder = "Provider"
	providerInput.SetValue(cfg.Provider)
	providerInput.Prompt = "Provider: "

	baseURLInput := textinput.New()
	baseURLInput.Placeholder = "Base URL"
	baseURLInput.SetValue(cfg.BaseURL)
	baseURLInput.Prompt = "Base URL: "

	spin := spinner.New()
	spin.Spinner = spinner.Dot

	out := viewport.New(80, 10)
	out.SetContent("Podsumowanie pojawi się tutaj.")

	return model{
		cfg:           cfg,
		service:       service,
		providerInfo:  providerLabel(cfg),
		input:         input,
		providerInput: providerInput,
		modelInput:    modelInput,
		baseURLInput:  baseURLInput,
		output:        out,
		spinner:       spin,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.input.SetWidth(max(30, msg.Width-6))
		m.output.Width = max(30, msg.Width-6)
		m.output.Height = max(8, msg.Height-18)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "tab", "shift+tab":
			m.cycleFocus(msg.String() == "shift+tab")
			return m, nil
		case "ctrl+s":
			if m.loading {
				return m, nil
			}
			m.loading = true
			m.err = ""
			m.output.SetContent("Generowanie podsumowania...")
			return m, tea.Batch(m.spinner.Tick, m.summarizeCmd())
		}

	case summaryResultMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err.Error()
			m.output.SetContent("Brak wyniku.")
			return m, nil
		}

		m.err = ""
		m.output.SetContent(msg.response.Summary)
		m.providerInfo = fmt.Sprintf("%s | %s", msg.response.Provider, msg.response.Model)
		return m, nil

	case spinner.TickMsg:
		if !m.loading {
			return m, nil
		}
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch m.activeField {
	case 0:
		m.input, cmd = m.input.Update(msg)
	case 1:
		m.providerInput, cmd = m.providerInput.Update(msg)
	case 2:
		m.modelInput, cmd = m.modelInput.Update(msg)
	case 3:
		m.baseURLInput, cmd = m.baseURLInput.Update(msg)
	}
	cmds = append(cmds, cmd)

	m.output, cmd = m.output.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	boxStyle := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(0, 1)
	helpStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("9"))

	status := m.providerInfo
	if m.loading {
		status = fmt.Sprintf("%s Generowanie przez %s", m.spinner.View(), m.providerInfo)
	}

	var b strings.Builder
	b.WriteString(titleStyle.Render("Summary TUI"))
	b.WriteString("\n")
	b.WriteString(helpStyle.Render("tab: zmiana pola | ctrl+s: podsumuj | esc: wyjdz"))
	b.WriteString("\n\n")
	b.WriteString(boxStyle.Render(status))
	b.WriteString("\n\n")
	b.WriteString(m.providerInput.View())
	b.WriteString("\n")
	b.WriteString(m.modelInput.View())
	b.WriteString("\n")
	b.WriteString(m.baseURLInput.View())
	b.WriteString("\n\n")
	b.WriteString(boxStyle.Render(m.input.View()))
	b.WriteString("\n\n")
	b.WriteString(boxStyle.Render(m.output.View()))

	if m.err != "" {
		b.WriteString("\n\n")
		b.WriteString(errorStyle.Render("Blad: " + m.err))
	}

	return b.String()
}

func (m *model) cycleFocus(reverse bool) {
	fields := []*textinput.Model{&m.providerInput, &m.modelInput, &m.baseURLInput}

	m.input.Blur()
	for _, field := range fields {
		field.Blur()
	}

	if reverse {
		m.activeField--
	} else {
		m.activeField++
	}

	if m.activeField < 0 {
		m.activeField = 3
	}
	if m.activeField > 3 {
		m.activeField = 0
	}

	switch m.activeField {
	case 0:
		m.input.Focus()
	case 1:
		m.providerInput.Focus()
	case 2:
		m.modelInput.Focus()
	case 3:
		m.baseURLInput.Focus()
	}
}

func (m model) summarizeCmd() tea.Cmd {
	text := m.input.Value()
	providerName := strings.TrimSpace(m.providerInput.Value())
	modelName := strings.TrimSpace(m.modelInput.Value())
	baseURL := strings.TrimSpace(m.baseURLInput.Value())
	cfg := m.cfg
	cfg.Provider = providerName
	cfg.Model = modelName
	cfg.BaseURL = baseURL

	return func() tea.Msg {
		provider, err := app.NewProvider(cfg)
		if err != nil {
			return summaryResultMsg{err: err}
		}

		service, err := summary.New(provider)
		if err != nil {
			return summaryResultMsg{err: err}
		}

		ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
		defer cancel()

		res, err := service.Summarize(
			ctx,
			text,
			summary.WithModel(modelName),
			summary.WithTemperature(0.2),
			summary.WithMaxTokens(400),
		)
		return summaryResultMsg{response: res, err: err}
	}
}

func providerLabel(cfg app.Config) string {
	modelName := cfg.Model
	if modelName == "" {
		modelName = "-"
	}
	providerName := cfg.Provider
	if providerName == "" {
		providerName = "-"
	}
	return fmt.Sprintf("%s | %s", providerName, modelName)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
