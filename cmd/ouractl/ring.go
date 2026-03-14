package ouractl

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/hagelstam/ouractl/internal/api"
	"github.com/hagelstam/ouractl/internal/auth"
	"github.com/hagelstam/ouractl/internal/tui"
	"github.com/spf13/cobra"
)

func fmtRingField(v *string) string {
	if v == nil {
		return "-"
	}
	return *v
}

var ringCmd = &cobra.Command{
	Use:   "ring",
	Short: "View ring hardware information",
	RunE: func(cmd *cobra.Command, args []string) error {
		token, err := auth.LoadToken()
		if err != nil {
			return err
		}

		client := api.NewClient(token)
		rings, err := client.GetRingConfig()
		if err != nil {
			return err
		}
		if len(rings) == 0 {
			fmt.Println("No ring configuration found.")
			return nil
		}

		ring := rings[len(rings)-1] // Show the most recently set up ring
		hardware := fmtRingField(ring.HardwareType)
		design := fmtRingField(ring.Design)
		color := fmtRingField(ring.Color)
		size := tui.FmtScore(ring.Size)
		firmware := fmtRingField(ring.FirmwareVersion)
		setupAt := fmtRingField(ring.SetUpAt)
		if ring.SetUpAt != nil {
			setupAt = tui.FmtTime(*ring.SetUpAt)
			// If it has a date portion, show the full date.
			if len(*ring.SetUpAt) >= 10 {
				setupAt = (*ring.SetUpAt)[:10]
			}
		}

		title := tui.HeaderStyle.Render("Ring Info")
		labelStyle := tui.LabelStyle.Width(12)
		lines := []struct{ label, value string }{
			{"Hardware", strings.ReplaceAll(hardware, "_", " ")},
			{"Design", design},
			{"Color", color},
			{"Size", size},
			{"Firmware", firmware},
			{"Set up", setupAt},
		}

		var rows []string
		for _, l := range lines {
			rows = append(rows, fmt.Sprintf("  %s  %s", labelStyle.Render(l.label), l.value))
		}

		box := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(tui.Subtle).
			Padding(0, 1).
			Width(36)

		content := title + "\n" + strings.Join(rows, "\n")
		fmt.Println(box.Render(content))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(ringCmd)
}
