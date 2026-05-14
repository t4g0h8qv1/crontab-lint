// Package normalizer provides crontab expression normalization utilities.
//
// It supports:
//   - Expanding shorthand aliases such as @daily, @weekly, @monthly, etc.
//   - Replacing named month abbreviations (jan–dec) with numeric values (1–12).
//   - Replacing named day-of-week abbreviations (sun–sat) with numeric values (0–6).
//   - Normalizing extra whitespace in expressions.
//
// Normalized expressions use only numeric values and standard cron syntax,
// making them suitable for further parsing and validation.
//
// Example usage:
//
//	expr := normalizer.Normalize("@daily")
//	// expr == "0 0 * * *"
//
//	expr = normalizer.Normalize("0 0 1 jan *")
//	// expr == "0 0 1 1 *"
package normalizer
