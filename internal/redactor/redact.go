// Package redactor replaces specific field values in a cron expression
// with wildcards or neutral placeholders, useful for anonymising or
// generalising expressions before sharing or storing them.
package redactor

import (
	"errors"
	"fmt"
	"strings"
)

// Result holds the outcome of a Redact operation.
type Result struct {
	// Original is the input expression unchanged.
	Original string
	// Redacted is the expression after applying replacements.
	Redacted string
	// Fields lists the per-field values after redaction.
	Fields []string
	// Changed reports which field indices were altered (0-based).
	Changed []int
	// Errors contains any validation problems encountered.
	Errors []string
}

var fieldNames = []string{"minute", "hour", "day-of-month", "month", "day-of-week"}

// Redact replaces the fields identified by the indices in mask with
// replacement. Indices are 0-based (0=minute … 4=day-of-week).
// Pass replacement="*" to generalise the masked fields.
func Redact(expr, replacement string, mask []int) Result {
	parts := strings.Fields(expr)
	if len(parts) != 5 {
		return Result{
			Original: expr,
			Redacted: expr,
			Errors: []string{fmt.Sprintf("expected 5 fields, got %d", len(parts))},
		}
	}

	if err := validateMask(mask); err != nil {
		return Result{
			Original: expr,
			Redacted: expr,
			Errors: []string{err.Error()},
		}
	}

	if replacement == "" {
		replacement = "*"
	}

	result := make([]string, 5)
	copy(result, parts)

	changed := make([]int, 0, len(mask))
	for _, idx := range mask {
		if result[idx] != replacement {
			result[idx] = replacement
			changed = append(changed, idx)
		}
	}

	return Result{
		Original: expr,
		Redacted: strings.Join(result, " "),
		Fields:   result,
		Changed:  changed,
	}
}

func validateMask(mask []int) error {
	seen := make(map[int]bool, len(mask))
	for _, idx := range mask {
		if idx < 0 || idx > 4 {
			return errors.New(fmt.Sprintf("field index %d out of range (0-4)", idx))
		}
		if seen[idx] {
			return errors.New(fmt.Sprintf("duplicate field index %d in mask", idx))
		}
		seen[idx] = true
	}
	return nil
}
