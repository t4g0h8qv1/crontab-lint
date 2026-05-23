package canonicalizer

import (
	"encoding/json"
	"fmt"
	"strings"
)

// FormatText returns a human-readable representation of a canonicalization Result.
func FormatText(r Result) string {
	if r.Error != "" {
		return fmt.Sprintf("Error: %s\n", r.Error)
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Original:  %s\n", r.Original))
	sb.WriteString(fmt.Sprintf("Canonical: %s\n", r.Canonical))

	if len(r.Changes) == 0 {
		sb.WriteString("No changes.\n")
	} else {
		sb.WriteString(fmt.Sprintf("Changes (%d):\n", len(r.Changes)))
		for _, c := range r.Changes {
			sb.WriteString(fmt.Sprintf("  - %s\n", c))
		}
	}

	return sb.String()
}

// FormatJSON returns a JSON representation of a canonicalization Result.
func FormatJSON(r Result) (string, error) {
	b, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return "", fmt.Errorf("canonicalizer: marshal: %w", err)
	}
	return string(b), nil
}
