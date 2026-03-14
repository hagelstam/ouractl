package tui

import (
	"fmt"

	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// DetailModel is a bubbletea model that shows a spinner while fetching,
// then renders the detail content and quits.
type DetailModel struct {
	spinner spinner.Model
	fetch   tea.Cmd
	content string
	done    bool
}

// NewDetailModel creates a model that shows a spinner while fetch runs.
func NewDetailModel(fetch tea.Cmd) DetailModel {
	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(Accent)

	return DetailModel{
		spinner: sp,
		fetch:   fetch,
	}
}

func (m DetailModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, m.fetch)
}

func (m DetailModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case DetailData:
		if msg.Err != nil {
			m.content = lipgloss.NewStyle().Foreground(Bad).
				Render("Error: " + msg.Err.Error())
		} else {
			m.content = msg.Content
		}
		m.done = true
		return m, tea.Quit
	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	case spinner.TickMsg:
		if !m.done {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
	}
	return m, nil
}

func (m DetailModel) View() tea.View {
	if m.done {
		return tea.NewView("\n" + m.content + "\n")
	}
	return tea.NewView(fmt.Sprintf("\n  %s Loading...\n", m.spinner.View()))
}
