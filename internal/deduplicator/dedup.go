// Package deduplicator removes duplicate fields and redundant entries
// from crontab expressions, returning a canonical minimal form.
package deduplicator

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// Result holds the output of a deduplication pass.
type Result struct {
	// Original is the input expression unchanged.
	Original string
	// Deduped is the cleaned expression.
	Deduped string
	// Changes lists human-readable descriptions of what was removed.
	Changes []string
	// Changed is true when Deduped differs from Original.
	Changed bool
}

// Dedup removes duplicate values and redundant list entries from each
// field of expr. It returns a Result describing any changes made.
// expr must have exactly 5 space-separated fields.
func Dedup(expr string) (Result, error) {
	fields := strings.Fields(expr)
	if len(fields) != 5 {
		return Result{}, fmt.Errorf("expected 5 fields, got %d", len(fields))
	}

	result := Result{Original: expr}
	cleaned := make([]string, 5)

	for i, f := range fields {
		deduped, changes := dedupField(f)
		cleaned[i] = deduped
		for _, c := range changes {
			result.Changes = append(result.Changes, fmt.Sprintf("field %d: %s", i+1, c))
		}
	}

	result.Deduped = strings.Join(cleaned, " ")
	result.Changed = result.Deduped != expr
	return result, nil
}

// dedupField removes duplicate numeric values from a comma-separated list.
// Non-list tokens (wildcards, ranges, steps) are passed through unchanged.
func dedupField(field string) (string, []string) {
	if !strings.Contains(field, ",") {
		return field, nil
	}

	parts := strings.Split(field, ",")
	seen := make(map[string]bool)
	var kept []string
	var changes []string

	for _, p := range parts {
		norm := normalizeValue(p)
		if seen[norm] {
			changes = append(changes, fmt.Sprintf("removed duplicate %q", p))
			continue
		}
		seen[norm] = true
		kept = append(kept, p)
	}

	// Sort purely numeric lists for canonical output.
	if allNumeric(kept) {
		sort.Slice(kept, func(i, j int) bool {
			a, _ := strconv.Atoi(kept[i])
			b, _ := strconv.Atoi(kept[j])
			return a < b
		})
	}

	return strings.Join(kept, ","), changes
}

func normalizeValue(s string) string { return strings.TrimSpace(strings.ToLower(s)) }

func allNumeric(parts []string) bool {
	for _, p := range parts {
		if _, err := strconv.Atoi(p); err != nil {
			return false
		}
	}
	return true
}
