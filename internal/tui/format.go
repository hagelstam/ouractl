package tui

import (
	"fmt"
	"strconv"
	"time"
)

// NextDay returns the day after the given "YYYY-MM-DD" date string.
// Falls back to returning the input unchanged on parse error.
func NextDay(day string) string {
	t, err := time.Parse("2006-01-02", day)
	if err != nil {
		return day
	}
	return t.AddDate(0, 0, 1).Format("2006-01-02")
}

// Tomorrow returns tomorrow's date as "YYYY-MM-DD".
func Tomorrow() string {
	return time.Now().AddDate(0, 0, 1).Format("2006-01-02")
}

// FmtScore formats an optional score integer.
func FmtScore(v *int) string {
	if v == nil {
		return "-"
	}
	return strconv.Itoa(*v)
}

// FmtDuration formats seconds into "Xh Ym" format.
func FmtDuration(seconds int) string {
	h := seconds / 3600
	m := (seconds % 3600) / 60
	if h > 0 {
		return fmt.Sprintf("%dh %dm", h, m)
	}
	return fmt.Sprintf("%dm", m)
}

// FmtDurationPtr formats an optional duration in seconds.
func FmtDurationPtr(v *int) string {
	if v == nil {
		return "-"
	}
	return FmtDuration(*v)
}

// FmtTime parses an ISO timestamp and returns "HH:MM".
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

// FmtTemp formats an optional temperature deviation.
func FmtTemp(v *float64) string {
	if v == nil {
		return "-"
	}
	if *v >= 0 {
		return fmt.Sprintf("+%.1f°C", *v)
	}
	return fmt.Sprintf("%.1f°C", *v)
}

// FmtDistance formats meters into "X.X km".
func FmtDistance(meters int) string {
	km := float64(meters) / 1000.0
	return fmt.Sprintf("%.1f km", km)
}

// FmtFloat formats an optional float to one decimal.
func FmtFloat(v *float64) string {
	if v == nil {
		return "-"
	}
	return fmt.Sprintf("%.1f", *v)
}

// FmtPercent formats an optional integer as a percentage.
func FmtPercent(v *int) string {
	if v == nil {
		return "-"
	}
	return fmt.Sprintf("%d%%", *v)
}

// WithUnit appends a unit suffix only when the value is present.
func WithUnit(formatted, unit string) string {
	if formatted == "-" {
		return "-"
	}
	return formatted + " " + unit
}

// ValidateDays checks that days is between 1 and 30.
func ValidateDays(days int) error {
	if days < 1 || days > 30 {
		return fmt.Errorf("flag --days must be between 1 and 30, got %d", days)
	}
	return nil
}
