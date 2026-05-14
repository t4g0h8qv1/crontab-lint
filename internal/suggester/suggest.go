// Package suggester provides fix suggestions for invalid or suspicious crontab expressions.
package suggester

import (
	"fmt"
	"strings"
)

// Suggestion represents a human-readable fix recommendation.
type Suggestion struct {
	Field   string
	Issue   string
	Fix     string
}

// fieldNames maps field index to its canonical name.
var fieldNames = []string{"minute", "hour", "day-of-month", "month", "day-of-week"}

// fieldRanges maps field index to its valid [min, max].
var fieldRanges = [][2]int{
	{0, 59},
	{0, 23},
	{1, 31},
	{1, 12},
	{0, 7},
}

// Suggest returns a list of suggestions for the given crontab expression.
// It expects a 5-field standard cron expression.
func Suggest(expr string) []Suggestion {
	fields := strings.Fields(expr)
	if len(fields) != 5 {
		return []Suggestion{{
			Field: "expression",
			Issue: fmt.Sprintf("expected 5 fields, got %d", len(fields)),
			Fix:   "Provide exactly 5 fields: minute hour day-of-month month day-of-week",
		}}
	}

	var suggestions []Suggestion

	for i, field := range fields {
		name := fieldNames[i]
		min, max := fieldRanges[i][0], fieldRanges[i][1]

		if s := suggestField(field, name, min, max); s != nil {
			suggestions = append(suggestions, *s)
		}
	}

	// Check for redundant */1 pattern
	for i, field := range fields {
		if field == "*/1" {
			suggestions = append(suggestions, Suggestion{
				Field: fieldNames[i],
				Issue: "*/1 is redundant",
				Fix:   fmt.Sprintf("Replace '*/1' with '*' in the %s field", fieldNames[i]),
			})
		}
	}

	// Check DOM and DOW both set (not wildcards)
	if fields[2] != "*" && fields[4] != "*" {
		suggestions = append(suggestions, Suggestion{
			Field: "day-of-month/day-of-week",
			Issue: "both day-of-month and day-of-week are restricted",
			Fix:   "Set one of day-of-month or day-of-week to '*' to avoid ambiguous scheduling",
		})
	}

	return suggestions
}

// suggestField checks a single field value and returns a suggestion if needed.
func suggestField(field, name string, min, max int) *Suggestion {
	if field == "*" {
		return nil
	}
	var val int
	if _, err := fmt.Sscanf(field, "%d", &val); err == nil {
		if val < min || val > max {
			return &Suggestion{
				Field: name,
				Issue: fmt.Sprintf("value %d is out of range [%d, %d]", val, min, max),
				Fix:   fmt.Sprintf("Use a value between %d and %d for the %s field", min, max, name),
			}
		}
	}
	return nil
}
