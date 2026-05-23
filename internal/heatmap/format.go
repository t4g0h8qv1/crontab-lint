package heatmap

import (
	"encoding/json"
	"fmt"
	"strings"
)

var (
	dayNames  = []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}
	density   = []rune{' ', '░', '▒', '▓', '█'}
)

// FormatText renders the heatmap as an ASCII grid (hours across, days down).
func FormatText(r Result) string {
	if len(r.Errors) > 0 {
		return fmt.Sprintf("heatmap error: %s", strings.Join(r.Errors, "; "))
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Heatmap for: %s\n", r.Expression))

	// Header row: hours 0-23
	sb.WriteString("     ")
	for h := 0; h < 24; h++ {
		sb.WriteString(fmt.Sprintf("%02d ", h))
	}
	sb.WriteByte('\n')

	// One row per day
	for d := 0; d < 7; d++ {
		sb.WriteString(fmt.Sprintf("%s  ", dayNames[d]))
		for h := 0; h < 24; h++ {
			cell := r.Cells[d*24+h]
			sb.WriteRune(rune(block(cell.Count, r.MaxCount)))
			sb.WriteString("  ")
		}
		sb.WriteByte('\n')
	}

	sb.WriteString(fmt.Sprintf("Legend: %s (low→high)\n",
		string(density)))
	return sb.String()
}

// FormatJSON returns a JSON representation of the heatmap result.
func FormatJSON(r Result) (string, error) {
	type jsonCell struct {
		Day   string `json:"day"`
		Hour  int    `json:"hour"`
		Count int    `json:"count"`
	}
	type jsonResult struct {
		Expression string      `json:"expression"`
		MaxCount   int         `json:"max_count"`
		Cells      []jsonCell  `json:"cells"`
		Errors     []string    `json:"errors,omitempty"`
	}

	jCells := make([]jsonCell, len(r.Cells))
	for i, c := range r.Cells {
		jCells[i] = jsonCell{Day: dayNames[c.DOW], Hour: c.Hour, Count: c.Count}
	}

	out := jsonResult{
		Expression: r.Expression,
		MaxCount:   r.MaxCount,
		Cells:      jCells,
		Errors:     r.Errors,
	}
	b, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// block maps a count to a density rune index.
func block(count, max int) rune {
	if max == 0 || count == 0 {
		return density[0]
	}
	idx := 1 + (count-1)*(len(density)-1)/max
	if idx >= len(density) {
		idx = len(density) - 1
	}
	return density[idx]
}
