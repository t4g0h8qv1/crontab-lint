package trimmer_test

import (
	"strings"
	"testing"

	"github.com/user/crontab-lint/internal/trimmer"
)

func TestTrim_StepOne(t *testing.T) {
	res, err := trimmer.Trim("*/1 * * * *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Trimmed != "* * * * *" {
		t.Errorf("expected '* * * * *', got %q", res.Trimmed)
	}
	if len(res.Changes) == 0 {
		t.Error("expected at least one change recorded")
	}
}

func TestTrim_NoChange(t *testing.T) {
	res, err := trimmer.Trim("0 9 * * 1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Trimmed != "0 9 * * 1" {
		t.Errorf("expected expression unchanged, got %q", res.Trimmed)
	}
	if len(res.Changes) != 0 {
		t.Errorf("expected no changes, got %v", res.Changes)
	}
}

func TestTrim_DuplicateValues(t *testing.T) {
	res, err := trimmer.Trim("0,0,30 * * * *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Trimmed != "0,30 * * * *" {
		t.Errorf("expected '0,30 * * * *', got %q", res.Trimmed)
	}
	if len(res.Changes) == 0 {
		t.Error("expected a change to be recorded for duplicate removal")
	}
}

func TestTrim_MultipleStepOnes(t *testing.T) {
	res, err := trimmer.Trim("*/1 */1 * * *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Trimmed != "* * * * *" {
		t.Errorf("expected '* * * * *', got %q", res.Trimmed)
	}
	if len(res.Changes) < 2 {
		t.Errorf("expected 2 changes, got %d", len(res.Changes))
	}
}

func TestTrim_InvalidExpression(t *testing.T) {
	_, err := trimmer.Trim("* * *")
	if err == nil {
		t.Error("expected error for invalid field count")
	}
	if !strings.Contains(err.Error(), "trim:") {
		t.Errorf("expected error to contain 'trim:', got %q", err.Error())
	}
}

func TestTrim_OriginalPreserved(t *testing.T) {
	expr := "*/1 0 * * *"
	res, err := trimmer.Trim(expr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Original != expr {
		t.Errorf("expected original to be %q, got %q", expr, res.Original)
	}
}
