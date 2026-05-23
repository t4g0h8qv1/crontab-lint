// Package rewriter transforms crontab expressions by applying
// field-level substitutions, producing a rewritten expression
// along with a record of what changed.
package rewriter

import (
	"fmt"
	"strings"

	"github.com/user/crontab-lint/internal/splitter"
)

// Change records a single field substitution.
type Change struct {
	Field    string
	Original string
	Rewritten string
}

// Result holds the output of a Rewrite operation.
type Result struct {
	Original  string
	Rewritten string
	Changes   []Change
	Errors    []string
}

// Replacements maps field names to their replacement values.
// Valid keys: "minute", "hour", "dom", "month", "dow".
type Replacements map[string]string

// Rewrite applies the given replacements to the crontab expression.
// Fields not present in replacements are left unchanged.
// Returns a Result describing the transformation.
func Rewrite(expr string, replacements Replacements) Result {
	result := Result{Original: expr}

	fields, err := splitter.Split(expr)
	if err != nil {
		result.Errors = append(result.Errors, err.Error())
		result.Rewritten = expr
		return result
	}

	names := splitter.FieldNames()
	output := make([]string, len(fields))

	for i, name := range names {
		original := fields[i]
		if replacement, ok := replacements[name]; ok {
			if replacement == "" {
				result.Errors = append(result.Errors,
					fmt.Sprintf("replacement for field %q must not be empty", name))
				output[i] = original
				continue
			}
			if replacement != original {
				result.Changes = append(result.Changes, Change{
					Field:     name,
					Original:  original,
					Rewritten: replacement,
				})
			}
			output[i] = replacement
		} else {
			output[i] = original
		}
	}

	result.Rewritten = strings.Join(output, " ")
	return result
}
