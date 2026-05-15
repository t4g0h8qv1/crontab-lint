package matcher

import (
	"encoding/json"
	"fmt"
	"strings"
)

// FormatText returns a human-readable summary of a match Result.
func FormatText(r Result) string {
	var sb strings.Builder

	status := "NO MATCH"
	if r.Matches {
		status = "MATCH"
	}

	fmt.Fprintf(&sb, "Expression : %s\n", r.Expression)
	fmt.Fprintf(&sb, "Time       : %s\n", r.Time.Format("2006-01-02 15:04 MST"))
	fmt.Fprintf(&sb, "Result     : %s\n", status)
	sb.WriteString("\nField breakdown:\n")

	for _, f := range r.Fields {
		mark := "✓"
		if !f.Matches {
			mark = "✗"
		}
		fmt.Fprintf(&sb, "  %s %-14s value=%d\n", mark, f.Name, f.Value)
	}

	return sb.String()
}

// jsonResult is the JSON-serialisable form of Result.
type jsonResult struct {
	Expression string       `json:"expression"`
	Time       string       `json:"time"`
	Matches    bool         `json:"matches"`
	Fields     []jsonField  `json:"fields"`
}

type jsonField struct {
	Name    string `json:"name"`
	Value   int    `json:"value"`
	Matches bool   `json:"matches"`
}

// FormatJSON returns a JSON representation of a match Result.
func FormatJSON(r Result) (string, error) {
	fields := make([]jsonField, len(r.Fields))
	for i, f := range r.Fields {
		fields[i] = jsonField{Name: f.Name, Value: f.Value, Matches: f.Matches}
	}

	data := jsonResult{
		Expression: r.Expression,
		Time:       r.Time.Format("2006-01-02T15:04:05Z07:00"),
		Matches:    r.Matches,
		Fields:     fields,
	}

	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}
