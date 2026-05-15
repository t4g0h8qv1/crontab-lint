package statistics

import (
	"encoding/json"
	"fmt"
	"strings"
)

// FormatText returns a human-readable summary of the statistics.
func FormatText(s Stats) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Expression : %s\n", s.Expression)
	fmt.Fprintf(&b, "Sample size: %d occurrences\n", s.SampleSize)
	fmt.Fprintf(&b, "Fires/hour : %.2f\n", s.FiresPerHour)
	fmt.Fprintf(&b, "Fires/day  : %.2f\n", s.FiresPerDay)
	fmt.Fprintf(&b, "Fires/week : %.2f\n", s.FiresPerWeek)
	fmt.Fprintf(&b, "Min gap    : %s\n", s.MinInterval)
	fmt.Fprintf(&b, "Max gap    : %s\n", s.MaxInterval)
	return b.String()
}

// statsJSON is the JSON-serialisable mirror of Stats.
type statsJSON struct {
	Expression   string  `json:"expression"`
	SampleSize   int     `json:"sample_size"`
	FiresPerHour float64 `json:"fires_per_hour"`
	FiresPerDay  float64 `json:"fires_per_day"`
	FiresPerWeek float64 `json:"fires_per_week"`
	MinIntervalS string  `json:"min_interval"`
	MaxIntervalS string  `json:"max_interval"`
}

// FormatJSON returns a JSON-encoded summary of the statistics.
func FormatJSON(s Stats) (string, error) {
	v := statsJSON{
		Expression:   s.Expression,
		SampleSize:   s.SampleSize,
		FiresPerHour: s.FiresPerHour,
		FiresPerDay:  s.FiresPerDay,
		FiresPerWeek: s.FiresPerWeek,
		MinIntervalS: s.MinInterval.String(),
		MaxIntervalS: s.MaxInterval.String(),
	}
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}
