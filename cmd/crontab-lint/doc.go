// Package main provides the crontab-lint command-line interface.
//
// crontab-lint is a static analyzer and validator for cron expressions.
// It parses, validates, lints, and explains crontab entries in a
// human-readable format.
//
// Usage:
//
//	crontab-lint [options] "<cron expression>"
//
// Options:
//
//	-json     Output results as JSON instead of plain text
//	-verbose  Include a human-readable explanation of the expression
//
// Exit codes:
//
//	0  Expression is valid with no lint issues
//	1  Invalid usage or missing arguments
//	2  Validation errors or lint issues found
//
// Examples:
//
//	crontab-lint "*/5 * * * *"
//	crontab-lint -verbose "0 9 * * 1-5"
//	crontab-lint -json "30 18 1 * *"
package main
