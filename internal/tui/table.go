package tui

import (
	"fmt"

	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type viewState int

const (
	stateTable viewState = iota
	stateLoading
	stateDetail
)

// DetailData is the result returned from a detail fetch.
type DetailData struct {
	Content string
	Err     error
}

// TableConfig configures the generic table-detail model.
type TableConfig struct {
	Columns     []table.Column
	Rows        []table.Row
	Width       int
	FetchDetail func(row table.Row) tea.Cmd
}

// TableModel is a bubbletea model with table => loading => detail states.
type TableModel struct {
	config  TableConfig
	table   table.Model
	spinner spinner.Model
	state   viewState
	detail  string
}

// NewTableModel creates a new table-detail model from the given config.
func NewTableModel(cfg TableConfig) TableModel {
	t := table.New(
		table.WithColumns(cfg.Columns),
		table.WithRows(cfg.Rows),
		table.WithFocused(true),
		table.WithHeight(len(cfg.Rows)),
		table.WithWidth(cfg.Width),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(Subtle).
		BorderBottom(true).
		Bold(true).
		Foreground(Accent)
	s.Selected = s.Selected.
		Foreground(Highlight).
		Background(Accent).
		Bold(true)
	t.SetStyles(s)

	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(Accent)

	return TableModel{
		config:  cfg,
		table:   t,
		spinner: sp,
		state:   stateTable,
	}
}

func (m TableModel) Init() tea.Cmd { return nil }

func (m TableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case DetailData:
		if msg.Err != nil {
			m.detail = lipgloss.NewStyle().Foreground(Bad).Render("  Error: " + msg.Err.Error())
		} else {
			m.detail = msg.Content
		}
		m.state = stateDetail
		return m, nil
	case tea.KeyPressMsg:
		switch m.state {
		case stateTable:
			switch msg.String() {
			case "q", "ctrl+c":
				return m, tea.Quit
			case "enter":
				if m.config.FetchDetail != nil && len(m.config.Rows) > 0 {
					m.state = stateLoading
					row := m.table.SelectedRow()
					return m, tea.Batch(m.spinner.Tick, m.config.FetchDetail(row))
				}
			}
		case stateLoading:
			switch msg.String() {
			case "q", "ctrl+c":
				return m, tea.Quit
			case "esc":
				m.state = stateTable
				return m, nil
			}
		case stateDetail:
			switch msg.String() {
			case "q", "ctrl+c":
				return m, tea.Quit
			case "esc":
				m.state = stateTable
				return m, nil
			}
		}
	case spinner.TickMsg:
		if m.state == stateLoading {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
	}

	if m.state == stateTable {
		var cmd tea.Cmd
		m.table, cmd = m.table.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m TableModel) View() tea.View {
	switch m.state {
	case stateLoading:
		return tea.NewView(fmt.Sprintf("\n  %s Loading...\n", m.spinner.View()))
	case stateDetail:
		help := LabelStyle.Render("  esc: back • q: quit")
		return tea.NewView("\n" + m.detail + "\n" + help + "\n")
	default:
		return tea.NewView(
			BaseStyle.Render(m.table.View()) + "\n  " + m.table.HelpView() + "\n",
		)
	}
}
