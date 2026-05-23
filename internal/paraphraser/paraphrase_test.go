package paraphraser_test

import (
	"strings"
	"testing"

	"github.com/user/crontab-lint/internal/paraphraser"
)

func TestParaphrase_EveryMinute(t *testing.T) {
	r := paraphraser.Paraphrase("* * * * *")
	if len(r.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", r.Errors)
	}
	if r.Paraphrase != "every minute" {
		t.Errorf("expected 'every minute', got %q", r.Paraphrase)
	}
}

func TestParaphrase_HourlyAt30(t *testing.T) {
	r := paraphraser.Paraphrase("30 * * * *")
	if len(r.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", r.Errors)
	}
	if !strings.Contains(r.Paraphrase, "30 minute") {
		t.Errorf("expected minute reference, got %q", r.Paraphrase)
	}
}

func TestParaphrase_DailyAtMidnight(t *testing.T) {
	r := paraphraser.Paraphrase("0 0 * * *")
	if len(r.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", r.Errors)
	}
	if !strings.Contains(r.Paraphrase, "0 minute") {
		t.Errorf("expected minute reference, got %q", r.Paraphrase)
	}
	if !strings.Contains(r.Paraphrase, "0 hour") {
		t.Errorf("expected hour reference, got %q", r.Paraphrase)
	}
}

func TestParaphrase_StepField(t *testing.T) {
	r := paraphraser.Paraphrase("*/15 * * * *")
	if len(r.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", r.Errors)
	}
	if !strings.Contains(r.Paraphrase, "every 15 minutes") {
		t.Errorf("expected step description, got %q", r.Paraphrase)
	}
}

func TestParaphrase_Alias(t *testing.T) {
	r := paraphraser.Paraphrase("@daily")
	if len(r.Errors) > 0 {
		t.Fatalf("unexpected errors for @daily: %v", r.Errors)
	}
	if r.Paraphrase == "" {
		t.Error("expected non-empty paraphrase for @daily")
	}
}

func TestParaphrase_InvalidExpression(t *testing.T) {
	r := paraphraser.Paraphrase("not a cron")
	if len(r.Errors) == 0 {
		t.Error("expected errors for invalid expression")
	}
	if r.Expression != "not a cron" {
		t.Errorf("expression not preserved, got %q", r.Expression)
	}
}

func TestParaphrase_WithDow(t *testing.T) {
	r := paraphraser.Paraphrase("0 9 * * 1")
	if len(r.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", r.Errors)
	}
	if !strings.Contains(r.Paraphrase, "1 day-of-week") {
		t.Errorf("expected day-of-week reference, got %q", r.Paraphrase)
	}
}

func TestParaphrase_ResultPreservesExpression(t *testing.T) {
	expr := "5 4 * * *"
	r := paraphraser.Paraphrase(expr)
	if r.Expression != expr {
		t.Errorf("expected expression %q, got %q", expr, r.Expression)
	}
}
