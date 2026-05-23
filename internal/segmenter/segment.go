// Package segmenter breaks a cron expression into labeled segments
// with positional byte offsets, useful for editor integrations and
// hover-based documentation tools.
package segmenter

import (
	"fmt"
	"strings"
)

// Segment represents a single cron field with its name, value, and
// the byte range [Start, End) it occupies in the original expression.
type Segment struct {
	Index int    // 0-based field index
	Name  string // e.g. "minute", "hour"
	Value string // raw field text
	Start int    // byte offset in original expression
	End   int    // exclusive byte offset
}

var fieldNames = []string{"minute", "hour", "day-of-month", "month", "day-of-week"}

// Segment parses expr and returns one Segment per field.
// Returns an error if the expression does not contain exactly 5 fields.
func Segment(expr string) ([]Segment, error) {
	expr = strings.TrimSpace(expr)
	parts := strings.Fields(expr)
	if len(parts) != 5 {
		return nil, fmt.Errorf("segmenter: expected 5 fields, got %d", len(parts))
	}

	segments := make([]Segment, 0, 5)
	offset := 0
	for i, part := range parts {
		// Advance offset to the start of this token.
		for offset < len(expr) && expr[offset] != part[0] {
			offset++
		}
		start := offset
		end := start + len(part)
		segments = append(segments, Segment{
			Index: i,
			Name:  fieldNames[i],
			Value: part,
			Start: start,
			End:   end,
		})
		offset = end
	}
	return segments, nil
}

// AtOffset returns the Segment that contains the given byte offset, or
// nil if the offset falls outside all segments (e.g. on whitespace).
func AtOffset(segments []Segment, offset int) *Segment {
	for i := range segments {
		if offset >= segments[i].Start && offset < segments[i].End {
			return &segments[i]
		}
	}
	return nil
}
