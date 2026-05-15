// Package converter translates crontab expressions into alternative
// scheduling formats such as systemd OnCalendar strings.
package converter

import (
	"fmt"
	"strings"

	"github.com/user/crontab-lint/internal/normalizer"
	"github.com/user/crontab-lint/internal/validator"
)

// Result holds the output of a conversion attempt.
type Result struct {
	Expression string `json:"expression"`
	Format     string `json:"format"`
	Output     string `json:"output,omitempty"`
	Errors     []string `json:"errors,omitempty"`
}

// ToSystemd converts a crontab expression to a systemd OnCalendar string.
// Returns a Result containing the converted value or any errors encountered.
func ToSystemd(expr string) Result {
	norm := normalizer.Normalize(expr)
	errs := validator.Validate(norm)
	if len(errs) > 0 {
		return Result{Expression: expr, Format: "systemd", Errors: errs}
	}

	fields := strings.Fields(norm)
	if len(fields) != 5 {
		return Result{
			Expression: expr,
			Format:     "systemd",
			Errors:     []string{"expression must have exactly 5 fields"},
		}
	}

	minute := fields[0]
	hour := fields[1]
	dom := fields[2]
	month := fields[3]
	dow := fields[4]

	// systemd OnCalendar format: DayOfWeek Year-Month-Day Hour:Minute:Second
	dowPart := convertDow(dow)
	datePart := fmt.Sprintf("*-%s-%s", month, dom)
	timePart := fmt.Sprintf("%s:%s:00", hour, minute)

	var parts []string
	if dowPart != "" {
		parts = append(parts, dowPart)
	}
	parts = append(parts, datePart, timePart)

	return Result{
		Expression: expr,
		Format:     "systemd",
		Output:     strings.Join(parts, " "),
	}
}

// convertDow maps a cron day-of-week field to systemd day abbreviations.
// Returns an empty string when the field is a wildcard.
func convertDow(dow string) string {
	if dow == "*" {
		return ""
	}
	dowNames := map[string]string{
		"0": "Sun", "1": "Mon", "2": "Tue", "3": "Wed",
		"4": "Thu", "5": "Fri", "6": "Sat", "7": "Sun",
	}
	parts := strings.Split(dow, ",")
	var out []string
	for _, p := range parts {
		if name, ok := dowNames[p]; ok {
			out = append(out, name)
		} else {
			out = append(out, p)
		}
	}
	return strings.Join(out, ",")
}
