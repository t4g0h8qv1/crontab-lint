package matcher_test

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/user/crontab-lint/internal/matcher"
)

func TestFormatText_Match(t *testing.T) {
	res, err := matcher.Match("* * * * *", ts(2024, 3, 10, 14, 5))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := matcher.FormatText(res)
	if !strings.Contains(out, "MATCH") {
		t.Error("expected MATCH in output")
	}
	if !strings.Contains(out, "minute") {
		t.Error("expected field name 'minute' in output")
	}
}

func TestFormatText_NoMatch(t *testing.T) {
	res, err := matcher.Match("0 0 * * *", ts(2024, 3, 10, 14, 5))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := matcher.FormatText(res)
	if !strings.Contains(out, "NO MATCH") {
		t.Error("expected NO MATCH in output")
	}
	if !strings.Contains(out, "✗") {
		t.Error("expected ✗ symbol for non-matching field")
	}
}

func TestFormatJSON_ValidStructure(t *testing.T) {
	res, err := matcher.Match("30 6 * * 1", time.Date(2024, 1, 8, 6, 30, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out, err := matcher.FormatJSON(res)
	if err != nil {
		t.Fatalf("FormatJSON error: %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(out), &parsed); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if _, ok := parsed["matches"]; !ok {
		t.Error("expected 'matches' key in JSON")
	}
	if _, ok := parsed["fields"]; !ok {
		t.Error("expected 'fields' key in JSON")
	}
}

func TestFormatJSON_MatchTrue(t *testing.T) {
	res, err := matcher.Match("* * * * *", ts(2024, 1, 1, 0, 0))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out, err := matcher.FormatJSON(res)
	if err != nil {
		t.Fatalf("FormatJSON error: %v", err)
	}
	if !strings.Contains(out, `"matches": true`) {
		t.Error("expected matches:true in JSON output")
	}
}
