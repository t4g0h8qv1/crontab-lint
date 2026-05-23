package rotator

import (
	"testing"
)

func TestRotate_NoOffset(t *testing.T) {
	r := Rotate("30 6 * * 1", [5]int{0, 0, 0, 0, 0})
	if r.Rotated != "30 6 * * 1" {
		t.Errorf("expected no change, got %q", r.Rotated)
	}
	if len(r.Changed) != 0 {
		t.Errorf("expected no changed fields, got %v", r.Changed)
	}
}

func TestRotate_ShiftMinute(t *testing.T) {
	r := Rotate("30 6 * * *", [5]int{15, 0, 0, 0, 0})
	if r.Rotated != "45 6 * * *" {
		t.Errorf("expected '45 6 * * *', got %q", r.Rotated)
	}
	if len(r.Changed) != 1 || r.Changed[0] != "minute" {
		t.Errorf("expected changed=[minute], got %v", r.Changed)
	}
}

func TestRotate_WrapAround(t *testing.T) {
	// minute 50 + 20 = 70 => wraps to 10
	r := Rotate("50 6 * * *", [5]int{20, 0, 0, 0, 0})
	if r.Rotated != "10 6 * * *" {
		t.Errorf("expected '10 6 * * *', got %q", r.Rotated)
	}
}

func TestRotate_NegativeOffset(t *testing.T) {
	// minute 5 - 10 => wraps to 55
	r := Rotate("5 12 * * *", [5]int{-10, 0, 0, 0, 0})
	if r.Rotated != "55 12 * * *" {
		t.Errorf("expected '55 12 * * *', got %q", r.Rotated)
	}
}

func TestRotate_MultipleFields(t *testing.T) {
	r := Rotate("0 6 * * 1", [5]int{30, 2, 0, 0, 3})
	if r.Rotated != "30 8 * * 4" {
		t.Errorf("expected '30 8 * * 4', got %q", r.Rotated)
	}
	if len(r.Changed) != 3 {
		t.Errorf("expected 3 changed fields, got %d: %v", len(r.Changed), r.Changed)
	}
}

func TestRotate_CommaList(t *testing.T) {
	r := Rotate("0,30 * * * *", [5]int{15, 0, 0, 0, 0})
	if r.Rotated != "15,45 * * * *" {
		t.Errorf("expected '15,45 * * * *', got %q", r.Rotated)
	}
}

func TestRotate_WildcardUnchanged(t *testing.T) {
	r := Rotate("* * * * *", [5]int{5, 5, 5, 5, 5})
	if r.Rotated != "* * * * *" {
		t.Errorf("expected wildcards unchanged, got %q", r.Rotated)
	}
	if len(r.Changed) != 0 {
		t.Errorf("expected no changed fields for wildcards, got %v", r.Changed)
	}
}

func TestRotate_InvalidFieldCount(t *testing.T) {
	r := Rotate("* * *", [5]int{})
	if len(r.Errors) == 0 {
		t.Error("expected error for wrong field count")
	}
	if r.Rotated != "* * *" {
		t.Errorf("expected original returned on error, got %q", r.Rotated)
	}
}

func TestRotate_NonNumericField(t *testing.T) {
	// Step expressions are not supported for rotation
	r := Rotate("*/5 * * * *", [5]int{3, 0, 0, 0, 0})
	if len(r.Errors) == 0 {
		t.Error("expected error for non-numeric step field")
	}
}
