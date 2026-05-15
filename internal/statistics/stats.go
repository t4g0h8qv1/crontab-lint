// Package statistics provides frequency and distribution analysis
// for crontab expressions, summarising how often a schedule fires.
package statistics

import (
	"fmt"
	"time"

	"github.com/user/crontab-lint/internal/scheduler"
)

// Stats holds computed statistics for a crontab expression.
type Stats struct {
	// Expression is the original crontab string.
	Expression string
	// FiresPerDay is the average number of times the job fires per day.
	FiresPerDay float64
	// FiresPerHour is the average number of times the job fires per hour.
	FiresPerHour float64
	// FiresPerWeek is the average number of times the job fires per week.
	FiresPerWeek float64
	// MinInterval is the shortest gap between two consecutive firings.
	MinInterval time.Duration
	// MaxInterval is the longest gap between two consecutive firings.
	MaxInterval time.Duration
	// SampleSize is the number of future occurrences used to compute stats.
	SampleSize int
}

// Compute analyses the next sampleSize occurrences of expr starting from
// base and returns a populated Stats value. It returns an error if the
// expression cannot be scheduled.
func Compute(expr string, base time.Time, sampleSize int) (Stats, error) {
	if sampleSize < 2 {
		return Stats{}, fmt.Errorf("statistics: sampleSize must be at least 2, got %d", sampleSize)
	}

	times, err := scheduler.NextN(expr, base, sampleSize)
	if err != nil {
		return Stats{}, fmt.Errorf("statistics: %w", err)
	}
	if len(times) < 2 {
		return Stats{}, fmt.Errorf("statistics: could not compute enough occurrences")
	}

	var minGap, maxGap time.Duration
	var totalGap time.Duration
	for i := 1; i < len(times); i++ {
		gap := times[i].Sub(times[i-1])
		totalGap += gap
		if i == 1 || gap < minGap {
			minGap = gap
		}
		if gap > maxGap {
			maxGap = gap
		}
	}

	intervals := len(times) - 1
	avgGap := totalGap / time.Duration(intervals)

	firesPerHour := float64(time.Hour) / float64(avgGap)
	firesPerDay := firesPerHour * 24
	firesPerWeek := firesPerDay * 7

	return Stats{
		Expression:  expr,
		FiresPerDay: firesPerDay,
		FiresPerHour: firesPerHour,
		FiresPerWeek: firesPerWeek,
		MinInterval: minGap,
		MaxInterval: maxGap,
		SampleSize:  len(times),
	}, nil
}
