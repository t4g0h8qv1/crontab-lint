// Package comparator provides schedule frequency comparison between two cron expressions.
package comparator

import (
	"fmt"
	"time"

	"github.com/user/crontab-lint/internal/scheduler"
)

// Result holds the comparison outcome between two cron expressions.
type Result struct {
	ExpressionA  string
	ExpressionB  string
	FrequencyA   int // runs per day
	FrequencyB   int // runs per day
	Delta        int // FrequencyB - FrequencyA
	Relation     string // "faster", "slower", "equal"
	RatioPercent float64
}

// Compare computes how often each expression fires in a 24-hour window
// starting from the given base time, then returns a Result describing
// the relative frequency difference.
func Compare(exprA, exprB string, base time.Time) (Result, error) {
	countA, err := countFirings(exprA, base)
	if err != nil {
		return Result{}, fmt.Errorf("expression A: %w", err)
	}

	countB, err := countFirings(exprB, base)
	if err != nil {
		return Result{}, fmt.Errorf("expression B: %w", err)
	}

	result := Result{
		ExpressionA: exprA,
		ExpressionB: exprB,
		FrequencyA:  countA,
		FrequencyB:  countB,
		Delta:       countB - countA,
	}

	switch {
	case countA == 0 && countB == 0:
		result.Relation = "equal"
		result.RatioPercent = 0
	case countA == 0:
		result.Relation = "faster"
		result.RatioPercent = 100
	case countB == countA:
		result.Relation = "equal"
		result.RatioPercent = 0
	case countB > countA:
		result.Relation = "faster"
		result.RatioPercent = float64(countB-countA) / float64(countA) * 100
	default:
		result.Relation = "slower"
		result.RatioPercent = float64(countA-countB) / float64(countA) * 100
	}

	return result, nil
}

// countFirings returns how many times expr fires within 24 hours after base.
func countFirings(expr string, base time.Time) (int, error) {
	end := base.Add(24 * time.Hour)
	times, err := scheduler.NextN(expr, base, 1500) // generous upper bound
	if err != nil {
		return 0, err
	}
	count := 0
	for _, t := range times {
		if t.After(end) {
			break
		}
		count++
	}
	return count, nil
}
