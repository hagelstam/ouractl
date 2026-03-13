package cmd

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"time"

	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/hagelstam/oura-cli/internal/api"
	"github.com/hagelstam/oura-cli/internal/auth"
	"github.com/spf13/cobra"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type sleepModel struct {
	table table.Model
}

func (m sleepModel) Init() tea.Cmd { return nil }

func (m sleepModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m sleepModel) View() tea.View {
	return tea.NewView(baseStyle.Render(m.table.View()) + "\n  " + m.table.HelpView() + "\n")
}

func fmtScore(v *int) string {
	if v == nil {
		return "-"
	}
	return strconv.Itoa(*v)
}

var sleepCmd = &cobra.Command{
	Use:   "sleep",
	Short: "View daily sleep data for the last 7 days",
	RunE: func(cmd *cobra.Command, args []string) error {
		token, err := auth.LoadToken()
		if err != nil {
			return err
		}

		client := api.NewClient(token)
		endDate := time.Now().Format("2006-01-02")
		startDate := time.Now().AddDate(0, 0, -7).Format("2006-01-02")

		data, err := client.GetDailySleep(startDate, endDate)
		if err != nil {
			return err
		}
		if len(data) == 0 {
			fmt.Println("No sleep data found for the last 7 days.")
			return nil
		}

		sort.Slice(data, func(i, j int) bool {
			return data[i].Day > data[j].Day
		})

		columns := []table.Column{
			{Title: "Date", Width: 12},
			{Title: "Score", Width: 7},
			{Title: "Deep", Width: 6},
			{Title: "Effic.", Width: 7},
			{Title: "REM", Width: 6},
			{Title: "Rest.", Width: 6},
			{Title: "Total", Width: 7},
		}

		rows := make([]table.Row, len(data))
		for i, d := range data {
			rows[i] = table.Row{
				d.Day,
				fmtScore(d.Score),
				fmtScore(d.Contributors.DeepSleep),
				fmtScore(d.Contributors.Efficiency),
				fmtScore(d.Contributors.REMSleep),
				fmtScore(d.Contributors.Restfulness),
				fmtScore(d.Contributors.TotalSleep),
			}
		}

		t := table.New(
			table.WithColumns(columns),
			table.WithRows(rows),
			table.WithFocused(true),
			table.WithHeight(len(rows)),
			table.WithWidth(59),
		)

		s := table.DefaultStyles()
		s.Header = s.Header.
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240")).
			BorderBottom(true).
			Bold(false)
		s.Selected = s.Selected.
			Foreground(lipgloss.Color("229")).
			Background(lipgloss.Color("57")).
			Bold(false)
		t.SetStyles(s)

		if _, err := tea.NewProgram(sleepModel{t}).Run(); err != nil {
			fmt.Fprintln(os.Stderr, "Error running program:", err)
			os.Exit(1)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(sleepCmd)
}
