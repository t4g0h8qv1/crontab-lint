package parser

import (
	"fmt"
	"strconv"
	"strings"
)

// FieldType represents the type of a crontab field.
type FieldType int

const (
	FieldMinute FieldType = iota
	FieldHour
	FieldDayOfMonth
	FieldMonth
	FieldDayOfWeek
)

// fieldBounds defines the valid min/max for each crontab field.
var fieldBounds = map[FieldType][2]int{
	FieldMinute:     {0, 59},
	FieldHour:       {0, 23},
	FieldDayOfMonth: {1, 31},
	FieldMonth:      {1, 12},
	FieldDayOfWeek:  {0, 7},
}

// fieldNames maps FieldType to a human-readable name.
var fieldNames = map[FieldType]string{
	FieldMinute:     "minute",
	FieldHour:       "hour",
	FieldDayOfMonth: "day of month",
	FieldMonth:      "month",
	FieldDayOfWeek:  "day of week",
}

// ParseField validates a single crontab field value against its type constraints.
// It supports: wildcard (*), values, ranges (1-5), steps (*/2, 1-5/2), and lists (1,2,3).
func ParseField(value string, ft FieldType) error {
	bounds := fieldBounds[ft]
	name := fieldNames[ft]

	for _, part := range strings.Split(value, ",") {
		if err := parsePart(part, bounds[0], bounds[1], name); err != nil {
			return err
		}
	}
	return nil
}

func parsePart(part string, min, max int, name string) error {
	stepParts := strings.SplitN(part, "/", 2)
	base := stepParts[0]

	if len(stepParts) == 2 {
		step, err := strconv.Atoi(stepParts[1])
		if err != nil || step < 1 {
			return fmt.Errorf("%s: invalid step value %q", name, stepParts[1])
		}
	}

	if base == "*" {
		return nil
	}

	if strings.Contains(base, "-") {
		return parseRange(base, min, max, name)
	}

	v, err := strconv.Atoi(base)
	if err != nil {
		return fmt.Errorf("%s: invalid value %q", name, base)
	}
	if v < min || v > max {
		return fmt.Errorf("%s: value %d out of range [%d-%d]", name, v, min, max)
	}
	return nil
}

func parseRange(r string, min, max int, name string) error {
	parts := strings.SplitN(r, "-", 2)
	if len(parts) != 2 {
		return fmt.Errorf("%s: invalid range %q", name, r)
	}
	lo, err1 := strconv.Atoi(parts[0])
	hi, err2 := strconv.Atoi(parts[1])
	if err1 != nil || err2 != nil {
		return fmt.Errorf("%s: invalid range values in %q", name, r)
	}
	if lo < min || hi > max {
		return fmt.Errorf("%s: range %d-%d out of bounds [%d-%d]", name, lo, hi, min, max)
	}
	if lo > hi {
		return fmt.Errorf("%s: range start %d is greater than end %d", name, lo, hi)
	}
	return nil
}
