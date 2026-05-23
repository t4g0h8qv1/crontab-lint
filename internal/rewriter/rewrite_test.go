package rewriter_test

import (
	"testing"

	"github.com/user/crontab-lint/internal/rewriter"
)

func TestRewrite_NoReplacements(t *testing.T) {
	expr := "*/5 * * * *"
	res := rewriter.Rewrite(expr, rewriter.Replacements{})
	if res.Rewritten != expr {
		t.Errorf("expected %q, got %q", expr, res.Rewritten)
	}
	if len(res.Changes) != 0 {
		t.Errorf("expected no changes, got %d", len(res.Changes))
	}
}

func TestRewrite_SingleField(t *testing.T) {
	res := rewriter.Rewrite("*/5 * * * *", rewriter.Replacements{
		"minute": "0",
	})
	if res.Rewritten != "0 * * * *" {
		t.Errorf("unexpected rewrite: %q", res.Rewritten)
	}
	if len(res.Changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(res.Changes))
	}
	if res.Changes[0].Field != "minute" {
		t.Errorf("expected field 'minute', got %q", res.Changes[0].Field)
	}
	if res.Changes[0].Original != "*/5" {
		t.Errorf("expected original '*/5', got %q", res.Changes[0].Original)
	}
	if res.Changes[0].Rewritten != "0" {
		t.Errorf("expected rewritten '0', got %q", res.Changes[0].Rewritten)
	}
}

func TestRewrite_MultipleFields(t *testing.T) {
	res := rewriter.Rewrite("0 3 * * 1", rewriter.Replacements{
		"hour": "6",
		"dow":  "5",
	})
	if res.Rewritten != "0 6 * * 5" {
		t.Errorf("unexpected rewrite: %q", res.Rewritten)
	}
	if len(res.Changes) != 2 {
		t.Errorf("expected 2 changes, got %d", len(res.Changes))
	}
}

func TestRewrite_SameValueNoChange(t *testing.T) {
	res := rewriter.Rewrite("0 * * * *", rewriter.Replacements{
		"minute": "0",
	})
	if len(res.Changes) != 0 {
		t.Errorf("expected no changes when value is identical, got %d", len(res.Changes))
	}
}

func TestRewrite_InvalidExpression(t *testing.T) {
	res := rewriter.Rewrite("bad expr", rewriter.Replacements{
		"minute": "0",
	})
	if len(res.Errors) == 0 {
		t.Error("expected errors for invalid expression")
	}
	if res.Rewritten != "bad expr" {
		t.Errorf("expected original expression on error, got %q", res.Rewritten)
	}
}

func TestRewrite_EmptyReplacement(t *testing.T) {
	res := rewriter.Rewrite("* * * * *", rewriter.Replacements{
		"hour": "",
	})
	if len(res.Errors) == 0 {
		t.Error("expected error for empty replacement value")
	}
}
