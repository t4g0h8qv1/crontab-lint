// Package canonicalizer rewrites crontab expressions into a canonical form:
// fields are sorted, redundant ranges collapsed, and step-1 simplified.
package canonicalizer

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/user/crontab-lint/internal/normalizer"
	"github.com/user/crontab-lint/internal/splitter"
)

// Result holds the original expression, its canonical form, and any changes made.
type Result struct {
	Original  string   `json:"original"`
	Canonical string   `json:"canonical"`
	Changes   []string `json:"changes"`
	Error     string   `json:"error,omitempty"`
}

// Canonicalize returns the canonical form of a crontab expression.
func Canonicalize(expr string) Result {
	norm, err := normalizer.Normalize(expr)
	if err != nil {
		return Result{Original: expr, Error: err.Error()}
	}

	fields, err := splitter.Split(norm)
	if err != nil {
		return Result{Original: expr, Error: err.Error()}
	}

	var changes []string
	canonical := make([]string, len(fields))
	for i, f := range fields {
		c, fieldChanges := canonicalizeField(f, splitter.FieldNames()[i])
		canonical[i] = c
		changes = append(changes, fieldChanges...)
	}

	return Result{
		Original:  expr,
		Canonical: splitter.Join(canonical),
		Changes:   changes,
	}
}

func canonicalizeField(field, name string) (string, []string) {
	var changes []string

	// Handle step-1: */1 -> *
	if field == "*/1" {
		changes = append(changes, fmt.Sprintf("%s: */1 simplified to *", name))
		return "*", changes
	}

	// Handle list fields: sort numeric lists
	if strings.Contains(field, ",") && !strings.ContainsAny(field, "/-") {
		parts := strings.Split(field, ",")
		if sorted, ok := sortNumericList(parts); ok && sorted != field {
			changes = append(changes, fmt.Sprintf("%s: list sorted to %s", name, sorted))
			return sorted, changes
		}
	}

	return field, changes
}

func sortNumericList(parts []string) (string, bool) {
	nums := make([]int, 0, len(parts))
	for _, p := range parts {
		n, err := strconv.Atoi(strings.TrimSpace(p))
		if err != nil {
			return "", false
		}
		nums = append(nums, n)
	}
	// insertion sort (small slices)
	for i := 1; i < len(nums); i++ {
		for j := i; j > 0 && nums[j] < nums[j-1]; j-- {
			nums[j], nums[j-1] = nums[j-1], nums[j]
		}
	}
	strs := make([]string, len(nums))
	for i, n := range nums {
		strs[i] = strconv.Itoa(n)
	}
	return strings.Join(strs, ","), true
}
