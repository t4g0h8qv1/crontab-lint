package suggester_test

import (
	"testing"

	"github.com/user/crontab-lint/internal/suggester"
)

func TestSuggest_ValidExpression(t *testing.T) {
	suggestions := suggester.Suggest("0 12 * * 1")
	if len(suggestions) != 0 {
		t.Errorf("expected no suggestions, got %d: %+v", len(suggestions), suggestions)
	}
}

func TestSuggest_WrongFieldCount(t *testing.T) {
	suggestions := suggester.Suggest("0 12 *")
	if len(suggestions) != 1 {
		t.Fatalf("expected 1 suggestion, got %d", len(suggestions))
	}
	if suggestions[0].Field != "expression" {
		t.Errorf("expected field 'expression', got %q", suggestions[0].Field)
	}
}

func TestSuggest_RedundantStep(t *testing.T) {
	suggestions := suggester.Suggest("*/1 * * * *")
	if len(suggestions) == 0 {
		t.Fatal("expected at least one suggestion for */1")
	}
	found := false
	for _, s := range suggestions {
		if s.Field == "minute" && s.Issue == "*/1 is redundant" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected redundant */1 suggestion for minute field, got: %+v", suggestions)
	}
}

func TestSuggest_OutOfRangeMinute(t *testing.T) {
	suggestions := suggester.Suggest("99 * * * *")
	if len(suggestions) == 0 {
		t.Fatal("expected suggestion for out-of-range minute")
	}
	if suggestions[0].Field != "minute" {
		t.Errorf("expected field 'minute', got %q", suggestions[0].Field)
	}
}

func TestSuggest_DomDowConflict(t *testing.T) {
	suggestions := suggester.Suggest("0 12 15 * 1")
	found := false
	for _, s := range suggestions {
		if s.Field == "day-of-month/day-of-week" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected DOM/DOW conflict suggestion, got: %+v", suggestions)
	}
}

func TestSuggest_OutOfRangeHour(t *testing.T) {
	suggestions := suggester.Suggest("0 25 * * *")
	if len(suggestions) == 0 {
		t.Fatal("expected suggestion for out-of-range hour")
	}
	if suggestions[0].Field != "hour" {
		t.Errorf("expected field 'hour', got %q", suggestions[0].Field)
	}
}
