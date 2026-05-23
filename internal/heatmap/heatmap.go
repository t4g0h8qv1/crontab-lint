// Package heatmap generates a firing frequency heatmap for a cron expression,
// showing how often each hour-of-day and day-of-week combination is triggered.
package heatmap

import (
	"fmt"
	"time"

	"github.com/your-org/crontab-lint/internal/scheduler"
)

// Cell holds the firing count for a specific (day-of-week, hour) pair.
type Cell struct {
	DOW   int // 0=Sunday … 6=Saturday
	Hour  int
	Count int
}

// Result is the output of a heatmap computation.
type Result struct {
	Expression string
	Cells      []Cell // 7 days × 24 hours = 168 entries
	MaxCount   int
	Errors     []string
}

// Compute builds a heatmap by sampling the next sampleSize firings of expr
// starting from origin. A sampleSize of 0 uses the default (10 080 = one week
// at per-minute granularity).
func Compute(expr string, origin time.Time, sampleSize int) Result {
	if sampleSize <= 0 {
		sampleSize = 10_080
	}

	times, err := scheduler.NextN(expr, origin, sampleSize)
	if err != nil {
		return Result{
			Expression: expr,
			Errors:     []string{fmt.Sprintf("scheduler error: %v", err)},
		}
	}

	// grid[dow][hour] = count
	var grid [7][24]int
	for _, t := range times {
		dow := int(t.Weekday()) // time.Weekday: 0=Sunday
		grid[dow][t.Hour()]++
	}

	cells := make([]Cell, 0, 7*24)
	maxCount := 0
	for d := 0; d < 7; d++ {
		for h := 0; h < 24; h++ {
			c := grid[d][h]
			if c > maxCount {
				maxCount = c
			}
			cells = append(cells, Cell{DOW: d, Hour: h, Count: c})
		}
	}

	return Result{
		Expression: expr,
		Cells:      cells,
		MaxCount:   maxCount,
	}
}
