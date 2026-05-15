package summarizer

import (
	"encoding/json"
	"fmt"
	"strings"
)

// FormatText returns a plain-text representation of a summarizer Result.
func FormatText(r Result) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "Expression : %s\n", r.Expression)
	if r.Normalized != "" && r.Normalized != r.Expression {
		fmt.Fprintf(&sb, "Normalized : %s\n", r.Normalized)
	}
	if r.Class != "" {
		fmt.Fprintf(&sb, "Class      : %s\n", r.Class)
	}
	fmt.Fprintf(&sb, "Summary    : %s\n", r.Summary)
	if len(r.Errors) > 0 {
		fmt.Fprintf(&sb, "Errors     :\n")
		for _, e := range r.Errors {
			fmt.Fprintf(&sb, "  - %s\n", e)
		}
	}
	return sb.String()
}

// FormatJSON returns a JSON representation of a summarizer Result.
func FormatJSON(r Result) (string, error) {
	b, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return "", fmt.Errorf("summarizer: marshal error: %w", err)
	}
	return string(b), nil
}
