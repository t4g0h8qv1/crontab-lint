// Package splitter provides utilities for splitting a crontab expression
// into its individual fields and reconstructing expressions from parts.
package splitter

import (
	"errors"
	"strings"
)

// ErrInvalidFieldCount is returned when the expression does not contain
// exactly 5 space-separated fields.
var ErrInvalidFieldCount = errors.New("cron expression must have exactly 5 fields")

// Fields holds the five named fields of a cron expression.
type Fields struct {
	Minute     string
	Hour       string
	DayOfMonth string
	Month      string
	DayOfWeek  string
}

// Split parses a cron expression string into a Fields struct.
// It returns ErrInvalidFieldCount if the expression does not have exactly 5 fields.
func Split(expr string) (Fields, error) {
	parts := strings.Fields(expr)
	if len(parts) != 5 {
		return Fields{}, ErrInvalidFieldCount
	}
	return Fields{
		Minute:     parts[0],
		Hour:       parts[1],
		DayOfMonth: parts[2],
		Month:      parts[3],
		DayOfWeek:  parts[4],
	}, nil
}

// Join reconstructs a cron expression string from a Fields struct.
func Join(f Fields) string {
	return strings.Join([]string{
		f.Minute,
		f.Hour,
		f.DayOfMonth,
		f.Month,
		f.DayOfWeek,
	}, " ")
}

// Slice returns the fields as an ordered string slice:
// [minute, hour, day-of-month, month, day-of-week].
func Slice(f Fields) []string {
	return []string{
		f.Minute,
		f.Hour,
		f.DayOfMonth,
		f.Month,
		f.DayOfWeek,
	}
}

// FieldNames returns the canonical names for each cron field position.
func FieldNames() []string {
	return []string{"minute", "hour", "day-of-month", "month", "day-of-week"}
}
