package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/crontab-lint/internal/differ"
)

func TestDiffIntegration_EquivalentExpressions(t *testing.T) {
	r, err := differ.Diff("@hourly", "0 * * * *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !r.Equal {
		t.Errorf("@hourly and '0 * * * *' should be equivalent")
	}
	summary := differ.Summary(r)
	if !strings.Contains(summary, "equivalent") {
		t.Errorf("summary should mention equivalent, got: %q", summary)
	}
}

func TestDiffIntegration_DifferentExpressions(t *testing.T) {
	r, err := differ.Diff("0 9 * * 1", "0 17 * * 5")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Equal {
		t.Error("expected expressions to differ")
	}

	fields := map[string]bool{}
	for _, c := range r.Changes {
		fields[c.Field] = true
	}
	if !fields["hour"] {
		t.Error("expected 'hour' field to differ")
	}
	if !fields["day-of-week"] {
		t.Error("expected 'day-of-week' field to differ")
	}
}

func TestDiffSummary_OutputFormat(t *testing.T) {
	r, err := differ.Diff("0 0 * * *", "0 12 * * *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var buf bytes.Buffer
	buf.WriteString(differ.Summary(r))
	out := buf.String()

	if !strings.Contains(out, "→") {
		t.Errorf("expected arrow separator in diff output, got: %q", out)
	}
	if !strings.Contains(out, "hour") {
		t.Errorf("expected field name in diff output, got: %q", out)
	}
}
