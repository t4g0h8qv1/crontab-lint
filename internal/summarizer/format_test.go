package summarizer_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/crontab-lint/internal/summarizer"
)

func TestFormatText_ContainsSummary(t *testing.T) {
	res := summarizer.Summarize("* * * * *")
	out := summarizer.FormatText(res)
	if !strings.Contains(out, "Summary") {
		t.Errorf("expected 'Summary' label in output, got:\n%s", out)
	}
	if !strings.Contains(out, res.Summary) {
		t.Errorf("expected summary text in output, got:\n%s", out)
	}
}

func TestFormatText_ShowsErrors(t *testing.T) {
	res := summarizer.Summarize("bad expr")
	out := summarizer.FormatText(res)
	if !strings.Contains(out, "Errors") {
		t.Errorf("expected 'Errors' section in output, got:\n%s", out)
	}
}

func TestFormatText_ShowsNormalizedWhenDifferent(t *testing.T) {
	res := summarizer.Summarize("@daily")
	out := summarizer.FormatText(res)
	if res.Normalized != "" && res.Normalized != res.Expression {
		if !strings.Contains(out, "Normalized") {
			t.Errorf("expected 'Normalized' label when alias used, got:\n%s", out)
		}
	}
}

func TestFormatJSON_ValidStructure(t *testing.T) {
	res := summarizer.Summarize("0 12 * * *")
	out, err := summarizer.FormatJSON(res)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(out), &parsed); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	for _, key := range []string{"expression", "summary", "class"} {
		if _, ok := parsed[key]; !ok {
			t.Errorf("expected key %q in JSON output", key)
		}
	}
}

func TestFormatJSON_InvalidExpression(t *testing.T) {
	res := summarizer.Summarize("not-valid")
	out, err := summarizer.FormatJSON(res)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(out), &parsed); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if _, ok := parsed["errors"]; !ok {
		t.Error("expected 'errors' key in JSON for invalid expression")
	}
}
