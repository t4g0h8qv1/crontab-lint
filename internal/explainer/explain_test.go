package explainer_test

import (
	"strings"
	"testing"

	"github.com/user/crontab-lint/internal/explainer"
)

func TestExplainField(t *testing.T) {
	cases := []struct {
		token     string
		field     string
		contains  string
	}{
		{"*", "minute", "every minute"},
		{"*/15", "minute", "every 15 minute"},
		{"1-5", "day-of-week", "day-of-week from 1 to 5"},
		{"0", "hour", "hour 0"},
		{"1,3,5", "month", "month 1,3,5"},
	}

	for _, tc := range cases {
		t.Run(tc.token, func(t *testing.T) {
			got := explainer.ExplainField(tc.token, tc.field, 0, 59)
			if !strings.Contains(got, tc.contains) {
				t.Errorf("ExplainField(%q, %q) = %q; want it to contain %q", tc.token, tc.field, got, tc.contains)
			}
		})
	}
}

func TestExplain_Valid(t *testing.T) {
	ex, err := explainer.Explain("0 12 * * 1-5")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(ex.Minute, "minute 0") {
		t.Errorf("unexpected minute explanation: %s", ex.Minute)
	}
	if !strings.Contains(ex.Hour, "hour 12") {
		t.Errorf("unexpected hour explanation: %s", ex.Hour)
	}
	if ex.Summary == "" {
		t.Error("expected non-empty summary")
	}
}

func TestExplain_InvalidFieldCount(t *testing.T) {
	_, err := explainer.Explain("* * *")
	if err == nil {
		t.Error("expected error for wrong field count")
	}
}
