// Package linter provides static analysis rules for crontab expressions.
package linter

import (
	"fmt"
	"strings"
)

// Severity represents the severity level of a lint issue.
type Severity string

const (
	SeverityWarning Severity = "warning"
	SeverityError   Severity = "error"
	SeverityInfo    Severity = "info"
)

// Issue represents a single lint finding.
type Issue struct {
	Field    string
	Message  string
	Severity Severity
}

func (i Issue) String() string {
	return fmt.Sprintf("[%s] %s: %s", i.Severity, i.Field, i.Message)
}

// Lint runs all lint rules against a crontab expression and returns any issues found.
func Lint(expr string) []Issue {
	var issues []Issue

	fields := strings.Fields(expr)
	if len(fields) != 5 {
		return []Issue{{
			Field:    "expression",
			Message:  fmt.Sprintf("expected 5 fields, got %d", len(fields)),
			Severity: SeverityError,
		}}
	}

	fieldNames := []string{"minute", "hour", "day-of-month", "month", "day-of-week"}
	for idx, field := range fields {
		issues = append(issues, checkField(fieldNames[idx], field)...)
	}

	issues = append(issues, checkConflicts(fields)...)
	return issues
}

// checkField applies per-field lint rules.
func checkField(name, value string) []Issue {
	var issues []Issue

	if value == "*" {
		return issues
	}

	if strings.Contains(value, "**") {
		issues = append(issues, Issue{
			Field:    name,
			Message:  "double wildcard '**' is not valid",
			Severity: SeverityError,
		})
	}

	if strings.HasPrefix(value, "/") || strings.HasSuffix(value, "/") {
		issues = append(issues, Issue{
			Field:    name,
			Message:  "step expression must be in the form '*/n' or 'a-b/n'",
			Severity: SeverityError,
		})
	}

	if value == "0/1" || value == "*/1" {
		issues = append(issues, Issue{
			Field:    name,
			Message:  fmt.Sprintf("step of 1 in '%s' is redundant; use '*' instead", value),
			Severity: SeverityInfo,
		})
	}

	return issues
}

// checkConflicts detects logical conflicts between fields.
func checkConflicts(fields []string) []Issue {
	var issues []Issue
	dom := fields[2]
	dow := fields[4]

	if dom != "*" && dow != "*" {
		issues = append(issues, Issue{
			Field:    "day-of-month/day-of-week",
			Message:  "specifying both day-of-month and day-of-week may produce unexpected results (OR semantics)",
			Severity: SeverityWarning,
		})
	}
	return issues
}
