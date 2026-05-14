package main

import (
	"testing"

	"github.com/yourorg/crontab-lint/internal/formatter"
	"github.com/yourorg/crontab-lint/internal/linter"
)

func TestBuildResult_ValidExpression(t *testing.T) {
	result := buildResult("*/5 * * * *", nil, nil, "")

	if result.Expression != "*/5 * * * *" {
		t.Errorf("expected expression %q, got %q", "*/5 * * * *", result.Expression)
	}
	if !result.Valid {
		t.Error("expected result to be valid")
	}
	if len(result.Errors) != 0 {
		t.Errorf("expected no errors, got %v", result.Errors)
	}
	if len(result.Issues) != 0 {
		t.Errorf("expected no issues, got %v", result.Issues)
	}
}

func TestBuildResult_WithErrors(t *testing.T) {
	errors := []string{"invalid minute field"}
	result := buildResult("invalid", errors, nil, "")

	if result.Valid {
		t.Error("expected result to be invalid")
	}
	if len(result.Errors) != 1 {
		t.Errorf("expected 1 error, got %d", len(result.Errors))
	}
}

func TestBuildResult_WithIssues(t *testing.T) {
	issues := []linter.Issue{
		{Field: "minute", Message: "redundant step value", Severity: "warning"},
	}
	result := buildResult("*/1 * * * *", nil, issues, "")

	if !result.Valid {
		t.Error("expected result to be valid despite lint issues")
	}
	if len(result.Issues) != 1 {
		t.Errorf("expected 1 issue, got %d", len(result.Issues))
	}
}

func TestBuildResult_WithExplanation(t *testing.T) {
	explanation := "Every 5 minutes"
	result := buildResult("*/5 * * * *", nil, nil, explanation)

	if result.Explanation != explanation {
		t.Errorf("expected explanation %q, got %q", explanation, result.Explanation)
	}
}

func TestBuildResult_ReturnsFormatterResult(t *testing.T) {
	result := buildResult("0 9 * * 1", nil, nil, "At 09:00 on Monday")

	var _ formatter.Result = result
}
