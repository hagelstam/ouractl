package ouractl

import (
	"fmt"
	"sort"
	"time"

	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"github.com/hagelstam/ouractl/internal/api"
	"github.com/hagelstam/ouractl/internal/auth"
	"github.com/hagelstam/ouractl/internal/tui"
	"github.com/spf13/cobra"
)

func fetchReadinessDetail(client *api.Client) func(row table.Row) tea.Cmd {
	return func(row table.Row) tea.Cmd {
		return fetchReadinessDayDetail(client, row[0])
	}
}

func renderReadinessDetail(day string, data []api.DailyReadiness) string {
	title := tui.HeaderStyle.Render(fmt.Sprintf(" Readiness %s", day))

	if len(data) == 0 {
		return title + "\n\n  No readiness data found."
	}

	r := data[0]
	boxWidth := 38

	scoreBox := tui.RenderBox("Readiness", []tui.KeyValue{
		{Key: "Score", Value: tui.FmtScore(r.Score)},
		{Key: "Temp deviation", Value: tui.FmtTemp(r.TemperatureDeviation)},
		{Key: "Temp trend", Value: tui.FmtTemp(r.TemperatureTrendDeviation)},
	}, boxWidth)

	contribBox := tui.RenderBox("Contributors", []tui.KeyValue{
		{Key: "Activity bal.", Value: tui.FmtScore(r.Contributors.ActivityBalance)},
		{Key: "Body temp", Value: tui.FmtScore(r.Contributors.BodyTemperature)},
		{Key: "HRV balance", Value: tui.FmtScore(r.Contributors.HRVBalance)},
		{Key: "Prev. activity", Value: tui.FmtScore(r.Contributors.PreviousDayActivity)},
		{Key: "Prev. night", Value: tui.FmtScore(r.Contributors.PreviousNight)},
		{Key: "Recovery", Value: tui.FmtScore(r.Contributors.RecoveryIndex)},
		{Key: "Resting HR", Value: tui.FmtScore(r.Contributors.RestingHeartRate)},
		{Key: "Sleep balance", Value: tui.FmtScore(r.Contributors.SleepBalance)},
		{Key: "Sleep regularity", Value: tui.FmtScore(r.Contributors.SleepRegularity)},
	}, boxWidth)

	return title + "\n\n" + scoreBox + "\n" + contribBox
}

func fetchReadinessDayDetail(client *api.Client, day string) tea.Cmd {
	return func() tea.Msg {
		end := tui.NextDay(day)
		data, err := client.GetDailyReadiness(day, end)
		if err != nil {
			return tui.DetailData{Err: err}
		}
		return tui.DetailData{Content: renderReadinessDetail(day, data)}
	}
}

func fetchReadinessLatestDetail(client *api.Client) tea.Cmd {
	return func() tea.Msg {
		endDate := tui.Tomorrow()
		startDate := time.Now().AddDate(0, 0, -7).Format("2006-01-02")

		data, err := client.GetDailyReadiness(startDate, endDate)
		if err != nil {
			return tui.DetailData{Err: err}
		}
		if len(data) == 0 {
			return tui.DetailData{Content: "No recent readiness data found."}
		}

		sort.Slice(data, func(i, j int) bool {
			return data[i].Day > data[j].Day
		})

		return tui.DetailData{Content: renderReadinessDetail(data[0].Day, data[:1])}
	}
}

var (
	readinessDays   int
	readinessLatest bool
)

var readinessCmd = &cobra.Command{
	Use:   "readiness [date]",
	Short: "View daily readiness data",
	Long:  "View daily readiness data in a table. Optionally pass a date (YYYY-MM-DD) for a detail view.",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		token, err := auth.LoadToken()
		if err != nil {
			return err
		}

		client := api.NewClient(token)

		if len(args) == 1 {
			return runDetailProgram(fetchReadinessDayDetail(client, args[0]))
		}
		if readinessLatest {
			return runDetailProgram(fetchReadinessLatestDetail(client))
		}
		if err := tui.ValidateDays(readinessDays); err != nil {
			return err
		}

		endDate := tui.Tomorrow()
		startDate := time.Now().AddDate(0, 0, -readinessDays).Format("2006-01-02")

		data, err := client.GetDailyReadiness(startDate, endDate)
		if err != nil {
			return err
		}
		if len(data) == 0 {
			fmt.Printf("No readiness data found for the last %d days.\n", readinessDays)
			return nil
		}

		columns := []table.Column{
			{Title: "Date", Width: 12},
			{Title: "Score", Width: 7},
			{Title: "Temp Dev", Width: 10},
			{Title: "Temp Trend", Width: 12},
		}

		rows := make([]table.Row, len(data))
		for i, d := range data {
			rows[i] = table.Row{
				d.Day,
				tui.FmtScore(d.Score),
				tui.FmtTemp(d.TemperatureDeviation),
				tui.FmtTemp(d.TemperatureTrendDeviation),
			}
		}
		rows = tui.FillDateGaps(rows, startDate, endDate, len(columns))

		model := tui.NewTableModel(tui.TableConfig{
			Columns:     columns,
			Rows:        rows,
			Width:       47,
			FetchDetail: fetchReadinessDetail(client),
		})

		_, err = tea.NewProgram(model).Run()
		return err
	},
}

func init() {
	readinessCmd.Flags().IntVar(&readinessDays, "days", 7, "Number of days to display (1-30)")
	readinessCmd.Flags().
		BoolVar(&readinessLatest, "latest", false, "Show details for the most recent day")
	rootCmd.AddCommand(readinessCmd)
}
