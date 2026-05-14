// Package highlighter provides syntax highlighting for crontab expressions,
// marking fields with colors or symbols to aid readability in terminal output.
package highlighter

import (
	"fmt"
	"strings"
)

// ANSI color codes for terminal output.
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
	colorBold   = "\033[1m"
)

// FieldColor maps each crontab field index to a distinct ANSI color.
var FieldColor = []string{
	colorCyan,   // minute
	colorGreen,  // hour
	colorYellow, // day-of-month
	colorRed,    // month
	colorBold,   // day-of-week
}

// FieldLabel maps each crontab field index to its human-readable name.
var FieldLabel = []string{
	"minute",
	"hour",
	"dom",
	"month",
	"dow",
}

// HighlightResult holds the highlighted expression and a legend.
type HighlightResult struct {
	// Highlighted is the expression with ANSI color codes applied per field.
	Highlighted string
	// Legend maps field names to their colored sample.
	Legend map[string]string
}

// Highlight applies per-field ANSI colors to a crontab expression string.
// It returns an error if the expression does not contain exactly 5 fields.
func Highlight(expression string) (HighlightResult, error) {
	fields := strings.Fields(expression)
	if len(fields) != 5 {
		return HighlightResult{}, fmt.Errorf(
			"highlighter: expected 5 fields, got %d", len(fields),
		)
	}

	colored := make([]string, 5)
	legend := make(map[string]string, 5)

	for i, field := range fields {
		colored[i] = FieldColor[i] + field + colorReset
		legend[FieldLabel[i]] = FieldColor[i] + FieldLabel[i] + colorReset
	}

	return HighlightResult{
		Highlighted: strings.Join(colored, " "),
		Legend:      legend,
	}, nil
}

// Strip removes all ANSI escape codes from a string, returning plain text.
func Strip(s string) string {
	replacer := strings.NewReplacer(
		colorReset, "",
		colorRed, "",
		colorGreen, "",
		colorYellow, "",
		colorCyan, "",
		colorBold, "",
	)
	return replacer.Replace(s)
}
