package linter_test

import (
	"testing"

	"github.com/user/crontab-lint/internal/linter"
)

func TestLint_Valid(t *testing.T) {
	tests := []struct {
		expr string
	}{
		{"* * * * *"},
		{"0 12 * * 1"},
		{"*/5 * * * *"},
		{"0 0 1 1 *"},
	}

	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			issues := linter.Lint(tt.expr)
			for _, issue := range issues {
				if issue.Severity == linter.SeverityError {
					t.Errorf("unexpected error for %q: %s", tt.expr, issue)
				}
			}
		})
	}
}

func TestLint_FieldCountError(t *testing.T) {
	issues := linter.Lint("* * * *")
	if len(issues) == 0 {
		t.Fatal("expected issues for wrong field count, got none")
	}
	if issues[0].Severity != linter.SeverityError {
		t.Errorf("expected error severity, got %s", issues[0].Severity)
	}
}

func TestLint_RedundantStep(t *testing.T) {
	issues := linter.Lint("*/1 * * * *")
	if len(issues) == 0 {
		t.Fatal("expected info issue for redundant step, got none")
	}
	found := false
	for _, issue := range issues {
		if issue.Severity == linter.SeverityInfo {
			found = true
		}
	}
	if !found {
		t.Error("expected at least one info-level issue")
	}
}

func TestLint_DomDowConflict(t *testing.T) {
	issues := linter.Lint("0 12 15 * 5")
	if len(issues) == 0 {
		t.Fatal("expected warning for dom+dow conflict, got none")
	}
	found := false
	for _, issue := range issues {
		if issue.Severity == linter.SeverityWarning {
			found = true
		}
	}
	if !found {
		t.Error("expected at least one warning-level issue")
	}
}

func TestLint_DoubleWildcard(t *testing.T) {
	issues := linter.Lint("** * * * *")
	if len(issues) == 0 {
		t.Fatal("expected error for double wildcard, got none")
	}
	if issues[0].Severity != linter.SeverityError {
		t.Errorf("expected error severity, got %s", issues[0].Severity)
	}
}
