package deduplicator

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestFormatText_NoChange(t *testing.T) {
	r := Result{
		Original: "* * * * *",
		Deduped:  "* * * * *",
		Changed:  false,
	}
	out := FormatText(r)
	if !strings.Contains(out, "no duplicates found") {
		t.Errorf("expected no-duplicate message, got:\n%s", out)
	}
}

func TestFormatText_WithChanges(t *testing.T) {
	r := Result{
		Original: "0,0 * * * *",
		Deduped:  "0 * * * *",
		Changed:  true,
		Changes:  []string{"field 1: removed duplicate \"0\""},
	}
	out := FormatText(r)
	if !strings.Contains(out, "1 change(s)") {
		t.Errorf("expected change count, got:\n%s", out)
	}
	if !strings.Contains(out, "field 1") {
		t.Errorf("expected field reference, got:\n%s", out)
	}
}

func TestFormatJSON_ValidStructure(t *testing.T) {
	r := Result{
		Original: "0,0 * * * *",
		Deduped:  "0 * * * *",
		Changed:  true,
		Changes:  []string{"field 1: removed duplicate \"0\""},
	}
	out, err := FormatJSON(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(out), &m); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if _, ok := m["changed"]; !ok {
		t.Error("missing 'changed' key")
	}
	if _, ok := m["changes"]; !ok {
		t.Error("missing 'changes' key")
	}
}

func TestFormatJSON_EmptyChangesArray(t *testing.T) {
	r := Result{Original: "* * * * *", Deduped: "* * * * *", Changed: false}
	out, err := FormatJSON(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, `"changes": []`) {
		t.Errorf("expected empty array, got:\n%s", out)
	}
}
