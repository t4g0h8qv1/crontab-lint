// Package linter implements static analysis rules for crontab expressions.
//
// It operates independently of the parser and validator, focusing on
// code-quality style issues, redundancies, and logical conflicts that
// are technically valid but likely unintentional.
//
// # Severity Levels
//
//   - SeverityError   – the expression is malformed or invalid.
//   - SeverityWarning – the expression is valid but may behave unexpectedly.
//   - SeverityInfo    – the expression can be simplified or improved.
//
// # Usage
//
//	issues := linter.Lint("0 12 15 * 5")
//	for _, issue := range issues {
//		fmt.Println(issue)
//	}
package linter
