package deduplicator

import (
	"strings"
	"testing"
)

func TestDedup_NoDuplicates(t *testing.T) {
	r, err := Dedup("* * * * *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Changed {
		t.Errorf("expected no change, got %q", r.Deduped)
	}
}

func TestDedup_DuplicateMinutes(t *testing.T) {
	r, err := Dedup("0,30,0 * * * *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !r.Changed {
		t.Fatal("expected a change")
	}
	if r.Deduped != "0,30 * * * *" {
		t.Errorf("got %q", r.Deduped)
	}
	if len(r.Changes) != 1 {
		t.Errorf("expected 1 change entry, got %d", len(r.Changes))
	}
}

func TestDedup_SortsNumericList(t *testing.T) {
	r, err := Dedup("30,0,15 * * * *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Deduped != "0,15,30 * * * *" {
		t.Errorf("expected sorted output, got %q", r.Deduped)
	}
}

func TestDedup_MultipleDuplicates(t *testing.T) {
	r, err := Dedup("1,2,1,3,2 * * * *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Deduped != "1,2,3 * * * *" {
		t.Errorf("got %q", r.Deduped)
	}
	if len(r.Changes) != 2 {
		t.Errorf("expected 2 change entries, got %d", len(r.Changes))
	}
}

func TestDedup_WrongFieldCount(t *testing.T) {
	_, err := Dedup("* * *")
	if err == nil {
		t.Fatal("expected error for wrong field count")
	}
}

func TestDedup_NonListFieldPassthrough(t *testing.T) {
	r, err := Dedup("*/5 0-23 * * 1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Changed {
		t.Errorf("expected no change, got %q", r.Deduped)
	}
}

func TestDedup_ChangesContainFieldIndex(t *testing.T) {
	r, err := Dedup("0,0 * * * *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Changes) == 0 {
		t.Fatal("expected at least one change")
	}
	if !strings.HasPrefix(r.Changes[0], "field 1:") {
		t.Errorf("expected field index prefix, got %q", r.Changes[0])
	}
}
