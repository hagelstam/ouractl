package ouractl

import (
	"fmt"
	"sort"
	"strconv"
	"time"

	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/hagelstam/ouractl/internal/api"
	"github.com/hagelstam/ouractl/internal/auth"
	"github.com/hagelstam/ouractl/internal/tui"
	"github.com/spf13/cobra"
)

func fetchActivityDetail(client *api.Client) func(row table.Row) tea.Cmd {
	return func(row table.Row) tea.Cmd {
		return fetchActivityDayDetail(client, row[0])
	}
}

func renderActivityDetail(day string, data []api.DailyActivity) string {
	title := tui.HeaderStyle.Render(fmt.Sprintf(" Activity %s", day))

	if len(data) == 0 {
		return title + "\n\n  No activity data found."
	}

	a := data[0]
	boxWidth := 38

	summaryBox := tui.RenderBox("Summary", []tui.KeyValue{
		{Key: "Score", Value: tui.FmtScore(a.Score)},
		{Key: "Steps", Value: strconv.Itoa(a.Steps)},
		{Key: "Active cal", Value: strconv.Itoa(a.ActiveCalories) + " kcal"},
		{Key: "Total cal", Value: strconv.Itoa(a.TotalCalories) + " kcal"},
		{Key: "Walking dist", Value: tui.FmtDistance(a.EquivalentWalkingDistance)},
		{Key: "Target cal", Value: strconv.Itoa(a.TargetCalories) + " kcal"},
	}, boxWidth)

	durationBox := tui.RenderBox("Activity Time", []tui.KeyValue{
		{Key: "High", Value: tui.FmtDuration(a.HighActivityTime)},
		{Key: "Medium", Value: tui.FmtDuration(a.MediumActivityTime)},
		{Key: "Low", Value: tui.FmtDuration(a.LowActivityTime)},
		{Key: "Resting", Value: tui.FmtDuration(a.RestingTime)},
		{Key: "Sedentary", Value: tui.FmtDuration(a.SedentaryTime)},
		{Key: "Non-wear", Value: tui.FmtDuration(a.NonWearTime)},
	}, boxWidth)

	topRow := lipgloss.JoinHorizontal(lipgloss.Top, summaryBox, " ", durationBox)

	contribBox := tui.RenderBox("Contributors", []tui.KeyValue{
		{Key: "Daily targets", Value: tui.FmtScore(a.Contributors.MeetDailyTargets)},
		{Key: "Move hourly", Value: tui.FmtScore(a.Contributors.MoveEveryHour)},
		{Key: "Recovery time", Value: tui.FmtScore(a.Contributors.RecoveryTime)},
		{Key: "Stay active", Value: tui.FmtScore(a.Contributors.StayActive)},
		{Key: "Training freq", Value: tui.FmtScore(a.Contributors.TrainingFrequency)},
		{Key: "Training volume", Value: tui.FmtScore(a.Contributors.TrainingVolume)},
	}, boxWidth)

	return title + "\n\n" + topRow + "\n" + contribBox
}

func fetchActivityDayDetail(client *api.Client, day string) tea.Cmd {
	return func() tea.Msg {
		end := tui.NextDay(day)
		data, err := client.GetDailyActivity(day, end)
		if err != nil {
			return tui.DetailData{Err: err}
		}
		return tui.DetailData{Content: renderActivityDetail(day, data)}
	}
}

func fetchActivityLatestDetail(client *api.Client) tea.Cmd {
	return func() tea.Msg {
		endDate := tui.Tomorrow()
		startDate := time.Now().AddDate(0, 0, -7).Format("2006-01-02")

		data, err := client.GetDailyActivity(startDate, endDate)
		if err != nil {
			return tui.DetailData{Err: err}
		}
		if len(data) == 0 {
			return tui.DetailData{Content: "No recent activity data found."}
		}
		sort.Slice(data, func(i, j int) bool {
			return data[i].Day > data[j].Day
		})

		return tui.DetailData{Content: renderActivityDetail(data[0].Day, data[:1])}
	}
}

var (
	activityDays   int
	activityLatest bool
)

var activityCmd = &cobra.Command{
	Use:   "activity [date]",
	Short: "View daily activity data",
	Long:  "View daily activity data in a table. Optionally pass a date (YYYY-MM-DD) for a detail view.",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		token, err := auth.LoadToken()
		if err != nil {
			return err
		}

		client := api.NewClient(token)

		if len(args) == 1 {
			return runDetailProgram(fetchActivityDayDetail(client, args[0]))
		}
		if activityLatest {
			return runDetailProgram(fetchActivityLatestDetail(client))
		}
		if err := tui.ValidateDays(activityDays); err != nil {
			return err
		}

		endDate := tui.Tomorrow()
		startDate := time.Now().AddDate(0, 0, -activityDays).Format("2006-01-02")

		data, err := client.GetDailyActivity(startDate, endDate)
		if err != nil {
			return err
		}
		if len(data) == 0 {
			fmt.Printf("No activity data found for the last %d days.\n", activityDays)
			return nil
		}
		sort.Slice(data, func(i, j int) bool {
			return data[i].Day > data[j].Day
		})

		columns := []table.Column{
			{Title: "Date", Width: 12},
			{Title: "Score", Width: 7},
			{Title: "Steps", Width: 7},
			{Title: "Active Cal", Width: 12},
			{Title: "Walking Dist", Width: 12},
		}

		rows := make([]table.Row, len(data))
		for i, d := range data {
			rows[i] = table.Row{
				d.Day,
				tui.FmtScore(d.Score),
				strconv.Itoa(d.Steps),
				strconv.Itoa(d.ActiveCalories) + " kcal",
				tui.FmtDistance(d.EquivalentWalkingDistance),
			}
		}

		model := tui.NewTableModel(tui.TableConfig{
			Columns:     columns,
			Rows:        rows,
			Width:       54,
			FetchDetail: fetchActivityDetail(client),
		})

		_, err = tea.NewProgram(model).Run()
		return err
	},
}

func init() {
	activityCmd.Flags().IntVar(&activityDays, "days", 7, "Number of days to display (1-30)")
	activityCmd.Flags().
		BoolVar(&activityLatest, "latest", false, "Show details for the most recent day")
	rootCmd.AddCommand(activityCmd)
}
