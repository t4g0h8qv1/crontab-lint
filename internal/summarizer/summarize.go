// Package summarizer produces a concise one-line human-readable summary
// of a cron expression, suitable for display in terminals or log output.
package summarizer

import (
	"fmt"
	"strings"

	"github.com/user/crontab-lint/internal/classifier"
	"github.com/user/crontab-lint/internal/normalizer"
	"github.com/user/crontab-lint/internal/validator"
)

// Result holds the summary output for a cron expression.
type Result struct {
	Expression string `json:"expression"`
	Normalized string `json:"normalized"`
	Class      string `json:"class"`
	Summary    string `json:"summary"`
	Errors     []string `json:"errors,omitempty"`
}

// Summarize validates, normalizes and produces a human-readable summary
// of the given cron expression.
func Summarize(expr string) Result {
	res := Result{Expression: expr}

	errs := validator.Validate(expr)
	if len(errs) > 0 {
		for _, e := range errs {
			res.Errors = append(res.Errors, e.Error())
		}
		res.Summary = "invalid expression"
		return res
	}

	norm, err := normalizer.Normalize(expr)
	if err != nil {
		res.Errors = []string{err.Error()}
		res.Summary = "could not normalize expression"
		return res
	}
	res.Normalized = norm

	cls := classifier.Classify(norm)
	res.Class = cls
	res.Summary = buildSummary(norm, cls)
	return res
}

// buildSummary constructs a descriptive sentence from the normalized expression
// and its classification class.
func buildSummary(norm, cls string) string {
	fields := strings.Fields(norm)
	if len(fields) != 5 {
		return "unknown schedule"
	}
	min, hour, dom, month, dow := fields[0], fields[1], fields[2], fields[3], fields[4]

	switch cls {
	case "every-minute":
		return "runs every minute"
	case "hourly":
		return fmt.Sprintf("runs every hour at minute %s", min)
	case "daily":
		return fmt.Sprintf("runs daily at %s:%s", pad(hour), pad(min))
	case "weekly":
		return fmt.Sprintf("runs weekly on day %s at %s:%s", dow, pad(hour), pad(min))
	case "monthly":
		return fmt.Sprintf("runs monthly on day %s at %s:%s", dom, pad(hour), pad(min))
	default:
		return fmt.Sprintf("runs on schedule: min=%s hour=%s dom=%s month=%s dow=%s",
			min, hour, dom, month, dow)
	}
}

func pad(s string) string {
	if len(s) == 1 {
		return "0" + s
	}
	return s
}
