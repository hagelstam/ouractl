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

func renderSleepDetail(day string, sleeps []api.Sleep, readiness []api.DailyReadiness) string {
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

	topRow := lipgloss.JoinHorizontal(lipgloss.Top, durationBox, " ", vitalsBox)

	var readinessBox string
	if len(readiness) > 0 {
		r := readiness[0]
		readinessBox = tui.RenderBox("Readiness", []tui.KeyValue{
			{Key: "Score", Value: tui.FmtScore(r.Score)},
			{Key: "Temp deviation", Value: tui.FmtTemp(r.TemperatureDeviation)},
			{Key: "Temp trend", Value: tui.FmtTemp(r.TemperatureTrendDeviation)},
			{Key: "Activity bal.", Value: tui.FmtScore(r.Contributors.ActivityBalance)},
			{Key: "Body temp", Value: tui.FmtScore(r.Contributors.BodyTemperature)},
			{Key: "HRV balance", Value: tui.FmtScore(r.Contributors.HRVBalance)},
			{Key: "Prev. activity", Value: tui.FmtScore(r.Contributors.PreviousDayActivity)},
			{Key: "Prev. night", Value: tui.FmtScore(r.Contributors.PreviousNight)},
			{Key: "Recovery", Value: tui.FmtScore(r.Contributors.RecoveryIndex)},
			{Key: "Resting HR", Value: tui.FmtScore(r.Contributors.RestingHeartRate)},
			{Key: "Sleep balance", Value: tui.FmtScore(r.Contributors.SleepBalance)},
		}, boxWidth)
	}

	result := title + "\n\n" + topRow
	if readinessBox != "" {
		result += "\n" + readinessBox
	}

	return result
}

func fetchSleepDayDetail(client *api.Client, day string) tea.Cmd {
	return func() tea.Msg {
		end := tui.NextDay(day)

		sleeps, err := client.GetSleep(day, end)
		if err != nil {
			return tui.DetailData{Err: err}
		}

		readiness, err := client.GetDailyReadiness(day, end)
		if err != nil {
			return tui.DetailData{Err: err}
		}

		return tui.DetailData{Content: renderSleepDetail(day, sleeps, readiness)}
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

		readiness, err := client.GetDailyReadiness(day, tui.NextDay(day))
		if err != nil {
			return tui.DetailData{Err: err}
		}

		var daySleeps []api.Sleep
		for _, s := range sleeps {
			if s.Day == day {
				daySleeps = append(daySleeps, s)
			}
		}

		return tui.DetailData{Content: renderSleepDetail(day, daySleeps, readiness)}
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

		sort.Slice(data, func(i, j int) bool {
			return data[i].Day > data[j].Day
		})

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
