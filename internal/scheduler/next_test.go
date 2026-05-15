package scheduler_test

import (
	"testing"
	"time"

	"github.com/user/crontab-lint/internal/scheduler"
)

var ref = time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

func TestNext_EveryMinute(t *testing.T) {
	next, err := scheduler.Next("* * * * *", ref)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := ref.Add(time.Minute)
	if !next.Equal(want) {
		t.Errorf("got %v, want %v", next, want)
	}
}

func TestNext_HourlyAt30(t *testing.T) {
	next, err := scheduler.Next("30 * * * *", ref)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := time.Date(2024, 1, 15, 12, 30, 0, 0, time.UTC)
	if !next.Equal(want) {
		t.Errorf("got %v, want %v", next, want)
	}
}

func TestNext_DailyAtMidnight(t *testing.T) {
	next, err := scheduler.Next("0 0 * * *", ref)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := time.Date(2024, 1, 16, 0, 0, 0, 0, time.UTC)
	if !next.Equal(want) {
		t.Errorf("got %v, want %v", next, want)
	}
}

func TestNextN_ReturnsN(t *testing.T) {
	times, err := scheduler.NextN("0 * * * *", ref, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(times) != 3 {
		t.Fatalf("expected 3 times, got %d", len(times))
	}
	expected := []time.Time{
		time.Date(2024, 1, 15, 13, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 15, 14, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 15, 15, 0, 0, 0, time.UTC),
	}
	for i, got := range times {
		if !got.Equal(expected[i]) {
			t.Errorf("times[%d]: got %v, want %v", i, got, expected[i])
		}
	}
}

func TestNext_StepField(t *testing.T) {
	next, err := scheduler.Next("*/15 * * * *", ref)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := time.Date(2024, 1, 15, 12, 15, 0, 0, time.UTC)
	if !next.Equal(want) {
		t.Errorf("got %v, want %v", next, want)
	}
}

func TestNext_WrongFieldCount(t *testing.T) {
	_, err := scheduler.Next("* * * *", ref)
	if err == nil {
		t.Error("expected error for wrong field count")
	}
}

func TestNext_InvalidField(t *testing.T) {
	_, err := scheduler.Next("abc * * * *", ref)
	if err == nil {
		t.Error("expected error for invalid field")
	}
}

func TestNext_SpecificMonthAndDow(t *testing.T) {
	// First Monday of February 2024 is Feb 5
	next, err := scheduler.Next("0 9 * 2 1", ref)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if next.Month() != time.February {
		t.Errorf("expected February, got %v", next.Month())
	}
	if next.Weekday() != time.Monday {
		t.Errorf("expected Monday, got %v", next.Weekday())
	}
}
