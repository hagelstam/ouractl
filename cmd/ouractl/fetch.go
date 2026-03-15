package ouractl

import (
	"sort"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/hagelstam/ouractl/internal/api"
	"github.com/hagelstam/ouractl/internal/auth"
	"github.com/hagelstam/ouractl/internal/tui"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

var ringASCII = `
		████████████

        ████████████
      ██            ██
    ██                ██
   ██                  ██
   ██                  ██
   ██                  ██
    ██                ██
      ██            ██
        ████████████
`

func fetchAllData() tea.Cmd {
	return func() tea.Msg {
		token, err := auth.LoadToken()
		if err != nil {
			return tui.DetailData{Err: err}
		}

		client := api.NewClient(token)
		startDate := time.Now().AddDate(0, 0, -2).Format("2006-01-02")
		endDate := tui.Tomorrow()

		var (
			dailySleep     []api.DailySleep
			dailyReadiness []api.DailyReadiness
			dailyActivity  []api.DailyActivity
			sleepDetail    []api.Sleep
			personalInfo   *api.PersonalInfo
			rings          []api.RingConfig
		)

		g := new(errgroup.Group)

		g.Go(func() error {
			var err error
			dailySleep, err = client.GetDailySleep(startDate, endDate)
			return err
		})
		g.Go(func() error {
			var err error
			dailyReadiness, err = client.GetDailyReadiness(startDate, endDate)
			return err
		})
		g.Go(func() error {
			var err error
			dailyActivity, err = client.GetDailyActivity(startDate, endDate)
			return err
		})
		g.Go(func() error {
			var err error
			sleepDetail, err = client.GetSleep(startDate, endDate)
			return err
		})
		g.Go(func() error {
			var err error
			personalInfo, err = client.GetPersonalInfo()
			return err
		})
		g.Go(func() error {
			var err error
			rings, err = client.GetRingConfig()
			return err
		})

		if err := g.Wait(); err != nil {
			return tui.DetailData{Err: err}
		}

		return tui.DetailData{
			Content: buildFetchOutput(
				dailySleep,
				dailyReadiness,
				dailyActivity,
				sleepDetail,
				personalInfo,
				rings,
			),
		}
	}
}

func buildFetchOutput(
	dailySleep []api.DailySleep,
	dailyReadiness []api.DailyReadiness,
	dailyActivity []api.DailyActivity,
	sleepDetail []api.Sleep,
	personalInfo *api.PersonalInfo,
	rings []api.RingConfig,
) string {
	art := lipgloss.NewStyle().
		Foreground(tui.Accent).
		MarginRight(3).
		Render(ringASCII)

	var infoLines []string

	// Title line.
	title := "ouractl"
	if personalInfo != nil && len(personalInfo.ID) >= 8 {
		title = "oura:" + personalInfo.ID[:8]
	}
	infoLines = append(infoLines, tui.HeaderStyle.Render(title))
	infoLines = append(infoLines, tui.LabelStyle.Render(strings.Repeat("─", 34)))

	// Sleep score.
	if len(dailySleep) > 0 {
		sort.Slice(dailySleep, func(i, j int) bool {
			return dailySleep[i].Day > dailySleep[j].Day
		})
		infoLines = append(infoLines, fmtInfoLine("Sleep", tui.FmtScore(dailySleep[0].Score)))
	}

	// Readiness score.
	if len(dailyReadiness) > 0 {
		sort.Slice(dailyReadiness, func(i, j int) bool {
			return dailyReadiness[i].Day > dailyReadiness[j].Day
		})
		infoLines = append(
			infoLines,
			fmtInfoLine("Readiness", tui.FmtScore(dailyReadiness[0].Score)),
		)
	}

	// Activity score + steps.
	if len(dailyActivity) > 0 {
		sort.Slice(dailyActivity, func(i, j int) bool {
			return dailyActivity[i].Day > dailyActivity[j].Day
		})
		infoLines = append(infoLines, fmtInfoLine("Activity", tui.FmtScore(dailyActivity[0].Score)))
		yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
		for _, a := range dailyActivity {
			if a.Day == yesterday {
				infoLines = append(infoLines, fmtInfoLine("Steps", tui.FmtSteps(a.Steps)))
				break
			}
		}
	}

	// Sleep duration (longest sleep period from latest day).
	if len(sleepDetail) > 0 {
		sort.Slice(sleepDetail, func(i, j int) bool {
			return sleepDetail[i].Day > sleepDetail[j].Day
		})
		latestDay := sleepDetail[0].Day
		s := sleepDetail[0]
		for _, sd := range sleepDetail[1:] {
			if sd.Day != latestDay {
				break
			}
			if sd.TimeInBed > s.TimeInBed {
				s = sd
			}
		}
		infoLines = append(
			infoLines,
			fmtInfoLine("Slept", tui.FmtDurationPtr(s.TotalSleepDuration)),
		)
	}

	// Ring info.
	if len(rings) > 0 {
		ring := rings[len(rings)-1]
		if ring.HardwareType != nil {
			hw := strings.ReplaceAll(*ring.HardwareType, "_", " ")
			infoLines = append(infoLines, fmtInfoLine("Ring", hw))
		}
	}

	// Neofetch style color blocks.
	infoLines = append(infoLines, "")
	var blocks []string
	for _, c := range []string{"1", "2", "3", "4", "5", "6", "7", "8"} {
		blocks = append(blocks, lipgloss.NewStyle().Background(lipgloss.Color(c)).Render("   "))
	}
	infoLines = append(infoLines, strings.Join(blocks, ""))

	info := strings.Join(infoLines, "\n")
	return lipgloss.JoinHorizontal(lipgloss.Center, art, info) + "\n"
}

func fmtInfoLine(label, value string) string {
	l := tui.LabelStyle.Width(14).Render(label)
	return l + tui.ValueStyle.Render(value)
}

func init() {
	rootCmd.RunE = func(cmd *cobra.Command, args []string) error {
		return runDetailProgram(fetchAllData())
	}
}
