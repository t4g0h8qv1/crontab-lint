// Package ranker scores crontab expressions by complexity and readability.
package ranker

import (
	"strings"

	"github.com/user/crontab-lint/internal/normalizer"
	"github.com/user/crontab-lint/internal/parser"
)

// Score holds the ranking result for a crontab expression.
type Score struct {
	Expression string `json:"expression"`
	Total      int    `json:"total"`
	Complexity int    `json:"complexity"`
	Readability int   `json:"readability"`
	Notes      []string `json:"notes"`
}

const fieldNames = "minute hour dom month dow"

var fields = strings.Fields(fieldNames)

// Rank computes a readability/complexity score for the given crontab expression.
// Higher total score means more readable and less complex.
func Rank(expr string) (Score, error) {
	norm, err := normalizer.Normalize(expr)
	if err != nil {
		return Score{}, err
	}

	parts := strings.Fields(norm)
	if len(parts) != 5 {
		return Score{}, fmt.Errorf("expected 5 fields, got %d", len(parts))
	}

	var notes []string
	complexity := 0
	readability := 10

	for i, part := range parts {
		pf, err := parser.ParseField(part, fields[i])
		if err != nil {
			notes = append(notes, "unparseable field: "+fields[i])
			complexity += 3
			continue
		}
		fc := fieldComplexity(part, pf)
		complexity += fc
		if fc == 0 {
			readability++
		} else if fc >= 3 {
			readability--
			notes = append(notes, "complex field: "+fields[i])
		}
	}

	if readability < 0 {
		readability = 0
	}

	total := readability - complexity/2
	if total < 0 {
		total = 0
	}

	return Score{
		Expression:  expr,
		Total:       total,
		Complexity:  complexity,
		Readability: readability,
		Notes:       notes,
	}, nil
}

// fieldComplexity returns a complexity score for a single parsed field.
func fieldComplexity(raw string, _ interface{}) int {
	if raw == "*" {
		return 0
	}
	score := 0
	if strings.Contains(raw, ",") {
		score += len(strings.Split(raw, ",")) - 1
	}
	if strings.Contains(raw, "/") {
		score++
	}
	if strings.Contains(raw, "-") {
		score++
	}
	return score
}
