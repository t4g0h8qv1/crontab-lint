package highlighter_test

import (
	"strings"
	"testing"

	"github.com/user/crontab-lint/internal/highlighter"
)

func TestHighlight_Valid(t *testing.T) {
	result, err := highlighter.Highlight("*/5 * * * *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Stripped output should match original fields joined by spaces.
	stripped := highlighter.Strip(result.Highlighted)
	if stripped != "*/5 * * * *" {
		t.Errorf("stripped highlighted = %q, want %q", stripped, "*/5 * * * *")
	}

	// Highlighted string must contain ANSI codes.
	if result.Highlighted == stripped {
		t.Error("expected ANSI codes in highlighted output, got plain text")
	}
}

func TestHighlight_LegendKeys(t *testing.T) {
	result, err := highlighter.Highlight("0 12 * * 1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedKeys := []string{"minute", "hour", "dom", "month", "dow"}
	for _, key := range expectedKeys {
		if _, ok := result.Legend[key]; !ok {
			t.Errorf("legend missing key %q", key)
		}
	}
}

func TestHighlight_WrongFieldCount(t *testing.T) {
	_, err := highlighter.Highlight("* * *")
	if err == nil {
		t.Fatal("expected error for wrong field count, got nil")
	}
	if !strings.Contains(err.Error(), "expected 5 fields") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestHighlight_EmptyExpression(t *testing.T) {
	_, err := highlighter.Highlight("")
	if err == nil {
		t.Fatal("expected error for empty expression, got nil")
	}
}

func TestStrip_RemovesAllCodes(t *testing.T) {
	input := "\033[36mminute\033[0m \033[32mhour\033[0m"
	want := "minute hour"
	got := highlighter.Strip(input)
	if got != want {
		t.Errorf("Strip() = %q, want %q", got, want)
	}
}

func TestStrip_PlainText(t *testing.T) {
	input := "* * * * *"
	got := highlighter.Strip(input)
	if got != input {
		t.Errorf("Strip() changed plain text: got %q, want %q", got, input)
	}
}
