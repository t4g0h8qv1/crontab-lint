package ranker_test

import (
	"strings"
	"testing"

	"github.com/user/crontab-lint/internal/ranker"
)

func TestRank_SimpleWildcard(t *testing.T) {
	s, err := ranker.Rank("* * * * *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Complexity != 0 {
		t.Errorf("expected complexity 0, got %d", s.Complexity)
	}
	if s.Total <= 0 {
		t.Errorf("expected positive total, got %d", s.Total)
	}
}

func TestRank_ComplexExpression(t *testing.T) {
	s, err := ranker.Rank("1,15,30 0-6 */2 1,6,12 1-5")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Complexity <= 0 {
		t.Errorf("expected non-zero complexity, got %d", s.Complexity)
	}
	if len(s.Notes) == 0 {
		t.Error("expected notes for complex expression")
	}
}

func TestRank_InvalidFieldCount(t *testing.T) {
	_, err := ranker.Rank("* * *")
	if err == nil {
		t.Error("expected error for wrong field count")
	}
}

func TestFormatText_ContainsScore(t *testing.T) {
	s, _ := ranker.Rank("0 12 * * *")
	out := ranker.FormatText(s)
	if !strings.Contains(out, "Total Score") {
		t.Error("expected 'Total Score' in text output")
	}
	if !strings.Contains(out, s.Expression) {
		t.Error("expected expression in text output")
	}
}

func TestFormatJSON_ValidJSON(t *testing.T) {
	s, _ := ranker.Rank("0 0 * * 0")
	out, err := ranker.FormatJSON(s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "\"total\"") {
		t.Error("expected 'total' key in JSON output")
	}
}

func TestRank_StepField(t *testing.T) {
	s, err := ranker.Rank("*/5 * * * *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Complexity < 1 {
		t.Errorf("expected complexity >= 1 for step field, got %d", s.Complexity)
	}
}
