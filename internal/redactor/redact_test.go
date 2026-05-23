package redactor

import (
	"testing"
)

func TestRedact_NoMask(t *testing.T) {
	r := Redact("30 6 * * 1", "*", nil)
	if len(r.Errors) != 0 {
		t.Fatalf("unexpected errors: %v", r.Errors)
	}
	if r.Redacted != r.Original {
		t.Errorf("expected no change, got %q", r.Redacted)
	}
	if len(r.Changed) != 0 {
		t.Errorf("expected no changed fields, got %v", r.Changed)
	}
}

func TestRedact_SingleField(t *testing.T) {
	r := Redact("30 6 * * 1", "*", []int{0})
	if len(r.Errors) != 0 {
		t.Fatalf("unexpected errors: %v", r.Errors)
	}
	if r.Redacted != "* 6 * * 1" {
		t.Errorf("got %q", r.Redacted)
	}
	if len(r.Changed) != 1 || r.Changed[0] != 0 {
		t.Errorf("expected changed=[0], got %v", r.Changed)
	}
}

func TestRedact_MultipleFields(t *testing.T) {
	r := Redact("30 6 15 3 1", "*", []int{1, 3})
	if len(r.Errors) != 0 {
		t.Fatalf("unexpected errors: %v", r.Errors)
	}
	if r.Redacted != "30 * 15 * 1" {
		t.Errorf("got %q", r.Redacted)
	}
}

func TestRedact_AlreadyWildcard_NotReportedChanged(t *testing.T) {
	r := Redact("* 6 * * 1", "*", []int{0})
	if len(r.Changed) != 0 {
		t.Errorf("field already wildcard should not appear in Changed, got %v", r.Changed)
	}
}

func TestRedact_CustomReplacement(t *testing.T) {
	r := Redact("30 6 * * 1", "X", []int{0, 4})
	if r.Redacted != "X 6 * * X" {
		t.Errorf("got %q", r.Redacted)
	}
}

func TestRedact_InvalidFieldCount(t *testing.T) {
	r := Redact("30 6 *", "*", []int{0})
	if len(r.Errors) == 0 {
		t.Fatal("expected error for wrong field count")
	}
}

func TestRedact_OutOfRangeIndex(t *testing.T) {
	r := Redact("30 6 * * 1", "*", []int{5})
	if len(r.Errors) == 0 {
		t.Fatal("expected error for out-of-range index")
	}
}

func TestRedact_DuplicateIndex(t *testing.T) {
	r := Redact("30 6 * * 1", "*", []int{1, 1})
	if len(r.Errors) == 0 {
		t.Fatal("expected error for duplicate index")
	}
}

func TestRedact_FieldsSliceLength(t *testing.T) {
	r := Redact("0 0 1 1 0", "*", []int{2})
	if len(r.Fields) != 5 {
		t.Errorf("expected 5 fields, got %d", len(r.Fields))
	}
	if r.Fields[2] != "*" {
		t.Errorf("expected Fields[2]=*, got %q", r.Fields[2])
	}
}
