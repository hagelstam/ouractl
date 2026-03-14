package tui

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
)

// KeyValue represents a labeled value for display in a box.
type KeyValue struct {
	Key   string
	Value string
}

// RenderBox renders a titled box with key-value pairs.
func RenderBox(title string, items []KeyValue, width int) string {
	titleStr := HeaderStyle.Render(title)

	maxKeyLen := 0
	for _, item := range items {
		if len(item.Key) > maxKeyLen {
			maxKeyLen = len(item.Key)
		}
	}

	var lines []string
	for _, item := range items {
		key := LabelStyle.Render(fmt.Sprintf("%-*s", maxKeyLen, item.Key))
		lines = append(lines, fmt.Sprintf("%s  %s", key, item.Value))
	}

	content := titleStr + "\n\n" + strings.Join(lines, "\n")

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Subtle).
		Padding(0, 1).
		Width(width)

	return box.Render(content)
}
