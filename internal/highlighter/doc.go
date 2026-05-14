// Package highlighter provides ANSI color-based syntax highlighting for
// crontab expressions. Each of the five standard cron fields (minute, hour,
// day-of-month, month, day-of-week) is assigned a distinct terminal color,
// making expressions easier to parse visually in CLI output.
//
// Usage:
//
//	result, err := highlighter.Highlight("*/15 9-17 * * 1-5")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(result.Highlighted) // colored expression
//
// The Strip function removes all ANSI escape codes, useful when writing
// output to files or non-terminal destinations.
package highlighter
