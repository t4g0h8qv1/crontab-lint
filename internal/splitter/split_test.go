package splitter_test

import (
	"testing"

	"github.com/user/crontab-lint/internal/splitter"
)

func TestSplit_Valid(t *testing.T) {
	f, err := splitter.Split("30 6 * * 1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Minute != "30" {
		t.Errorf("Minute: want 30, got %s", f.Minute)
	}
	if f.Hour != "6" {
		t.Errorf("Hour: want 6, got %s", f.Hour)
	}
	if f.DayOfMonth != "*" {
		t.Errorf("DayOfMonth: want *, got %s", f.DayOfMonth)
	}
	if f.Month != "*" {
		t.Errorf("Month: want *, got %s", f.Month)
	}
	if f.DayOfWeek != "1" {
		t.Errorf("DayOfWeek: want 1, got %s", f.DayOfWeek)
	}
}

func TestSplit_InvalidFieldCount(t *testing.T) {
	cases := []string{
		"",
		"* * *",
		"* * * * * *",
		"30 6 1",
	}
	for _, c := range cases {
		_, err := splitter.Split(c)
		if err == nil {
			t.Errorf("Split(%q): expected error, got nil", c)
		}
		if err != splitter.ErrInvalidFieldCount {
			t.Errorf("Split(%q): want ErrInvalidFieldCount, got %v", c, err)
		}
	}
}

func TestJoin_RoundTrip(t *testing.T) {
	expr := "0 12 * * 1-5"
	f, err := splitter.Split(expr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := splitter.Join(f)
	if got != expr {
		t.Errorf("Join: want %q, got %q", expr, got)
	}
}

func TestSlice_Order(t *testing.T) {
	f := splitter.Fields{
		Minute:     "5",
		Hour:       "4",
		DayOfMonth: "3",
		Month:      "2",
		DayOfWeek:  "1",
	}
	s := splitter.Slice(f)
	expected := []string{"5", "4", "3", "2", "1"}
	if len(s) != len(expected) {
		t.Fatalf("Slice length: want %d, got %d", len(expected), len(s))
	}
	for i, v := range expected {
		if s[i] != v {
			t.Errorf("Slice[%d]: want %q, got %q", i, v, s[i])
		}
	}
}

func TestFieldNames_Length(t *testing.T) {
	names := splitter.FieldNames()
	if len(names) != 5 {
		t.Errorf("FieldNames: want 5 names, got %d", len(names))
	}
}
