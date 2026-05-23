// Package annotator attaches inline comments to each field of a crontab
// expression, producing an annotated representation suitable for display.
package annotator

import (
	"fmt"
	"strings"

	"github.com/user/crontab-lint/internal/explainer"
	"github.com/user/crontab-lint/internal/splitter"
)

// Field holds a single crontab field together with its annotation.
type Field struct {
	Name  string `json:"name"`
	Value string `json:"value"`
	Note  string `json:"note"`
}

// Result is the output of Annotate.
type Result struct {
	Expression string  `json:"expression"`
	Fields     []Field `json:"fields"`
	Errors     []string `json:"errors,omitempty"`
}

// Annotate parses expr and returns each field paired with a human-readable
// note derived from the explainer package.
func Annotate(expr string) Result {
	result := Result{Expression: expr}

	split, err := splitter.Split(expr)
	if err != nil {
		result.Errors = []string{err.Error()}
		return result
	}

	names := splitter.FieldNames()
	for i, val := range split {
		note := explainer.ExplainField(names[i], val)
		result.Fields = append(result.Fields, Field{
			Name:  names[i],
			Value: val,
			Note:  note,
		})
	}

	return result
}

// Inline returns the expression with a comment appended that summarises each
// field, e.g. "*/5 * * * *  # every 5 min | any hour | ..."
func Inline(expr string) string {
	r := Annotate(expr)
	if len(r.Errors) > 0 {
		return fmt.Sprintf("%s  # error: %s", expr, strings.Join(r.Errors, "; "))
	}

	parts := make([]string, len(r.Fields))
	for i, f := range r.Fields {
		parts[i] = f.Note
	}
	return fmt.Sprintf("%s  # %s", expr, strings.Join(parts, " | "))
}
