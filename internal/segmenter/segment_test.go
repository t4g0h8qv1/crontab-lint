package segmenter

import (
	"testing"
)

func TestSegment_Valid(t *testing.T) {
	segs, err := Segment("*/5 0 * * 1-5")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(segs) != 5 {
		t.Fatalf("expected 5 segments, got %d", len(segs))
	}

	expected := []struct {
		name  string
		value string
	}{
		{"minute", "*/5"},
		{"hour", "0"},
		{"day-of-month", "*"},
		{"month", "*"},
		{"day-of-week", "1-5"},
	}
	for i, e := range expected {
		if segs[i].Name != e.name {
			t.Errorf("seg[%d].Name = %q, want %q", i, segs[i].Name, e.name)
		}
		if segs[i].Value != e.value {
			t.Errorf("seg[%d].Value = %q, want %q", i, segs[i].Value, e.value)
		}
		if segs[i].Index != i {
			t.Errorf("seg[%d].Index = %d, want %d", i, segs[i].Index, i)
		}
	}
}

func TestSegment_Offsets(t *testing.T) {
	expr := "30 4 * * *"
	segs, err := Segment(expr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// "30" starts at 0, ends at 2
	if segs[0].Start != 0 || segs[0].End != 2 {
		t.Errorf("minute offsets: got [%d,%d), want [0,2)", segs[0].Start, segs[0].End)
	}
	// "4" starts at 3, ends at 4
	if segs[1].Start != 3 || segs[1].End != 4 {
		t.Errorf("hour offsets: got [%d,%d), want [3,4)", segs[1].Start, segs[1].End)
	}
}

func TestSegment_WrongFieldCount(t *testing.T) {
	_, err := Segment("* * *")
	if err == nil {
		t.Fatal("expected error for wrong field count")
	}
}

func TestSegment_EmptyExpression(t *testing.T) {
	_, err := Segment("")
	if err == nil {
		t.Fatal("expected error for empty expression")
	}
}

func TestAtOffset_HitsField(t *testing.T) {
	segs, _ := Segment("30 4 * * *")
	s := AtOffset(segs, 0)
	if s == nil || s.Name != "minute" {
		t.Errorf("expected minute field at offset 0, got %v", s)
	}
}

func TestAtOffset_OnWhitespace(t *testing.T) {
	segs, _ := Segment("30 4 * * *")
	s := AtOffset(segs, 2) // space between "30" and "4"
	if s != nil {
		t.Errorf("expected nil for whitespace offset, got %v", s)
	}
}

func TestAtOffset_LastField(t *testing.T) {
	expr := "0 0 * * 7"
	segs, _ := Segment(expr)
	s := AtOffset(segs, len(expr)-1)
	if s == nil || s.Name != "day-of-week" {
		t.Errorf("expected day-of-week at last offset, got %v", s)
	}
}
