// Package deduplicator provides utilities for detecting and removing
// duplicate values within crontab expression fields.
//
// A crontab field such as "0,30,0" contains a redundant second 0; this
// package normalises such lists to their shortest equivalent form and
// reports every change made so callers can surface actionable feedback
// to users.
//
// Usage:
//
//	result, err := deduplicator.Dedup("0,30,0 * * * *")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(deduplicator.FormatText(result))
package deduplicator
