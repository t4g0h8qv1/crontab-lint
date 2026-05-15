package ranker

import (
	"encoding/json"
	"fmt"
	"strings"
)

// FormatText returns a human-readable ranking report.
func FormatText(s Score) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Expression : %s\n", s.Expression))
	sb.WriteString(fmt.Sprintf("Total Score: %d\n", s.Total))
	sb.WriteString(fmt.Sprintf("Complexity : %d\n", s.Complexity))
	sb.WriteString(fmt.Sprintf("Readability: %d\n", s.Readability))
	if len(s.Notes) > 0 {
		sb.WriteString("Notes:\n")
		for _, n := range s.Notes {
			sb.WriteString(fmt.Sprintf("  - %s\n", n))
		}
	}
	return sb.String()
}

// FormatJSON returns a JSON-encoded ranking report.
func FormatJSON(s Score) (string, error) {
	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return "", fmt.Errorf("ranker: marshal: %w", err)
	}
	return string(b), nil
}
