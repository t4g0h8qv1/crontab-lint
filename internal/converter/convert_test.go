package converter

import (
	"strings"
	"testing"
)

func TestToSystemd_EveryMinute(t *testing.T) {
	res := ToSystemd("* * * * *")
	if len(res.Errors) != 0 {
		t.Fatalf("unexpected errors: %v", res.Errors)
	}
	if res.Output != "*-*-* *:*:00" {
		t.Errorf("got %q, want %q", res.Output, "*-*-* *:*:00")
	}
}

func TestToSystemd_HourlyAt30(t *testing.T) {
	res := ToSystemd("30 * * * *")
	if len(res.Errors) != 0 {
		t.Fatalf("unexpected errors: %v", res.Errors)
	}
	if res.Output != "*-*-* *:30:00" {
		t.Errorf("got %q, want %q", res.Output, "*-*-* *:30:00")
	}
}

func TestToSystemd_DailyMidnight(t *testing.T) {
	res := ToSystemd("0 0 * * *")
	if len(res.Errors) != 0 {
		t.Fatalf("unexpected errors: %v", res.Errors)
	}
	if res.Output != "*-*-* 0:0:00" {
		t.Errorf("got %q, want %q", res.Output, "*-*-* 0:0:00")
	}
}

func TestToSystemd_WithDow(t *testing.T) {
	res := ToSystemd("0 9 * * 1")
	if len(res.Errors) != 0 {
		t.Fatalf("unexpected errors: %v", res.Errors)
	}
	if !strings.HasPrefix(res.Output, "Mon ") {
		t.Errorf("expected output to start with 'Mon ', got %q", res.Output)
	}
}

func TestToSystemd_InvalidExpression(t *testing.T) {
	res := ToSystemd("99 * * * *")
	if len(res.Errors) == 0 {
		t.Error("expected errors for invalid expression, got none")
	}
	if res.Output != "" {
		t.Errorf("expected empty output on error, got %q", res.Output)
	}
}

func TestToSystemd_WrongFieldCount(t *testing.T) {
	res := ToSystemd("* * *")
	if len(res.Errors) == 0 {
		t.Error("expected errors for wrong field count")
	}
}

func TestToSystemd_FormatField(t *testing.T) {
	res := ToSystemd("0 12 15 6 *")
	if len(res.Errors) != 0 {
		t.Fatalf("unexpected errors: %v", res.Errors)
	}
	if res.Format != "systemd" {
		t.Errorf("expected format 'systemd', got %q", res.Format)
	}
	if res.Expression != "0 12 15 6 *" {
		t.Errorf("expected expression preserved, got %q", res.Expression)
	}
}

func TestToSystemd_Alias(t *testing.T) {
	res := ToSystemd("@daily")
	// @daily normalizes to '0 0 * * *'
	if len(res.Errors) != 0 {
		t.Fatalf("unexpected errors for @daily: %v", res.Errors)
	}
	if res.Output == "" {
		t.Error("expected non-empty output for @daily")
	}
}
