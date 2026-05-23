// Package truncator shortens crontab expressions by replacing explicit
// full-range lists with wildcards where semantically equivalent.
package truncator

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/user/crontab-lint/internal/splitter"
)

// fieldRange defines the inclusive min/max for each crontab field.
var fieldRange = []struct {
	min, max int
}{
	{0, 59}, // minute
	{0, 23}, // hour
	{1, 31}, // day-of-month
	{1, 12}, // month
	{0, 6},  // day-of-week
}

// Result holds the output of a Truncate call.
type Result struct {
	Original  string
	Truncated string
	Changes   []string
	Errors    []string
}

// Truncate replaces explicit full-range enumerations or step-1 wildcards with
// a bare wildcard (*) for each field that is semantically equivalent.
func Truncate(expr string) Result {
	fields, err := splitter.Split(expr)
	if err != nil {
		return Result{
			Original: expr,
			Errors:   []string{err.Error()},
		}
	}

	names := splitter.FieldNames()
	out := make([]string, len(fields))
	var changes []string

	for i, f := range fields {
		truncated, changed := truncateField(f, fieldRange[i].min, fieldRange[i].max)
		out[i] = truncated
		if changed {
			changes = append(changes, fmt.Sprintf("%s: %q → %q", names[i], f, truncated))
		}
	}

	return Result{
		Original:  expr,
		Truncated: strings.Join(out, " "),
		Changes:   changes,
	}
}

// truncateField returns the simplified form of a single field value.
func truncateField(field string, min, max int) (string, bool) {
	// */1 → *
	if field == "*/1" {
		return "*", true
	}

	// min-max or min-max/1 → *
	base := field
	if idx := strings.Index(field, "/1"); idx > 0 && field[idx+2:] == "" {
		base = field[:idx]
	}
	if lo, hi, ok := parseRangeParts(base); ok {
		if lo == min && hi == max {
			return "*", true
		}
	}

	return field, false
}

func parseRangeParts(s string) (int, int, bool) {
	parts := strings.SplitN(s, "-", 2)
	if len(parts) != 2 {
		return 0, 0, false
	}
	lo, err1 := strconv.Atoi(parts[0])
	hi, err2 := strconv.Atoi(parts[1])
	if err1 != nil || err2 != nil {
		return 0, 0, false
	}
	return lo, hi, true
}
