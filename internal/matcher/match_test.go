package matcher_test

import (
	"testing"
	"time"

	"github.com/user/crontab-lint/internal/matcher"
)

func ts(year, month, day, hour, min int) time.Time {
	return time.Date(year, time.Month(month), day, hour, min, 0, 0, time.UTC)
}

func TestMatch_EveryMinute(t *testing.T) {
	res, err := matcher.Match("* * * * *", ts(2024, 1, 15, 10, 30))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Matches {
		t.Error("expected match for wildcard expression")
	}
}

func TestMatch_SpecificTime_Matches(t *testing.T) {
	res, err := matcher.Match("30 10 15 1 *", ts(2024, 1, 15, 10, 30))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Matches {
		t.Error("expected match")
	}
}

func TestMatch_SpecificTime_NoMatch(t *testing.T) {
	res, err := matcher.Match("0 12 * * *", ts(2024, 1, 15, 10, 30))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Matches {
		t.Error("expected no match")
	}
}

func TestMatch_FieldDetails(t *testing.T) {
	res, err := matcher.Match("30 10 * * *", ts(2024, 1, 15, 10, 30))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Fields[0].Matches {
		t.Error("minute field should match")
	}
	if !res.Fields[1].Matches {
		t.Error("hour field should match")
	}
}

func TestMatch_Alias(t *testing.T) {
	// @hourly normalizes to "0 * * * *"
	res, err := matcher.Match("@hourly", ts(2024, 6, 1, 8, 0))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Matches {
		t.Error("expected @hourly to match at minute 0")
	}
}

func TestMatch_StepExpression(t *testing.T) {
	res, err := matcher.Match("*/15 * * * *", ts(2024, 1, 1, 0, 45))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Matches {
		t.Error("expected match at minute 45 for */15")
	}
}

func TestMatch_InvalidExpression(t *testing.T) {
	_, err := matcher.Match("invalid", ts(2024, 1, 1, 0, 0))
	if err == nil {
		t.Error("expected error for invalid expression")
	}
}
