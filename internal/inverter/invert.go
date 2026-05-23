// Package inverter computes the logical complement of a cron expression,
// returning a new expression that fires at all times the original does NOT.
package inverter

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/user/crontab-lint/internal/expander"
)

// fieldRanges defines the inclusive [min, max] for each cron field.
var fieldRanges = [5][2]int{
	{0, 59}, // minute
	{0, 23}, // hour
	{1, 31}, // day-of-month
	{1, 12}, // month
	{0, 6},  // day-of-week
}

// Result holds the output of an Invert call.
type Result struct {
	Original  string
	Inverted  string
	ChangedFields []string
	Errors    []string
}

// Invert returns a cron expression representing all times NOT covered by expr.
// Fields that are already wildcards remain wildcards in the complement because
// their complement would be the empty set within a larger cross-product context.
func Invert(expr string) Result {
	res := Result{Original: expr}

	expanded, err := expander.Expand(expr)
	if err != nil {
		res.Errors = append(res.Errors, err.Error())
		return res
	}

	parts := strings.Fields(expr)
	if len(parts) != 5 {
		res.Errors = append(res.Errors, fmt.Sprintf("expected 5 fields, got %d", len(parts)))
		return res
	}

	invertedParts := make([]string, 5)
	for i, field := range parts {
		if field == "*" {
			invertedParts[i] = "*"
			continue
		}

		present := toSet(expanded.Fields[i])
		min, max := fieldRanges[i][0], fieldRanges[i][1]

		var complement []string
		for v := min; v <= max; v++ {
			if !present[v] {
				complement = append(complement, strconv.Itoa(v))
			}
		}

		if len(complement) == 0 {
			res.Errors = append(res.Errors, fmt.Sprintf("field %d has no complement (covers all values)", i))
			return res
		}

		if len(complement) == max-min+1 {
			invertedParts[i] = "*"
		} else {
			invertedParts[i] = strings.Join(complement, ",")
			res.ChangedFields = append(res.ChangedFields, expander.FieldName(i))
		}
	}

	res.Inverted = strings.Join(invertedParts, " ")
	return res
}

func toSet(values []int) map[int]bool {
	s := make(map[int]bool, len(values))
	for _, v := range values {
		s[v] = true
	}
	return s
}
