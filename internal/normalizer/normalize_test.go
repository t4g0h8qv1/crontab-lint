package normalizer

import (
	"testing"
)

func TestNormalize_Aliases(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"@yearly", "0 0 1 1 *"},
		{"@annually", "0 0 1 1 *"},
		{"@monthly", "0 0 1 * *"},
		{"@weekly", "0 0 * * 0"},
		{"@daily", "0 0 * * *"},
		{"@midnight", "0 0 * * *"},
		{"@hourly", "0 * * * *"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := Normalize(tt.input)
			if got != tt.expected {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestNormalize_MonthNames(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"0 0 1 jan *", "0 0 1 1 *"},
		{"0 0 1 dec *", "0 0 1 12 *"},
		{"0 0 1 jan-jun *", "0 0 1 1-6 *"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := Normalize(tt.input)
			if got != tt.expected {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestNormalize_DowNames(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"0 0 * * sun", "0 0 * * 0"},
		{"0 0 * * mon-fri", "0 0 * * 1-5"},
		{"0 0 * * sat", "0 0 * * 6"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := Normalize(tt.input)
			if got != tt.expected {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestNormalize_Passthrough(t *testing.T) {
	input := "*/5 0 * * *"
	got := Normalize(input)
	if got != input {
		t.Errorf("Normalize(%q) = %q, want %q", input, got, input)
	}
}

func TestNormalize_InvalidFieldCount(t *testing.T) {
	input := "* * *"
	got := Normalize(input)
	if got != input {
		t.Errorf("Normalize(%q) = %q, want passthrough %q", input, got, input)
	}
}
