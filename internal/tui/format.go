package tui

import (
	"fmt"
	"strconv"
	"time"
)

func fmtPtr[T any](v *T, fn func(T) string) string {
	if v == nil {
		return "-"
	}
	return fn(*v)
}

func NextDay(day string) string {
	t, err := time.Parse("2006-01-02", day)
	if err != nil {
		return day
	}
	return t.AddDate(0, 0, 1).Format("2006-01-02")
}

func Tomorrow() string {
	return time.Now().AddDate(0, 0, 1).Format("2006-01-02")
}

func FmtScore(v *int) string {
	return fmtPtr(v, strconv.Itoa)
}

func FmtFloat(v *float64) string {
	return fmtPtr(v, func(f float64) string {
		return fmt.Sprintf("%.1f", f)
	})
}

func FmtPercent(v *int) string {
	return fmtPtr(v, func(n int) string {
		return fmt.Sprintf("%d%%", n)
	})
}

func FmtDurationPtr(v *int) string {
	return fmtPtr(v, FmtDuration)
}

func FmtTemp(v *float64) string {
	return fmtPtr(v, func(f float64) string {
		if f >= 0 {
			return fmt.Sprintf("+%.1f°C", f)
		}
		return fmt.Sprintf("%.1f°C", f)
	})
}

func FmtDuration(seconds int) string {
	h := seconds / 3600
	m := (seconds % 3600) / 60
	if h > 0 {
		return fmt.Sprintf("%dh %dm", h, m)
	}
	return fmt.Sprintf("%dm", m)
}

func FmtTime(isoTimestamp string) string {
	formats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05-07:00",
		"2006-01-02T15:04:05.000-07:00",
		"2006-01-02T15:04:05",
	}
	for _, layout := range formats {
		if t, err := time.Parse(layout, isoTimestamp); err == nil {
			return t.Format("15:04")
		}
	}
	return isoTimestamp
}

func FmtDistance(meters int) string {
	return fmt.Sprintf("%.1f km", float64(meters)/1000)
}

func WithUnit(formatted, unit string) string {
	if formatted == "-" {
		return "-"
	}
	return formatted + " " + unit
}

func ValidateDays(days int) error {
	if days < 1 || days > 30 {
		return fmt.Errorf("flag --days must be between 1 and 30, got %d", days)
	}
	return nil
}
