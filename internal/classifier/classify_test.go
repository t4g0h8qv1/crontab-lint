package classifier

import (
	"testing"
)

func TestClassify_EveryMinute(t *testing.T) {
	r := Classify("* * * * *")
	if r.Class != ClassEveryMinute {
		t.Errorf("expected every-minute, got %s", r.Class)
	}
}

func TestClassify_Hourly(t *testing.T) {
	r := Classify("0 * * * *")
	if r.Class != ClassHourly {
		t.Errorf("expected hourly, got %s", r.Class)
	}
	if r.Label != "hourly" {
		t.Errorf("unexpected label: %s", r.Label)
	}
}

func TestClassify_Daily(t *testing.T) {
	r := Classify("30 6 * * *")
	if r.Class != ClassDaily {
		t.Errorf("expected daily, got %s", r.Class)
	}
}

func TestClassify_Weekly(t *testing.T) {
	r := Classify("0 9 * * 1")
	if r.Class != ClassWeekly {
		t.Errorf("expected weekly, got %s", r.Class)
	}
}

func TestClassify_Monthly(t *testing.T) {
	r := Classify("0 0 1 * *")
	if r.Class != ClassMonthly {
		t.Errorf("expected monthly, got %s", r.Class)
	}
}

func TestClassify_Yearly(t *testing.T) {
	r := Classify("0 0 1 1 *")
	if r.Class != ClassYearly {
		t.Errorf("expected yearly, got %s", r.Class)
	}
}

func TestClassify_Custom(t *testing.T) {
	r := Classify("*/15 8-18 * * 1-5")
	if r.Class != ClassCustom {
		t.Errorf("expected custom, got %s", r.Class)
	}
	if r.Label != "custom schedule" {
		t.Errorf("unexpected label: %s", r.Label)
	}
}

func TestClassify_InvalidFieldCount(t *testing.T) {
	r := Classify("* * *")
	if r.Class != ClassUnknown {
		t.Errorf("expected unknown for invalid expression, got %s", r.Class)
	}
}

func TestClassify_ExpressionPreserved(t *testing.T) {
	expr := "5 4 * * 0"
	r := Classify(expr)
	if r.Expression != expr {
		t.Errorf("expression not preserved: got %s", r.Expression)
	}
}
