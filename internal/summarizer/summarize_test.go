package summarizer_test

import (
	"strings"
	"testing"

	"github.com/user/crontab-lint/internal/summarizer"
)

func TestSummarize_EveryMinute(t *testing.T) {
	res := summarizer.Summarize("* * * * *")
	if len(res.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", res.Errors)
	}
	if !strings.Contains(res.Summary, "every minute") {
		t.Errorf("expected 'every minute' in summary, got: %s", res.Summary)
	}
}

func TestSummarize_Hourly(t *testing.T) {
	res := summarizer.Summarize("30 * * * *")
	if len(res.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", res.Errors)
	}
	if !strings.Contains(res.Summary, "every hour") {
		t.Errorf("expected 'every hour' in summary, got: %s", res.Summary)
	}
	if !strings.Contains(res.Summary, "30") {
		t.Errorf("expected minute '30' in summary, got: %s", res.Summary)
	}
}

func TestSummarize_Daily(t *testing.T) {
	res := summarizer.Summarize("0 9 * * *")
	if len(res.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", res.Errors)
	}
	if !strings.Contains(res.Summary, "daily") {
		t.Errorf("expected 'daily' in summary, got: %s", res.Summary)
	}
}

func TestSummarize_Weekly(t *testing.T) {
	res := summarizer.Summarize("0 8 * * 1")
	if len(res.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", res.Errors)
	}
	if !strings.Contains(res.Summary, "weekly") {
		t.Errorf("expected 'weekly' in summary, got: %s", res.Summary)
	}
}

func TestSummarize_InvalidExpression(t *testing.T) {
	res := summarizer.Summarize("invalid")
	if len(res.Errors) == 0 {
		t.Error("expected errors for invalid expression, got none")
	}
	if res.Summary != "invalid expression" {
		t.Errorf("expected 'invalid expression', got: %s", res.Summary)
	}
}

func TestSummarize_SetsNormalized(t *testing.T) {
	res := summarizer.Summarize("@daily")
	if len(res.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", res.Errors)
	}
	if res.Normalized == "" {
		t.Error("expected normalized field to be set")
	}
}

func TestSummarize_SetsClass(t *testing.T) {
	res := summarizer.Summarize("* * * * *")
	if res.Class == "" {
		t.Error("expected class to be set")
	}
}
