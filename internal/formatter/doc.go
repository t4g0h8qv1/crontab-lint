// Package formatter provides output formatting utilities for crontab-lint
// analysis reports.
//
// A Report aggregates the results of parsing, validation, and explanation
// for a single crontab expression. The package exposes two rendering formats:
//
//   - FormatText: plain-text output suitable for terminal display.
//   - FormatJSON: JSON output suitable for machine consumption or editor
//     integrations.
//
// Example usage:
//
//	report := formatter.Report{
//		Expression:  "0 9 * * 1-5",
//		Valid:        true,
//		Explanation:  "At 09:00 on every weekday",
//	}
//	fmt.Print(formatter.FormatText(report))
package formatter
