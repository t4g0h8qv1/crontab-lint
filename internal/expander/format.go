package expander

import (
	"encoding/json"
	"fmt"
	"strings"
)

// FormatText returns a human-readable multi-line representation of
// an expansion Result.
func FormatText(r Result) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "Expression: %s\n", r.Expression)

	if len(r.Errors) > 0 {
		sb.WriteString("Errors:\n")
		for _, e := range r.Errors {
			fmt.Fprintf(&sb, "  - %s\n", e)
		}
		return sb.String()
	}

	sb.WriteString("Expanded fields:\n")
	for _, f := range r.Fields {
		fmt.Fprintf(&sb, "  %-8s %s\n", f.Name+":", JoinValues(f.Values))
	}
	return sb.String()
}

// jsonResult mirrors Result for JSON serialisation with string values.
type jsonResult struct {
	Expression string            `json:"expression"`
	Fields     []jsonField       `json:"fields,omitempty"`
	Errors     []string          `json:"errors,omitempty"`
}

type jsonField struct {
	Name   string `json:"name"`
	Values string `json:"values"`
	Count  int    `json:"count"`
}

// FormatJSON returns a JSON representation of the expansion Result.
func FormatJSON(r Result) (string, error) {
	jr := jsonResult{
		Expression: r.Expression,
		Errors:     r.Errors,
	}
	for _, f := range r.Fields {
		jr.Fields = append(jr.Fields, jsonField{
			Name:   f.Name,
			Values: JoinValues(f.Values),
			Count:  len(f.Values),
		})
	}
	b, err := json.MarshalIndent(jr, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}
