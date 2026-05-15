// Package differ compares two crontab expressions and reports their differences.
package differ

import (
	"fmt"
	"strings"

	"github.com/user/crontab-lint/internal/normalizer"
)

// FieldNames holds the canonical names for each crontab field position.
var FieldNames = []string{"minute", "hour", "day-of-month", "month", "day-of-week"}

// FieldDiff describes a difference between a single field in two expressions.
type FieldDiff struct {
	Field string
	From  string
	To    string
}

// Result holds the full diff between two crontab expressions.
type Result struct {
	From    string
	To      string
	Changes []FieldDiff
	Equal   bool
}

// Diff compares two crontab expressions field-by-field after normalizing both.
// Returns an error if either expression has an invalid field count.
func Diff(from, to string) (*Result, error) {
	normFrom, err := normalizer.Normalize(from)
	if err != nil {
		return nil, fmt.Errorf("invalid 'from' expression: %w", err)
	}
	normTo, err := normalizer.Normalize(to)
	if err != nil {
		return nil, fmt.Errorf("invalid 'to' expression: %w", err)
	}

	fieldsFrom := strings.Fields(normFrom)
	fieldsTo := strings.Fields(normTo)

	result := &Result{
		From: from,
		To:   to,
	}

	for i := 0; i < len(FieldNames); i++ {
		f := fieldsFrom[i]
		t := fieldsTo[i]
		if f != t {
			result.Changes = append(result.Changes, FieldDiff{
				Field: FieldNames[i],
				From:  f,
				To:    t,
			})
		}
	}

	result.Equal = len(result.Changes) == 0
	return result, nil
}

// Summary returns a short human-readable summary of the diff result.
func Summary(r *Result) string {
	if r.Equal {
		return "Expressions are equivalent."
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "Expressions differ in %d field(s):\n", len(r.Changes))
	for _, c := range r.Changes {
		fmt.Fprintf(&sb, "  %-14s %s  →  %s\n", c.Field+":", c.From, c.To)
	}
	return strings.TrimRight(sb.String(), "\n")
}
