// Package matcher provides functionality to check whether a cron expression
// matches a given time instant.
package matcher

import (
	"fmt"
	"time"

	"github.com/user/crontab-lint/internal/normalizer"
	"github.com/user/crontab-lint/internal/parser"
)

// Result holds the outcome of a match check.
type Result struct {
	Expression string
	Time       time.Time
	Matches    bool
	Fields     [5]FieldMatch
}

// FieldMatch describes whether a single cron field matched.
type FieldMatch struct {
	Name    string
	Value   int
	Matches bool
}

var fieldNames = [5]string{"minute", "hour", "day-of-month", "month", "day-of-week"}

// Match reports whether the given cron expression fires at the given time.
func Match(expr string, t time.Time) (Result, error) {
	norm, err := normalizer.Normalize(expr)
	if err != nil {
		return Result{}, fmt.Errorf("normalize: %w", err)
	}

	fields, err := splitFields(norm)
	if err != nil {
		return Result{}, err
	}

	values := [5]int{
		t.Minute(),
		t.Hour(),
		t.Day(),
		int(t.Month()),
		int(t.Weekday()),
	}

	ranges := [5][2]int{
		{0, 59},
		{0, 23},
		{1, 31},
		{1, 12},
		{0, 6},
	}

	res := Result{Expression: expr, Time: t}
	allMatch := true

	for i, f := range fields {
		set, err := parser.ParseField(f, ranges[i][0], ranges[i][1])
		if err != nil {
			return Result{}, fmt.Errorf("field %s: %w", fieldNames[i], err)
		}
		ok := contains(set, values[i])
		res.Fields[i] = FieldMatch{Name: fieldNames[i], Value: values[i], Matches: ok}
		if !ok {
			allMatch = false
		}
	}

	res.Matches = allMatch
	return res, nil
}

func contains(set []int, v int) bool {
	for _, s := range set {
		if s == v {
			return true
		}
	}
	return false
}

func splitFields(expr string) ([]string, error) {
	fields := make([]string, 0, 5)
	cur := ""
	for _, ch := range expr {
		if ch == ' ' || ch == '\t' {
			if cur != "" {
				fields = append(fields, cur)
				cur = ""
			}
		} else {
			cur += string(ch)
		}
	}
	if cur != "" {
		fields = append(fields, cur)
	}
	if len(fields) != 5 {
		return nil, fmt.Errorf("expected 5 fields, got %d", len(fields))
	}
	return fields, nil
}
