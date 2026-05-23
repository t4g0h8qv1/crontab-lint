package truncator_test

import (
	"strings"
	"testing"

	"github.com/user/crontab-lint/internal/truncator"
)

func TestTruncate_NoChange(t *testing.T) {
	r := truncator.Truncate("30 4 * * 1")
	if r.Truncated != "30 4 * * 1" {
		t.Errorf("expected no change, got %q", r.Truncated)
	}
	if len(r.Changes) != 0 {
		t.Errorf("expected no changes, got %v", r.Changes)
	}
	if len(r.Errors) != 0 {
		t.Errorf("expected no errors, got %v", r.Errors)
	}
}

func TestTruncate_StepOne(t *testing.T) {
	r := truncator.Truncate("*/1 */1 * * *")
	if r.Truncated != "* * * * *" {
		t.Errorf("expected '* * * * *', got %q", r.Truncated)
	}
	if len(r.Changes) != 2 {
		t.Errorf("expected 2 changes, got %d", len(r.Changes))
	}
}

func TestTruncate_FullMinuteRange(t *testing.T) {
	r := truncator.Truncate("0-59 * * * *")
	if r.Truncated != "* * * * *" {
		t.Errorf("expected '* * * * *', got %q", r.Truncated)
	}
	if len(r.Changes) != 1 {
		t.Errorf("expected 1 change, got %d", len(r.Changes))
	}
}

func TestTruncate_FullRangeWithStepOne(t *testing.T) {
	r := truncator.Truncate("0-23/1 * * * *")
	if r.Truncated != "* * * * *" {
		t.Errorf("expected '* * * * *', got %q", r.Truncated)
	}
}

func TestTruncate_PartialRangeUnchanged(t *testing.T) {
	r := truncator.Truncate("0-30 * * * *")
	if r.Truncated != "0-30 * * * *" {
		t.Errorf("expected '0-30 * * * *', got %q", r.Truncated)
	}
	if len(r.Changes) != 0 {
		t.Errorf("expected no changes, got %v", r.Changes)
	}
}

func TestTruncate_InvalidExpression(t *testing.T) {
	r := truncator.Truncate("* * *")
	if len(r.Errors) == 0 {
		t.Error("expected error for wrong field count")
	}
	if r.Truncated != "" {
		t.Errorf("expected empty truncated on error, got %q", r.Truncated)
	}
}

func TestTruncate_ChangesDescribeFields(t *testing.T) {
	r := truncator.Truncate("*/1 * * * *")
	if len(r.Changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(r.Changes))
	}
	if !strings.Contains(r.Changes[0], "minute") {
		t.Errorf("expected change to mention 'minute', got %q", r.Changes[0])
	}
}

func TestTruncate_OriginalPreserved(t *testing.T) {
	expr := "0-59/1 0-23/1 * * *"
	r := truncator.Truncate(expr)
	if r.Original != expr {
		t.Errorf("expected original %q, got %q", expr, r.Original)
	}
}
