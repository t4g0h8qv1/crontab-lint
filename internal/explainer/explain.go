package explainer

import (
	"fmt"
	"strings"
)

// Explanation holds human-readable descriptions for each crontab field.
type Explanation struct {
	Minute     string
	Hour       string
	DayOfMonth string
	Month      string
	DayOfWeek  string
	Summary    string
}

var monthNames = []string{
	"", "January", "February", "March", "April", "May", "June",
	"July", "August", "September", "October", "November", "December",
}

var dayNames = []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}

// ExplainField returns a human-readable description for a single crontab field token.
func ExplainField(token, fieldName string, min, max int) string {
	if token == "*" {
		return fmt.Sprintf("every %s", fieldName)
	}
	if strings.HasPrefix(token, "*/") {
		step := strings.TrimPrefix(token, "*/")
		return fmt.Sprintf("every %s %s(s)", step, fieldName)
	}
	if strings.Contains(token, "-") {
		parts := strings.SplitN(token, "-", 2)
		return fmt.Sprintf("%s from %s to %s", fieldName, parts[0], parts[1])
	}
	if strings.Contains(token, ",") {
		return fmt.Sprintf("%s %s", fieldName, token)
	}
	return fmt.Sprintf("%s %s", fieldName, token)
}

// Explain produces a full Explanation for a 5-field crontab expression.
func Explain(expr string) (Explanation, error) {
	parts := strings.Fields(expr)
	if len(parts) != 5 {
		return Explanation{}, fmt.Errorf("expected 5 fields, got %d", len(parts))
	}

	ex := Explanation{
		Minute:     ExplainField(parts[0], "minute", 0, 59),
		Hour:       ExplainField(parts[1], "hour", 0, 23),
		DayOfMonth: ExplainField(parts[2], "day-of-month", 1, 31),
		Month:      ExplainField(parts[3], "month", 1, 12),
		DayOfWeek:  ExplainField(parts[4], "day-of-week", 0, 7),
	}
	ex.Summary = fmt.Sprintf("Run at %s, %s, %s, %s, %s.",
		ex.Minute, ex.Hour, ex.DayOfMonth, ex.Month, ex.DayOfWeek)
	return ex, nil
}
