package expander

import (
	"strings"
	"testing"
)

func TestExpand_EveryMinute(t *testing.T) {
	r := Expand("* * * * *")
	if len(r.Errors) != 0 {
		t.Fatalf("unexpected errors: %v", r.Errors)
	}
	if len(r.Fields) != 5 {
		t.Fatalf("expected 5 fields, got %d", len(r.Fields))
	}
	// minute field should have 60 values
	if len(r.Fields[0].Values) != 60 {
		t.Errorf("minute: expected 60 values, got %d", len(r.Fields[0].Values))
	}
}

func TestExpand_SpecificValues(t *testing.T) {
	r := Expand("0,30 9 * * *")
	if len(r.Errors) != 0 {
		t.Fatalf("unexpected errors: %v", r.Errors)
	}
	if got := r.Fields[0].Values; len(got) != 2 || got[0] != 0 || got[1] != 30 {
		t.Errorf("minute: expected [0 30], got %v", got)
	}
	if got := r.Fields[1].Values; len(got) != 1 || got[0] != 9 {
		t.Errorf("hour: expected [9], got %v", got)
	}
}

func TestExpand_StepField(t *testing.T) {
	r := Expand("*/15 * * * *")
	if len(r.Errors) != 0 {
		t.Fatalf("unexpected errors: %v", r.Errors)
	}
	expected := []int{0, 15, 30, 45}
	got := r.Fields[0].Values
	if len(got) != len(expected) {
		t.Fatalf("minute: expected %v, got %v", expected, got)
	}
	for i, v := range expected {
		if got[i] != v {
			t.Errorf("minute[%d]: expected %d, got %d", i, v, got[i])
		}
	}
}

func TestExpand_InvalidFieldCount(t *testing.T) {
	r := Expand("* * *")
	if len(r.Errors) == 0 {
		t.Fatal("expected error for wrong field count")
	}
	if !strings.Contains(r.Errors[0], "expected 5 fields") {
		t.Errorf("unexpected error message: %s", r.Errors[0])
	}
}

func TestExpand_RangeField(t *testing.T) {
	r := Expand("0 8-10 * * *")
	if len(r.Errors) != 0 {
		t.Fatalf("unexpected errors: %v", r.Errors)
	}
	expected := []int{8, 9, 10}
	got := r.Fields[1].Values
	if len(got) != len(expected) {
		t.Fatalf("hour: expected %v, got %v", expected, got)
	}
}

func TestFormatText_ContainsFields(t *testing.T) {
	r := Expand("0 0 * * *")
	out := FormatText(r)
	if !strings.Contains(out, "minute") {
		t.Error("expected 'minute' in text output")
	}
	if !strings.Contains(out, "hour") {
		t.Error("expected 'hour' in text output")
	}
}

func TestFormatJSON_ValidStructure(t *testing.T) {
	r := Expand("0 0 * * *")
	out, err := FormatJSON(r)
	if err != nil {
		t.Fatalf("FormatJSON error: %v", err)
	}
	if !strings.Contains(out, `"expression"`) {
		t.Error("expected 'expression' key in JSON")
	}
	if !strings.Contains(out, `"fields"`) {
		t.Error("expected 'fields' key in JSON")
	}
}

func TestJoinValues(t *testing.T) {
	got := JoinValues([]int{1, 2, 3})
	if got != "1,2,3" {
		t.Errorf("expected '1,2,3', got %q", got)
	}
}
