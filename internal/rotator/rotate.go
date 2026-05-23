// Package rotator shifts crontab field values by a given offset,
// wrapping within the valid range for each field.
package rotator

import (
	"fmt"
	"strconv"
	"strings"
)

// fieldMeta holds the valid min/max for each cron field.
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

// Result holds the outcome of a Rotate call.
type Result struct {
	Original string
	Rotated  string
	Offsets  [5]int
	Changed  []string
	Errors   []string
}

// Rotate shifts each field's numeric values by the corresponding offset,
// wrapping within the field's valid range. Pass 0 to leave a field unchanged.
func Rotate(expr string, offsets [5]int) Result {
	fields := strings.Fields(expr)
	result := Result{Original: expr, Offsets: offsets}

	if len(fields) != 5 {
		result.Errors = append(result.Errors,
			fmt.Sprintf("expected 5 fields, got %d", len(fields)))
		result.Rotated = expr
		return result
	}

	out := make([]string, 5)
	for i, f := range fields {
		rotated, changed, err := rotateField(f, offsets[i], fieldMeta[i].min, fieldMeta[i].max)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("%s: %v", fieldMeta[i].name, err))
			out[i] = f
		} else {
			out[i] = rotated
			if changed {
				result.Changed = append(result.Changed, fieldMeta[i].name)
			}
		}
	}

	result.Rotated = strings.Join(out, " ")
	return result
}

// rotateField shifts a single field token by offset within [min, max].
func rotateField(field string, offset, min, max int) (string, bool, error) {
	if offset == 0 || field == "*" {
		return field, false, nil
	}
	span := max - min + 1
	parts := strings.Split(field, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		v, err := strconv.Atoi(strings.TrimSpace(p))
		if err != nil {
			return field, false, fmt.Errorf("non-numeric value %q cannot be rotated", p)
		}
		shifted := ((v-min+offset)%span+span)%span + min
		result = append(result, strconv.Itoa(shifted))
	}
	return strings.Join(result, ","), true, nil
}
