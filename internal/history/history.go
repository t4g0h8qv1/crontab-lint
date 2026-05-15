// Package history provides functionality for recording and retrieving
// previously validated crontab expressions within a session.
package history

import (
	"fmt"
	"strings"
	"time"
)

// Entry represents a single history record of a validated crontab expression.
type Entry struct {
	Expression string
	Timestamp  time.Time
	Valid      bool
	Issues     []string
}

// History holds a bounded list of crontab validation entries.
type History struct {
	entries []Entry
	maxSize int
}

// New creates a new History with the given maximum capacity.
func New(maxSize int) *History {
	if maxSize <= 0 {
		maxSize = 50
	}
	return &History{maxSize: maxSize}
}

// Add appends a new entry to the history, evicting the oldest if at capacity.
func (h *History) Add(expr string, valid bool, issues []string) {
	if len(h.entries) >= h.maxSize {
		h.entries = h.entries[1:]
	}
	h.entries = append(h.entries, Entry{
		Expression: expr,
		Timestamp:  time.Now(),
		Valid:      valid,
		Issues:     issues,
	})
}

// All returns a copy of all history entries in chronological order.
func (h *History) All() []Entry {
	result := make([]Entry, len(h.entries))
	copy(result, h.entries)
	return result
}

// Last returns the most recent entry, or nil if history is empty.
func (h *History) Last() *Entry {
	if len(h.entries) == 0 {
		return nil
	}
	e := h.entries[len(h.entries)-1]
	return &e
}

// Clear removes all entries from the history.
func (h *History) Clear() {
	h.entries = nil
}

// Len returns the number of entries currently stored.
func (h *History) Len() int {
	return len(h.entries)
}

// Format returns a human-readable summary of all history entries.
func (h *History) Format() string {
	if len(h.entries) == 0 {
		return "No history recorded."
	}
	var sb strings.Builder
	for i, e := range h.entries {
		status := "OK"
		if !e.Valid {
			status = "INVALID"
		} else if len(e.Issues) > 0 {
			status = fmt.Sprintf("WARN(%d)", len(e.Issues))
		}
		fmt.Fprintf(&sb, "%2d. [%s] %-20s %s\n",
			i+1,
			e.Timestamp.Format("15:04:05"),
			status,
			e.Expression,
		)
	}
	return sb.String()
}
