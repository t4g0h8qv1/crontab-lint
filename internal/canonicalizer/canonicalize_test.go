package canonicalizer_test

import (
	"strings"
	"testing"

	"github.com/user/crontab-lint/internal/canonicalizer"
)

func TestCanonicalize_NoChange(t *testing.T) {
	r := canonicalizer.Canonicalize("0 * * * *")
	if r.Error != "" {
		t.Fatalf("unexpected error: %s", r.Error)
	}
	if r.Canonical != "0 * * * *" {
		t.Errorf("expected no change, got %q", r.Canonical)
	}
	if len(r.Changes) != 0 {
		t.Errorf("expected no changes, got %v", r.Changes)
	}
}

func TestCanonicalize_StepOne(t *testing.T) {
	r := canonicalizer.Canonicalize("*/1 * * * *")
	if r.Error != "" {
		t.Fatalf("unexpected error: %s", r.Error)
	}
	if r.Canonical != "* * * * *" {
		t.Errorf("expected '* * * * *', got %q", r.Canonical)
	}
	if len(r.Changes) == 0 {
		t.Error("expected at least one change")
	}
}

func TestCanonicalize_SortsList(t *testing.T) {
	r := canonicalizer.Canonicalize("30,5,15 * * * *")
	if r.Error != "" {
		t.Fatalf("unexpected error: %s", r.Error)
	}
	if r.Canonical != "5,15,30 * * * *" {
		t.Errorf("expected sorted list, got %q", r.Canonical)
	}
	if len(r.Changes) == 0 {
		t.Error("expected a change for sorted list")
	}
}

func TestCanonicalize_NormalizesAlias(t *testing.T) {
	r := canonicalizer.Canonicalize("@hourly")
	if r.Error != "" {
		t.Fatalf("unexpected error: %s", r.Error)
	}
	if r.Canonical == "@hourly" {
		t.Error("expected alias to be expanded")
	}
}

func TestCanonicalize_InvalidExpression(t *testing.T) {
	r := canonicalizer.Canonicalize("bad expression here")
	if r.Error == "" {
		t.Error("expected error for invalid expression")
	}
}

func TestCanonicalize_AlreadySorted(t *testing.T) {
	r := canonicalizer.Canonicalize("5,15,30 * * * *")
	if r.Error != "" {
		t.Fatalf("unexpected error: %s", r.Error)
	}
	if len(r.Changes) != 0 {
		t.Errorf("expected no changes for already-sorted list, got %v", r.Changes)
	}
}

func TestCanonicalize_MultipleChanges(t *testing.T) {
	r := canonicalizer.Canonicalize("*/1 */1 * * *")
	if r.Error != "" {
		t.Fatalf("unexpected error: %s", r.Error)
	}
	if !strings.Contains(r.Canonical, "* *") {
		t.Errorf("expected both fields simplified, got %q", r.Canonical)
	}
	if len(r.Changes) < 2 {
		t.Errorf("expected at least 2 changes, got %d", len(r.Changes))
	}
}
