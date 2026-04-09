package ouractl

import (
	"fmt"
	"sort"
	"time"

	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/hagelstam/ouractl/internal/api"
	"github.com/hagelstam/ouractl/internal/auth"
	"github.com/hagelstam/ouractl/internal/tui"
	"github.com/spf13/cobra"
)

func fetchSleepDetail(client *api.Client) func(row table.Row) tea.Cmd {
	return func(row table.Row) tea.Cmd {
		return fetchSleepDayDetail(client, row[0])
	}
}

func renderSleepDetail(day string, sleeps []api.Sleep) string {
	title := tui.HeaderStyle.Render(fmt.Sprintf(" Sleep %s", day))

	// Find the longest sleep period.
	var s *api.Sleep
	for i := range sleeps {
		if s == nil || sleeps[i].TimeInBed > s.TimeInBed {
			s = &sleeps[i]
		}
	}
	if s == nil {
		return title + "\n\n  No detailed sleep data found."
	}

	boxWidth := 38

	durationBox := tui.RenderBox("Duration", []tui.KeyValue{
		{Key: "Bedtime", Value: tui.FmtTime(s.BedtimeStart) + " → " + tui.FmtTime(s.BedtimeEnd)},
		{Key: "Total", Value: tui.FmtDurationPtr(s.TotalSleepDuration)},
		{Key: "Time in bed", Value: tui.FmtDuration(s.TimeInBed)},
		{Key: "Deep", Value: tui.FmtDurationPtr(s.DeepSleepDuration)},
		{Key: "REM", Value: tui.FmtDurationPtr(s.REMSleepDuration)},
		{Key: "Light", Value: tui.FmtDurationPtr(s.LightSleepDuration)},
		{Key: "Awake", Value: tui.FmtDurationPtr(s.AwakeTime)},
	}, boxWidth)

	vitalsBox := tui.RenderBox("Vitals", []tui.KeyValue{
		{Key: "Avg HR", Value: tui.WithUnit(tui.FmtFloat(s.AverageHeartRate), "BPM")},
		{Key: "Lowest HR", Value: tui.WithUnit(tui.FmtScore(s.LowestHeartRate), "BPM")},
		{Key: "HRV", Value: tui.WithUnit(tui.FmtScore(s.AverageHRV), "ms")},
		{Key: "Breathing", Value: tui.WithUnit(tui.FmtFloat(s.AverageBreath), "br/min")},
		{Key: "Efficiency", Value: tui.FmtPercent(s.Efficiency)},
		{Key: "Latency", Value: tui.FmtDurationPtr(s.Latency)},
	}, boxWidth)

	return title + "\n\n" + lipgloss.JoinHorizontal(lipgloss.Top, durationBox, " ", vitalsBox)
}

func fetchSleepDayDetail(client *api.Client, day string) tea.Cmd {
	return func() tea.Msg {
		sleeps, err := client.GetSleep(tui.PrevDay(day), day)
		if err != nil {
			return tui.DetailData{Err: err}
		}
		return tui.DetailData{Content: renderSleepDetail(day, sleeps)}
	}
}

func fetchSleepLatestDetail(client *api.Client) tea.Cmd {
	return func() tea.Msg {
		endDate := tui.Tomorrow()
		startDate := time.Now().AddDate(0, 0, -7).Format("2006-01-02")

		sleeps, err := client.GetSleep(startDate, endDate)
		if err != nil {
			return tui.DetailData{Err: err}
		}
		if len(sleeps) == 0 {
			return tui.DetailData{Content: "No recent sleep data found."}
		}

		sort.Slice(sleeps, func(i, j int) bool {
			return sleeps[i].Day > sleeps[j].Day
		})
		day := sleeps[0].Day

		var daySleeps []api.Sleep
		for _, s := range sleeps {
			if s.Day == day {
				daySleeps = append(daySleeps, s)
			}
		}

		return tui.DetailData{Content: renderSleepDetail(day, daySleeps)}
	}
}

func runDetailProgram(fetch tea.Cmd) error {
	_, err := tea.NewProgram(tui.NewDetailModel(fetch)).Run()
	return err
}

var (
	sleepDays   int
	sleepLatest bool
)

var sleepCmd = &cobra.Command{
	Use:   "sleep [date]",
	Short: "View daily sleep data",
	Long:  "View daily sleep data in a table. Optionally pass a date (YYYY-MM-DD) for a detail view.",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		token, err := auth.LoadToken()
		if err != nil {
			return err
		}

		client := api.NewClient(token)

		if len(args) == 1 {
			return runDetailProgram(fetchSleepDayDetail(client, args[0]))
		}
		if sleepLatest {
			return runDetailProgram(fetchSleepLatestDetail(client))
		}
		if err := tui.ValidateDays(sleepDays); err != nil {
			return err
		}

		endDate := tui.Tomorrow()
		startDate := time.Now().AddDate(0, 0, -sleepDays).Format("2006-01-02")

		data, err := client.GetDailySleep(startDate, endDate)
		if err != nil {
			return err
		}
		if len(data) == 0 {
			fmt.Printf("No sleep data found for the last %d days.\n", sleepDays)
			return nil
		}

		columns := []table.Column{
			{Title: "Date", Width: 12},
			{Title: "Score", Width: 6},
			{Title: "Deep", Width: 6},
			{Title: "REM", Width: 6},
			{Title: "Effic.", Width: 6},
			{Title: "Rest.", Width: 6},
			{Title: "Total", Width: 7},
		}

		rows := make([]table.Row, len(data))
		for i, d := range data {
			rows[i] = table.Row{
				d.Day,
				tui.FmtScore(d.Score),
				tui.FmtScore(d.Contributors.DeepSleep),
				tui.FmtScore(d.Contributors.REMSleep),
				tui.FmtScore(d.Contributors.Efficiency),
				tui.FmtScore(d.Contributors.Restfulness),
				tui.FmtScore(d.Contributors.TotalSleep),
			}
		}
		rows = tui.FillDateGaps(rows, startDate, endDate, len(columns))

		model := tui.NewTableModel(tui.TableConfig{
			Columns:     columns,
			Rows:        rows,
			Width:       59,
			FetchDetail: fetchSleepDetail(client),
		})

		_, err = tea.NewProgram(model).Run()
		return err
	},
}

func init() {
	sleepCmd.Flags().IntVar(&sleepDays, "days", 7, "Number of days to display (1-30)")
	sleepCmd.Flags().BoolVar(&sleepLatest, "latest", false, "Show details for the most recent day")
	rootCmd.AddCommand(sleepCmd)
}
