package normalizer

import "testing"

func TestIsAlias_Known(t *testing.T) {
	known := []string{
		"@yearly", "@annually", "@monthly",
		"@weekly", "@daily", "@midnight", "@hourly",
	}
	for _, alias := range known {
		if !IsAlias(alias) {
			t.Errorf("IsAlias(%q) = false, want true", alias)
		}
	}
}

func TestIsAlias_Unknown(t *testing.T) {
	unknown := []string{"*/5 * * * *", "@reboot", "", "daily"}
	for _, expr := range unknown {
		if IsAlias(expr) {
			t.Errorf("IsAlias(%q) = true, want false", expr)
		}
	}
}

func TestDescribeAlias_Found(t *testing.T) {
	tests := []struct {
		alias string
		want  string
	}{
		{"@daily", "Run once a day at midnight"},
		{"@hourly", "Run once an hour at the beginning of the hour"},
		{"@weekly", "Run once a week at midnight on Sunday"},
	}
	for _, tt := range tests {
		t.Run(tt.alias, func(t *testing.T) {
			desc, ok := DescribeAlias(tt.alias)
			if !ok {
				t.Fatalf("DescribeAlias(%q): ok = false, want true", tt.alias)
			}
			if desc != tt.want {
				t.Errorf("DescribeAlias(%q) = %q, want %q", tt.alias, desc, tt.want)
			}
		})
	}
}

func TestDescribeAlias_NotFound(t *testing.T) {
	_, ok := DescribeAlias("@reboot")
	if ok {
		t.Error("DescribeAlias(@reboot): ok = true, want false")
	}
}

func TestKnownAliases_AllHaveExpansions(t *testing.T) {
	for _, a := range KnownAliases {
		if a.Expansion == "" {
			t.Errorf("alias %q has empty expansion", a.Alias)
		}
		if a.Description == "" {
			t.Errorf("alias %q has empty description", a.Alias)
		}
	}
}
