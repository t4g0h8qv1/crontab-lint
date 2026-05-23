// Package expander expands compact crontab expressions into their full
// enumerated form, listing every discrete value for each field.
package expander

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/user/crontab-lint/internal/parser"
)

// fieldMeta holds the allowed min/max for each cron field.
var fieldMeta = []struct {
	name string
	min  int
	max  int
}{
	{"minute", 0, 59},
	{"hour", 0, 23},
	{"dom", 1, 31},
	{"month", 1, 12},
	{"dow", 0, 6},
}

// Result holds the expanded values for a single crontab expression.
type Result struct {
	Expression string
	Fields     []ExpandedField
	Errors     []string
}

// ExpandedField contains the name and sorted list of concrete values
// that a single cron field resolves to.
type ExpandedField struct {
	Name   string
	Values []int
}

// Expand parses expr and returns the full set of integers each field
// matches within its legal range.
func Expand(expr string) Result {
	res := Result{Expression: expr}

	parts := strings.Fields(expr)
	if len(parts) != 5 {
		res.Errors = append(res.Errors, fmt.Sprintf("expected 5 fields, got %d", len(parts)))
		return res
	}

	for i, part := range parts {
		meta := fieldMeta[i]
		field, err := parser.ParseField(part, meta.min, meta.max)
		if err != nil {
			res.Errors = append(res.Errors, fmt.Sprintf("%s: %v", meta.name, err))
			continue
		}
		values := enumerate(field, meta.min, meta.max)
		res.Fields = append(res.Fields, ExpandedField{Name: meta.name, Values: values})
	}
	return res
}

// enumerate converts a parsed field into a sorted slice of matching integers.
func enumerate(field []parser.Range, min, max int) []int {
	seen := make(map[int]struct{})
	for _, r := range field {
		start := r.Start
		end := r.End
		step := r.Step
		if step == 0 {
			step = 1
		}
		for v := start; v <= end; v += step {
			if v >= min && v <= max {
				seen[v] = struct{}{}
			}
		}
	}
	out := make([]int, 0, len(seen))
	for v := range seen {
		out = append(out, v)
	}
	sort.Ints(out)
	return out
}

// JoinValues formats a slice of ints as a comma-separated string.
func JoinValues(vals []int) string {
	strs := make([]string, len(vals))
	for i, v := range vals {
		strs[i] = strconv.Itoa(v)
	}
	return strings.Join(strs, ",")
}
