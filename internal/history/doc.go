// Package history provides a bounded, in-memory record of crontab expressions
// that have been validated during a session.
//
// Each Entry captures the expression string, whether it was valid, any lint
// issues found, and the time it was recorded. The History type enforces a
// configurable maximum capacity, evicting the oldest entry when full.
//
// Example usage:
//
//	h := history.New(100)
//	h.Add("0 * * * *", true, nil)
//	h.Add("bad expr", false, []string{"invalid field count"})
//	fmt.Print(h.Format())
package history
