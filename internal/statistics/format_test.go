package statistics_test

import (
	"strings"
	"testing"

	"github.com/user/crontab-lint/internal/statistics"
)

func sampleStats() statistics.Stats {
	s, _ := statistics.Compute("*/10 * * * *", epoch, 20)
	return s
}

func TestFormatText_ContainsExpression(t *testing.T) {
	s := sampleStats()
	out := statistics.FormatText(s)
	if !strings.Contains(out, "*/10 * * * *") {
		t.Errorf("FormatText output missing expression:\n%s", out)
	}
}

func TestFormatText_ContainsKeys(t *testing.T) {
	s := sampleStats()
	out := statistics.FormatText(s)
	for _, key := range []string{"Fires/hour", "Fires/day", "Fires/week", "Min gap", "Max gap"} {
		if !strings.Contains(out, key) {
			t.Errorf("FormatText missing key %q", key)
		}
	}
}

func TestFormatJSON_ValidJSON(t *testing.T) {
	s := sampleStats()
	out, err := statistics.FormatJSON(s)
	if err != nil {
		t.Fatalf("FormatJSON error: %v", err)
	}
	if !strings.Contains(out, `"expression"`) {
		t.Errorf("JSON missing expression key:\n%s", out)
	}
	if !strings.Contains(out, `"fires_per_hour"`) {
		t.Errorf("JSON missing fires_per_hour key:\n%s", out)
	}
}

func TestFormatJSON_MinIntervalString(t *testing.T) {
	s := sampleStats()
	out, err := statistics.FormatJSON(s)
	if err != nil {
		t.Fatalf("FormatJSON error: %v", err)
	}
	if !strings.Contains(out, `"min_interval"`) {
		t.Errorf("JSON missing min_interval:\n%s", out)
	}
}
