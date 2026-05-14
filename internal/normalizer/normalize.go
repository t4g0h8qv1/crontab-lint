// Package normalizer provides functionality to normalize crontab expressions
// by replacing common aliases and shorthand notations with their canonical forms.
package normalizer

import "strings"

// aliases maps common crontab shorthand expressions to their canonical forms.
var aliases = map[string]string{
	"@yearly":   "0 0 1 1 *",
	"@annually": "0 0 1 1 *",
	"@monthly":  "0 0 1 * *",
	"@weekly":   "0 0 * * 0",
	"@daily":    "0 0 * * *",
	"@midnight": "0 0 * * *",
	"@hourly":   "0 * * * *",
}

// monthNames maps month name abbreviations to their numeric equivalents.
var monthNames = map[string]string{
	"jan": "1", "feb": "2", "mar": "3", "apr": "4",
	"may": "5", "jun": "6", "jul": "7", "aug": "8",
	"sep": "9", "oct": "10", "nov": "11", "dec": "12",
}

// dowNames maps day-of-week name abbreviations to their numeric equivalents.
var dowNames = map[string]string{
	"sun": "0", "mon": "1", "tue": "2", "wed": "3",
	"thu": "4", "fri": "5", "sat": "6",
}

// Normalize takes a raw crontab expression and returns its canonical form.
// It expands aliases, replaces named months and weekdays with numbers,
// and normalizes whitespace.
func Normalize(expr string) string {
	trimmed := strings.TrimSpace(expr)
	lower := strings.ToLower(trimmed)

	if canonical, ok := aliases[lower]; ok {
		return canonical
	}

	fields := strings.Fields(trimmed)
	if len(fields) != 5 {
		return trimmed
	}

	fields[3] = replaceNames(fields[3], monthNames)
	fields[4] = replaceNames(fields[4], dowNames)

	return strings.Join(fields, " ")
}

// replaceNames substitutes named tokens in a field with their numeric equivalents.
func replaceNames(field string, names map[string]string) string {
	result := field
	for name, num := range names {
		result = strings.ReplaceAll(strings.ToLower(result), name, num)
	}
	return result
}
