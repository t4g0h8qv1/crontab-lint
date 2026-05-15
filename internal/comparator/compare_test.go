package comparator_test

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/user/crontab-lint/internal/comparator"
)

var baseTime = time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

func TestCompare_EqualExpressions(t *testing.T) {
	r, err := comparator.Compare("* * * * *", "* * * * *", baseTime)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Relation != "equal" {
		t.Errorf("expected equal, got %s", r.Relation)
	}
	if r.Delta != 0 {
		t.Errorf("expected delta 0, got %d", r.Delta)
	}
}

func TestCompare_BFasterThanA(t *testing.T) {
	// A fires once/hour (24/day), B fires every minute (1440/day)
	r, err := comparator.Compare("0 * * * *", "* * * * *", baseTime)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Relation != "faster" {
		t.Errorf("expected faster, got %s", r.Relation)
	}
	if r.FrequencyB <= r.FrequencyA {
		t.Errorf("expected B (%d) > A (%d)", r.FrequencyB, r.FrequencyA)
	}
	if r.Delta <= 0 {
		t.Errorf("expected positive delta, got %d", r.Delta)
	}
}

func TestCompare_BSlowerThanA(t *testing.T) {
	r, err := comparator.Compare("* * * * *", "0 * * * *", baseTime)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Relation != "slower" {
		t.Errorf("expected slower, got %s", r.Relation)
	}
	if r.Delta >= 0 {
		t.Errorf("expected negative delta, got %d", r.Delta)
	}
}

func TestCompare_InvalidExpressionA(t *testing.T) {
	_, err := comparator.Compare("invalid", "* * * * *", baseTime)
	if err == nil {
		t.Fatal("expected error for invalid expression A")
	}
	if !strings.Contains(err.Error(), "expression A") {
		t.Errorf("error should mention expression A, got: %v", err)
	}
}

func TestFormatText_ContainsRelation(t *testing.T) {
	r, _ := comparator.Compare("0 * * * *", "* * * * *", baseTime)
	out := comparator.FormatText(r)
	if !strings.Contains(out, "faster") {
		t.Errorf("expected 'faster' in output, got: %s", out)
	}
}

func TestFormatJSON_ValidStructure(t *testing.T) {
	r, _ := comparator.Compare("0 * * * *", "0 0 * * *", baseTime)
	out, err := comparator.FormatJSON(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(out), &m); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	for _, key := range []string{"expression_a", "expression_b", "relation", "delta"} {
		if _, ok := m[key]; !ok {
			t.Errorf("missing key %q in JSON output", key)
		}
	}
}
