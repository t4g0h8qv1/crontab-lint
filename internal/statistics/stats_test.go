package statistics_test

import (
	"testing"
	"time"

	"github.com/user/crontab-lint/internal/statistics"
)

var epoch = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func TestCompute_EveryMinute(t *testing.T) {
	s, err := statistics.Compute("* * * * *", epoch, 60)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.SampleSize != 60 {
		t.Errorf("SampleSize = %d, want 60", s.SampleSize)
	}
	if s.MinInterval != time.Minute {
		t.Errorf("MinInterval = %v, want 1m", s.MinInterval)
	}
	if s.MaxInterval != time.Minute {
		t.Errorf("MaxInterval = %v, want 1m", s.MaxInterval)
	}
	if s.FiresPerHour < 59 || s.FiresPerHour > 61 {
		t.Errorf("FiresPerHour = %.2f, want ~60", s.FiresPerHour)
	}
}

func TestCompute_Hourly(t *testing.T) {
	s, err := statistics.Compute("0 * * * *", epoch, 24)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.MinInterval != time.Hour {
		t.Errorf("MinInterval = %v, want 1h", s.MinInterval)
	}
	if s.FiresPerDay < 23 || s.FiresPerDay > 25 {
		t.Errorf("FiresPerDay = %.2f, want ~24", s.FiresPerDay)
	}
}

func TestCompute_InvalidExpression(t *testing.T) {
	_, err := statistics.Compute("bad expr", epoch, 10)
	if err == nil {
		t.Error("expected error for invalid expression, got nil")
	}
}

func TestCompute_SampleSizeTooSmall(t *testing.T) {
	_, err := statistics.Compute("* * * * *", epoch, 1)
	if err == nil {
		t.Error("expected error for sampleSize < 2")
	}
}

func TestCompute_ExpressionStored(t *testing.T) {
	expr := "*/5 * * * *"
	s, err := statistics.Compute(expr, epoch, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Expression != expr {
		t.Errorf("Expression = %q, want %q", s.Expression, expr)
	}
}
