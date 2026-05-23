package heatmap

import (
	"strings"
	"testing"
	"time"
)

// fixed origin: Monday 2024-01-01 00:00:00 UTC
var origin = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func TestCompute_EveryMinute(t *testing.T) {
	r := Compute("* * * * *", origin, 10_080)
	if len(r.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", r.Errors)
	}
	if len(r.Cells) != 7*24 {
		t.Fatalf("expected 168 cells, got %d", len(r.Cells))
	}
	// Every hour should have the same count (60 firings).
	for _, c := range r.Cells {
		if c.Count != 60 {
			t.Errorf("DOW=%d Hour=%d: want 60, got %d", c.DOW, c.Hour, c.Count)
		}
	}
	if r.MaxCount != 60 {
		t.Errorf("MaxCount: want 60, got %d", r.MaxCount)
	}
}

func TestCompute_HourlyAt30(t *testing.T) {
	r := Compute("30 * * * *", origin, 10_080)
	if len(r.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", r.Errors)
	}
	for _, c := range r.Cells {
		if c.Count != 1 {
			t.Errorf("DOW=%d Hour=%d: want 1, got %d", c.DOW, c.Hour, c.Count)
		}
	}
}

func TestCompute_InvalidExpression(t *testing.T) {
	r := Compute("invalid", origin, 100)
	if len(r.Errors) == 0 {
		t.Fatal("expected errors for invalid expression")
	}
	if len(r.Cells) != 0 {
		t.Errorf("expected no cells on error, got %d", len(r.Cells))
	}
}

func TestFormatText_ContainsHeader(t *testing.T) {
	r := Compute("* * * * *", origin, 10_080)
	out := FormatText(r)
	if !strings.Contains(out, "Heatmap for:") {
		t.Error("FormatText missing 'Heatmap for:' header")
	}
	for _, d := range dayNames {
		if !strings.Contains(out, d) {
			t.Errorf("FormatText missing day label %q", d)
		}
	}
}

func TestFormatJSON_ValidStructure(t *testing.T) {
	r := Compute("0 9 * * 1", origin, 10_080)
	out, err := FormatJSON(r)
	if err != nil {
		t.Fatalf("FormatJSON error: %v", err)
	}
	if !strings.Contains(out, `"expression"`) {
		t.Error("JSON missing 'expression' key")
	}
	if !strings.Contains(out, `"cells"`) {
		t.Error("JSON missing 'cells' key")
	}
	if !strings.Contains(out, `"max_count"`) {
		t.Error("JSON missing 'max_count' key")
	}
}

func TestFormatText_ErrorResult(t *testing.T) {
	r := Result{Expression: "bad", Errors: []string{"scheduler error: invalid"}}
	out := FormatText(r)
	if !strings.Contains(out, "heatmap error") {
		t.Errorf("expected error text, got: %s", out)
	}
}
