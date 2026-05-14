package validator

import (
	"fmt"
	"strings"

	"github.com/user/crontab-lint/internal/parser"
)

// FieldSpec defines the valid range for a crontab field.
type FieldSpec struct {
	Name string
	Min  int
	Max  int
}

// Standard crontab field specifications.
var standardFields = []FieldSpec{
	{Name: "minute", Min: 0, Max: 59},
	{Name: "hour", Min: 0, Max: 23},
	{Name: "day-of-month", Min: 1, Max: 31},
	{Name: "month", Min: 1, Max: 12},
	{Name: "day-of-week", Min: 0, Max: 7},
}

// ValidationError holds a field-level validation error.
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("field %q: %s", e.Field, e.Message)
}

// Result holds the outcome of validating a crontab expression.
type Result struct {
	Expression string
	Valid      bool
	Errors     []ValidationError
}

// Validate parses and validates a crontab expression string.
func Validate(expr string) Result {
	result := Result{Expression: expr}

	parts := strings.Fields(expr)
	if len(parts) != 5 {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "expression",
			Message: fmt.Sprintf("expected 5 fields, got %d", len(parts)),
		})
		result.Valid = false
		return result
	}

	for i, spec := range standardFields {
		field, err := parser.ParseField(parts[i], spec.Min, spec.Max)
		if err != nil {
			result.Errors = append(result.Errors, ValidationError{
				Field:   spec.Name,
				Message: err.Error(),
			})
			continue
		}
		_ = field
	}

	result.Valid = len(result.Errors) == 0
	return result
}
