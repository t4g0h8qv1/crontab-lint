package annotator

import (
	"strings"
	"testing"
)

func TestAnnotate_ValidExpression(t *testing.T) {
	r := Annotate("0 9 * * 1")
	if len(r.Errors) != 0 {
		t.Fatalf("unexpected errors: %v", r.Errors)
	}
	if len(r.Fields) != 5 {
		t.Fatalf("expected 5 fields, got %d", len(r.Fields))
	}
	if r.Fields[0].Name != "minute" {
		t.Errorf("first field name: got %q, want \"minute\"", r.Fields[0].Name)
	}
	if r.Fields[0].Value != "0" {
		t.Errorf("first field value: got %q, want \"0\"", r.Fields[0].Value)
	}
	if r.Fields[0].Note == "" {
		t.Error("expected non-empty note for minute field")
	}
}

func TestAnnotate_InvalidFieldCount(t *testing.T) {
	r := Annotate("* * *")
	if len(r.Errors) == 0 {
		t.Fatal("expected errors for wrong field count")
	}
	if len(r.Fields) != 0 {
		t.Errorf("expected no fields on error, got %d", len(r.Fields))
	}
}

func TestAnnotate_EveryMinute(t *testing.T) {
	r := Annotate("* * * * *")
	if len(r.Errors) != 0 {
		t.Fatalf("unexpected errors: %v", r.Errors)
	}
	for _, f := range r.Fields {
		if f.Note == "" {
			t.Errorf("field %q has empty note", f.Name)
		}
	}
}

func TestAnnotate_FieldNamesPresent(t *testing.T) {
	r := Annotate("30 6 * * *")
	expected := []string{"minute", "hour", "day-of-month", "month", "day-of-week"}
	for i, f := range r.Fields {
		if f.Name != expected[i] {
			t.Errorf("field[%d] name: got %q, want %q", i, f.Name, expected[i])
		}
	}
}

func TestInline_ValidExpression(t *testing.T) {
	out := Inline("0 0 * * *")
	if !strings.Contains(out, "0 0 * * *") {
		t.Errorf("output should contain original expression, got %q", out)
	}
	if !strings.Contains(out, "#") {
		t.Errorf("output should contain inline comment, got %q", out)
	}
	if !strings.Contains(out, "|") {
		t.Errorf("output should contain field separator '|', got %q", out)
	}
}

func TestInline_InvalidExpression(t *testing.T) {
	out := Inline("bad expression here")
	if !strings.Contains(out, "error:") {
		t.Errorf("expected error annotation, got %q", out)
	}
}
