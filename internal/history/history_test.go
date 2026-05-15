package history_test

import (
	"strings"
	"testing"

	"github.com/user/crontab-lint/internal/history"
)

func TestHistory_AddAndLen(t *testing.T) {
	h := history.New(10)
	if h.Len() != 0 {
		t.Fatalf("expected empty history, got %d", h.Len())
	}
	h.Add("* * * * *", true, nil)
	h.Add("0 * * * *", true, nil)
	if h.Len() != 2 {
		t.Fatalf("expected 2 entries, got %d", h.Len())
	}
}

func TestHistory_Eviction(t *testing.T) {
	h := history.New(3)
	h.Add("a", true, nil)
	h.Add("b", true, nil)
	h.Add("c", true, nil)
	h.Add("d", true, nil)
	if h.Len() != 3 {
		t.Fatalf("expected 3 entries after eviction, got %d", h.Len())
	}
	entries := h.All()
	if entries[0].Expression != "b" {
		t.Errorf("expected oldest entry to be 'b', got %q", entries[0].Expression)
	}
	if entries[2].Expression != "d" {
		t.Errorf("expected newest entry to be 'd', got %q", entries[2].Expression)
	}
}

func TestHistory_Last(t *testing.T) {
	h := history.New(10)
	if h.Last() != nil {
		t.Fatal("expected nil for empty history")
	}
	h.Add("0 12 * * *", true, nil)
	h.Add("bad expr", false, []string{"invalid field count"})
	last := h.Last()
	if last == nil {
		t.Fatal("expected non-nil last entry")
	}
	if last.Expression != "bad expr" {
		t.Errorf("expected last expression 'bad expr', got %q", last.Expression)
	}
	if last.Valid {
		t.Error("expected last entry to be invalid")
	}
}

func TestHistory_Clear(t *testing.T) {
	h := history.New(10)
	h.Add("* * * * *", true, nil)
	h.Clear()
	if h.Len() != 0 {
		t.Fatalf("expected empty history after clear, got %d", h.Len())
	}
	if h.Last() != nil {
		t.Fatal("expected nil last after clear")
	}
}

func TestHistory_Format_Empty(t *testing.T) {
	h := history.New(10)
	out := h.Format()
	if out != "No history recorded." {
		t.Errorf("unexpected empty format output: %q", out)
	}
}

func TestHistory_Format_WithEntries(t *testing.T) {
	h := history.New(10)
	h.Add("* * * * *", true, nil)
	h.Add("0 * * * *", true, []string{"some warning"})
	h.Add("bad", false, []string{"invalid"})
	out := h.Format()
	if !strings.Contains(out, "* * * * *") {
		t.Error("expected expression in format output")
	}
	if !strings.Contains(out, "INVALID") {
		t.Error("expected INVALID status in format output")
	}
	if !strings.Contains(out, "WARN(1)") {
		t.Error("expected WARN status in format output")
	}
}

func TestHistory_DefaultMaxSize(t *testing.T) {
	h := history.New(0)
	for i := 0; i < 60; i++ {
		h.Add("* * * * *", true, nil)
	}
	if h.Len() != 50 {
		t.Errorf("expected default max size 50, got %d", h.Len())
	}
}
