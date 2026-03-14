package tui

import "charm.land/lipgloss/v2"

var (
	Accent    = lipgloss.Color("99")
	Subtle    = lipgloss.Color("242")
	Highlight = lipgloss.Color("230")
	Good      = lipgloss.Color("78")
	Warn      = lipgloss.Color("214")
	Bad       = lipgloss.Color("196")

	BaseStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(Subtle)

	HeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Accent)

	SelectedStyle = lipgloss.NewStyle().
			Foreground(Highlight).
			Background(Accent).
			Bold(true)

	LabelStyle = lipgloss.NewStyle().Foreground(Subtle)

	ValueStyle = lipgloss.NewStyle().Bold(true)
)
