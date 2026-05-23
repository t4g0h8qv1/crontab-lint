package inverter_test

import (
	"strings"
	"testing"

	"github.com/user/crontab-lint/internal/inverter"
)

func TestInvert_SingleMinute(t *testing.T) {
	res := inverter.Invert("0 * * * *")
	if len(res.Errors) != 0 {
		t.Fatalf("unexpected errors: %v", res.Errors)
	}
	// minute 0 inverted => 1-59
	parts := strings.Fields(res.Inverted)
	if parts[0] == "*" || parts[0] == "0" {
		t.Errorf("expected inverted minute field, got %q", parts[0])
	}
	// other fields should remain wildcard
	for _, p := range parts[1:] {
		if p != "*" {
			t.Errorf("expected wildcard for non-minute field, got %q", p)
		}
	}
}

func TestInvert_AllWildcards(t *testing.T) {
	res := inverter.Invert("* * * * *")
	if len(res.Errors) != 0 {
		t.Fatalf("unexpected errors: %v", res.Errors)
	}
	if res.Inverted != "* * * * *" {
		t.Errorf("expected all wildcards, got %q", res.Inverted)
	}
	if len(res.ChangedFields) != 0 {
		t.Errorf("expected no changed fields, got %v", res.ChangedFields)
	}
}

func TestInvert_InvalidFieldCount(t *testing.T) {
	res := inverter.Invert("0 * *")
	if len(res.Errors) == 0 {
		t.Fatal("expected error for wrong field count")
	}
}

func TestInvert_InvalidExpression(t *testing.T) {
	res := inverter.Invert("99 * * * *")
	if len(res.Errors) == 0 {
		t.Fatal("expected error for out-of-range minute")
	}
}

func TestInvert_ChangedFieldsReported(t *testing.T) {
	res := inverter.Invert("30 12 * * *")
	if len(res.Errors) != 0 {
		t.Fatalf("unexpected errors: %v", res.Errors)
	}
	if len(res.ChangedFields) < 2 {
		t.Errorf("expected at least 2 changed fields, got %v", res.ChangedFields)
	}
}

func TestInvert_OriginalPreserved(t *testing.T) {
	expr := "0 9 * * 1"
	res := inverter.Invert(expr)
	if res.Original != expr {
		t.Errorf("expected original %q, got %q", expr, res.Original)
	}
}

func TestInvert_StepExpression(t *testing.T) {
	res := inverter.Invert("*/15 * * * *")
	if len(res.Errors) != 0 {
		t.Fatalf("unexpected errors: %v", res.Errors)
	}
	// */15 fires at 0,15,30,45 — complement is 56 values
	parts := strings.Fields(res.Inverted)
	vals := strings.Split(parts[0], ",")
	if len(vals) != 56 {
		t.Errorf("expected 56 complement values for */15, got %d", len(vals))
	}
}
