package validator_test

import (
	"testing"

	"github.com/user/crontab-lint/internal/validator"
)

func TestValidate_Valid(t *testing.T) {
	cases := []struct {
		expr string
	}{
		{"* * * * *"},
		{"0 * * * *"},
		{"0 12 * * 1"},
		{"*/15 * * * *"},
		{"0 0 1 1 *"},
		{"30 6 */2 * 1-5"},
	}

	for _, tc := range cases {
		t.Run(tc.expr, func(t *testing.T) {
			res := validator.Validate(tc.expr)
			if !res.Valid {
				t.Errorf("expected valid, got errors: %v", res.Errors)
			}
		})
	}
}

func TestValidate_Invalid(t *testing.T) {
	cases := []struct {
		expr    string
		errCount int
	}{
		{"* * * *", 1},       // too few fields
		{"60 * * * *", 1},   // minute out of range
		{"* 25 * * *", 1},   // hour out of range
		{"* * 0 * *", 1},    // day-of-month out of range
		{"* * * 13 *", 1},   // month out of range
		{"* * * * 8", 1},    // day-of-week out of range
	}

	for _, tc := range cases {
		t.Run(tc.expr, func(t *testing.T) {
			res := validator.Validate(tc.expr)
			if res.Valid {
				t.Errorf("expected invalid expression %q to fail", tc.expr)
			}
			if len(res.Errors) != tc.errCount {
				t.Errorf("expected %d error(s), got %d: %v", tc.errCount, len(res.Errors), res.Errors)
			}
		})
	}
}
