package formatter_test

import (
	"strings"
	"testing"

	"github.com/user/crontab-lint/internal/formatter"
)

func TestFormatText_Valid(t *testing.T) {
	r := formatter.Report{
		Expression:  "0 * * * *",
		Valid:        true,
		Explanation:  "At minute 0 of every hour",
		Issues:       nil,
	}

	out := formatter.FormatText(r)

	if !strings.Contains(out, "0 * * * *") {
		t.Errorf("expected expression in output, got: %s", out)
	}
	if !strings.Contains(out, "true") {
		t.Errorf("expected valid=true in output, got: %s", out)
	}
	if !strings.Contains(out, "At minute 0") {
		t.Errorf("expected explanation in output, got: %s", out)
	}
	if !strings.Contains(out, "none") {
		t.Errorf("expected 'none' issues in output, got: %s", out)
	}
}

func TestFormatText_WithIssues(t *testing.T) {
	r := formatter.Report{
		Expression: "99 * * * *",
		Valid:       false,
		Issues: []formatter.Issue{
			{Field: "minute", Message: "value 99 out of range [0-59]", Severity: "error"},
		},
	}

	out := formatter.FormatText(r)

	if !strings.Contains(out, "[ERROR]") {
		t.Errorf("expected [ERROR] tag in output, got: %s", out)
	}
	if !strings.Contains(out, "minute") {
		t.Errorf("expected field name in output, got: %s", out)
	}
	if !strings.Contains(out, "out of range") {
		t.Errorf("expected issue message in output, got: %s", out)
	}
}

func TestFormatJSON_Valid(t *testing.T) {
	r := formatter.Report{
		Expression:  "*/5 * * * *",
		Valid:        true,
		Explanation:  "Every 5 minutes",
		Issues:       nil,
	}

	out := formatter.FormatJSON(r)

	if !strings.Contains(out, `"valid":true`) {
		t.Errorf("expected valid:true in JSON, got: %s", out)
	}
	if !strings.Contains(out, `"issues":[]`) {
		t.Errorf("expected empty issues array in JSON, got: %s", out)
	}
	if !strings.Contains(out, `"explanation":"Every 5 minutes"`) {
		t.Errorf("expected explanation in JSON, got: %s", out)
	}
}

func TestFormatJSON_WithIssues(t *testing.T) {
	r := formatter.Report{
		Expression: "60 * * * *",
		Valid:       false,
		Issues: []formatter.Issue{
			{Field: "minute", Message: "value 60 out of range", Severity: "error"},
		},
	}

	out := formatter.FormatJSON(r)

	if !strings.Contains(out, `"severity":"error"`) {
		t.Errorf("expected severity in JSON, got: %s", out)
	}
	if !strings.Contains(out, `"field":"minute"`) {
		t.Errorf("expected field in JSON, got: %s", out)
	}
}
