// Package paraphraser rewrites a cron expression in plain English,
// producing a concise human-readable sentence for each valid expression.
package paraphraser

import (
	"fmt"
	"strings"

	"github.com/user/crontab-lint/internal/normalizer"
	"github.com/user/crontab-lint/internal/splitter"
)

// Result holds the paraphrased output for a cron expression.
type Result struct {
	Expression  string `json:"expression"`
	Paraphrase  string `json:"paraphrase"`
	Errors      []string `json:"errors,omitempty"`
}

// Paraphrase converts a cron expression into a plain-English sentence.
// It normalises aliases before processing so that e.g. "@daily" and
// "0 0 * * *" produce the same output.
func Paraphrase(expr string) Result {
	norm, err := normalizer.Normalize(expr)
	if err != nil {
		return Result{
			Expression: expr,
			Errors:     []string{err.Error()},
		}
	}

	fields, err := splitter.Split(norm)
	if err != nil {
		return Result{
			Expression: expr,
			Errors:     []string{err.Error()},
		}
	}

	sentence := buildSentence(fields)
	return Result{
		Expression: expr,
		Paraphrase: sentence,
	}
}

func buildSentence(f splitter.Fields) string {
	minute := describeField(f.Minute, "minute", 0, 59)
	hour := describeField(f.Hour, "hour", 0, 23)
	dom := describeField(f.Dom, "day-of-month", 1, 31)
	month := describeField(f.Month, "month", 1, 12)
	dow := describeField(f.Dow, "day-of-week", 0, 6)

	parts := []string{}
	if minute != "" {
		parts = append(parts, "at "+minute)
	}
	if hour != "" {
		parts = append(parts, "past "+hour)
	}
	if dom != "" {
		parts = append(parts, "on "+dom)
	}
	if month != "" {
		parts = append(parts, "in "+month)
	}
	if dow != "" {
		parts = append(parts, "on "+dow)
	}

	if len(parts) == 0 {
		return "every minute"
	}
	return strings.Join(parts, ", ")
}

func describeField(value, name string, min, max int) string {
	if value == "*" {
		return ""
	}
	if strings.HasPrefix(value, "*/") {
		step := strings.TrimPrefix(value, "*/")
		if step == "1" {
			return ""
		}
		return fmt.Sprintf("every %s %ss", step, name)
	}
	if strings.Contains(value, "-") && strings.Contains(value, "/") {
		return fmt.Sprintf("%s %s", value, name)
	}
	if strings.Contains(value, "-") {
		return fmt.Sprintf("%s %ss", value, name)
	}
	if strings.Contains(value, ",") {
		return fmt.Sprintf("%s %ss", value, name)
	}
	return fmt.Sprintf("%s %s", value, name)
}
