// Package trimmer simplifies crontab expressions by removing redundant
// components while preserving their semantic meaning.
package trimmer

import (
	"fmt"
	"strings"

	"github.com/user/crontab-lint/internal/splitter"
)

// Result holds the trimmed expression and a log of changes made.
type Result struct {
	Original string
	Trimmed  string
	Changes  []string
}

// Trim takes a crontab expression and returns a simplified version
// with redundant constructs removed.
func Trim(expr string) (Result, error) {
	fields, err := splitter.Split(expr)
	if err != nil {
		return Result{}, fmt.Errorf("trim: %w", err)
	}

	result := Result{Original: expr}
	out := make([]string, len(fields))
	copy(out, fields)

	for i, f := range out {
		trimmed, changes := trimField(f, i)
		out[i] = trimmed
		result.Changes = append(result.Changes, changes...)
	}

	result.Trimmed = strings.Join(out, " ")
	return result, nil
}

// trimField simplifies a single crontab field value.
func trimField(field string, index int) (string, []string) {
	name := splitter.FieldNames()[index]
	var changes []string

	// */1 -> *
	if field == "*/1" {
		changes = append(changes, fmt.Sprintf("%s: */1 simplified to *", name))
		return "*", changes
	}

	// Handle comma-separated list with a single element
	parts := strings.Split(field, ",")
	if len(parts) == 1 && !strings.Contains(field, "-") && !strings.Contains(field, "/") && !strings.Contains(field, "*") {
		return field, changes
	}

	// Deduplicate comma-separated values
	seen := make(map[string]bool)
	unique := parts[:0]
	for _, p := range parts {
		if !seen[p] {
			seen[p] = true
			unique = append(unique, p)
		}
	}
	if len(unique) < len(parts) {
		changes = append(changes, fmt.Sprintf("%s: removed duplicate values", name))
		field = strings.Join(unique, ",")
	}

	return field, changes
}
