package differ_test

import (
	"strings"
	"testing"

	"github.com/user/crontab-lint/internal/differ"
)

func TestDiff_EqualExpressions(t *testing.T) {
	r, err := differ.Diff("0 * * * *", "0 * * * *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !r.Equal {
		t.Errorf("expected expressions to be equal")
	}
	if len(r.Changes) != 0 {
		t.Errorf("expected no changes, got %d", len(r.Changes))
	}
}

func TestDiff_SingleFieldChange(t *testing.T) {
	r, err := differ.Diff("0 * * * *", "30 * * * *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Equal {
		t.Error("expected expressions to differ")
	}
	if len(r.Changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(r.Changes))
	}
	c := r.Changes[0]
	if c.Field != "minute" {
		t.Errorf("expected field 'minute', got %q", c.Field)
	}
	if c.From != "0" || c.To != "30" {
		t.Errorf("unexpected values: from=%q to=%q", c.From, c.To)
	}
}

func TestDiff_NormalizesAliases(t *testing.T) {
	// @daily normalizes to "0 0 * * *"
	r, err := differ.Diff("@daily", "0 0 * * *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !r.Equal {
		t.Errorf("expected @daily and '0 0 * * *' to be equivalent, changes: %+v", r.Changes)
	}
}

func TestDiff_MultipleFieldChanges(t *testing.T) {
	r, err := differ.Diff("0 0 * * *", "30 6 1 * 0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Changes) != 3 {
		t.Errorf("expected 3 changes, got %d", len(r.Changes))
	}
}

func TestDiff_InvalidFrom(t *testing.T) {
	_, err := differ.Diff("bad", "0 * * * *")
	if err == nil {
		t.Error("expected error for invalid 'from' expression")
	}
}

func TestDiff_InvalidTo(t *testing.T) {
	_, err := differ.Diff("0 * * * *", "also bad")
	if err == nil {
		t.Error("expected error for invalid 'to' expression")
	}
}

func TestSummary_Equal(t *testing.T) {
	r, _ := differ.Diff("* * * * *", "* * * * *")
	s := differ.Summary(r)
	if !strings.Contains(s, "equivalent") {
		t.Errorf("expected 'equivalent' in summary, got: %q", s)
	}
}

func TestSummary_WithChanges(t *testing.T) {
	r, _ := differ.Diff("0 * * * *", "30 12 * * *")
	s := differ.Summary(r)
	if !strings.Contains(s, "minute") {
		t.Errorf("expected 'minute' in summary, got: %q", s)
	}
	if !strings.Contains(s, "hour") {
		t.Errorf("expected 'hour' in summary, got: %q", s)
	}
}
