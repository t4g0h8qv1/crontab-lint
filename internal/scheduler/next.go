// Package scheduler computes the next scheduled run time for a cron expression.
package scheduler

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// NextN returns the next n scheduled times after the given reference time
// for the provided cron expression (5-field standard format).
func NextN(expr string, after time.Time, n int) ([]time.Time, error) {
	fields := strings.Fields(expr)
	if len(fields) != 5 {
		return nil, fmt.Errorf("expected 5 fields, got %d", len(fields))
	}

	minutes, err := expand(fields[0], 0, 59)
	if err != nil {
		return nil, fmt.Errorf("minute: %w", err)
	}
	hours, err := expand(fields[1], 0, 23)
	if err != nil {
		return nil, fmt.Errorf("hour: %w", err)
	}
	doms, err := expand(fields[2], 1, 31)
	if err != nil {
		return nil, fmt.Errorf("day-of-month: %w", err)
	}
	months, err := expand(fields[3], 1, 12)
	if err != nil {
		return nil, fmt.Errorf("month: %w", err)
	}
	dows, err := expand(fields[4], 0, 6)
	if err != nil {
		return nil, fmt.Errorf("day-of-week: %w", err)
	}

	var results []time.Time
	t := after.Truncate(time.Minute).Add(time.Minute)

	for len(results) < n {
		if !inSet(months, int(t.Month())) {
			t = advanceToNextMonth(t)
			continue
		}
		domWild := fields[2] == "*"
		dowWild := fields[4] == "*"
		domMatch := inSet(doms, t.Day())
		dowMatch := inSet(dows, int(t.Weekday()))
		var dayMatch bool
		if domWild && dowWild {
			dayMatch = true
		} else if !domWild && !dowWild {
			dayMatch = domMatch || dowMatch
		} else {
			dayMatch = domMatch && dowMatch
		}
		if !dayMatch {
			t = t.Add(24 * time.Hour).Truncate(24 * time.Hour)
			continue
		}
		if !inSet(hours, t.Hour()) {
			t = t.Add(time.Hour).Truncate(time.Hour)
			continue
		}
		if !inSet(minutes, t.Minute()) {
			t = t.Add(time.Minute)
			continue
		}
		results = append(results, t)
		t = t.Add(time.Minute)
	}
	return results, nil
}

// Next returns the single next scheduled time after the given reference time.
func Next(expr string, after time.Time) (time.Time, error) {
	times, err := NextN(expr, after, 1)
	if err != nil {
		return time.Time{}, err
	}
	return times[0], nil
}

func expand(field string, min, max int) (map[int]bool, error) {
	set := make(map[int]bool)
	if field == "*" {
		for i := min; i <= max; i++ {
			set[i] = true
		}
		return set, nil
	}
	for _, part := range strings.Split(field, ",") {
		if strings.Contains(part, "/") {
			sub := strings.SplitN(part, "/", 2)
			step, err := strconv.Atoi(sub[1])
			if err != nil || step <= 0 {
				return nil, fmt.Errorf("invalid step in %q", part)
			}
			start := min
			if sub[0] != "*" {
				start, err = strconv.Atoi(sub[0])
				if err != nil {
					return nil, fmt.Errorf("invalid range start in %q", part)
				}
			}
			for i := start; i <= max; i += step {
				set[i] = true
			}
		} else if strings.Contains(part, "-") {
			sub := strings.SplitN(part, "-", 2)
			lo, err1 := strconv.Atoi(sub[0])
			hi, err2 := strconv.Atoi(sub[1])
			if err1 != nil || err2 != nil {
				return nil, fmt.Errorf("invalid range %q", part)
			}
			for i := lo; i <= hi; i++ {
				set[i] = true
			}
		} else {
			v, err := strconv.Atoi(part)
			if err != nil {
				return nil, fmt.Errorf("invalid value %q", part)
			}
			set[v] = true
		}
	}
	return set, nil
}

func inSet(set map[int]bool, v int) bool { return set[v] }

func advanceToNextMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month()+1, 1, 0, 0, 0, 0, t.Location())
}
