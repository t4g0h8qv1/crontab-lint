package formatter

import (
	"fmt"
	"strings"
)

// Issue represents a validation or lint issue found in a crontab expression.
type Issue struct {
	Field   string
	Message string
	Severity string // "error" or "warning"
}

// Report holds the full lint report for a crontab expression.
type Report struct {
	Expression  string
	Valid        bool
	Explanation  string
	Issues       []Issue
}

// FormatText returns a human-readable plain-text representation of a Report.
func FormatText(r Report) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Expression : %s\n", r.Expression))
	sb.WriteString(fmt.Sprintf("Valid      : %v\n", r.Valid))

	if r.Explanation != "" {
		sb.WriteString(fmt.Sprintf("Explanation: %s\n", r.Explanation))
	}

	if len(r.Issues) == 0 {
		sb.WriteString("Issues     : none\n")
		return sb.String()
	}

	sb.WriteString("Issues:\n")
	for _, issue := range r.Issues {
		sb.WriteString(fmt.Sprintf("  [%s] %s: %s\n",
			strings.ToUpper(issue.Severity), issue.Field, issue.Message))
	}

	return sb.String()
}

// FormatJSON returns a JSON string representation of a Report.
func FormatJSON(r Report) string {
	issuesJSON := "[]"
	if len(r.Issues) > 0 {
		parts := make([]string, 0, len(r.Issues))
		for _, issue := range r.Issues {
			parts = append(parts, fmt.Sprintf(
				`{"field":%q,"message":%q,"severity":%q}`,
				issue.Field, issue.Message, issue.Severity,
			))
		}
		issuesJSON = "[" + strings.Join(parts, ",") + "]"
	}

	return fmt.Sprintf(
		`{"expression":%q,"valid":%v,"explanation":%q,"issues":%s}`,
		r.Expression, r.Valid, r.Explanation, issuesJSON,
	)
}

// ErrorCount returns the number of issues with severity "error".
func (r Report) ErrorCount() int {
	count := 0
	for _, issue := range r.Issues {
		if issue.Severity == "error" {
			count++
		}
	}
	return count
}

// WarningCount returns the number of issues with severity "warning".
func (r Report) WarningCount() int {
	count := 0
	for _, issue := range r.Issues {
		if issue.Severity == "warning" {
			count++
		}
	}
	return count
}
