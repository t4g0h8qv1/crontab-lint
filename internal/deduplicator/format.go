package deduplicator

import (
	"encoding/json"
	"fmt"
	"strings"
)

// FormatText returns a human-readable summary of a deduplication Result.
func FormatText(r Result) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Original : %s\n", r.Original))
	sb.WriteString(fmt.Sprintf("Deduped  : %s\n", r.Deduped))
	if !r.Changed {
		sb.WriteString("Status   : no duplicates found\n")
		return sb.String()
	}
	sb.WriteString(fmt.Sprintf("Status   : %d change(s) made\n", len(r.Changes)))
	for _, c := range r.Changes {
		sb.WriteString(fmt.Sprintf("  - %s\n", c))
	}
	return sb.String()
}

type jsonResult struct {
	Original string   `json:"original"`
	Deduped  string   `json:"deduped"`
	Changed  bool     `json:"changed"`
	Changes  []string `json:"changes"`
}

// FormatJSON returns a JSON-encoded representation of a deduplication Result.
func FormatJSON(r Result) (string, error) {
	j := jsonResult{
		Original: r.Original,
		Deduped:  r.Deduped,
		Changed:  r.Changed,
		Changes:  r.Changes,
	}
	if j.Changes == nil {
		j.Changes = []string{}
	}
	b, err := json.MarshalIndent(j, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}
