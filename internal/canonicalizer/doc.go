// Package canonicalizer rewrites crontab expressions into a stable,
// canonical form without changing their semantics.
//
// Transformations applied:
//   - Aliases and named constants (e.g. @hourly, JAN, MON) are expanded
//     via the normalizer before further processing.
//   - Step-1 shorthand (*/1, 0-59/1) is simplified to a bare wildcard or range.
//   - Comma-separated numeric lists are sorted in ascending order.
//
// Example:
//
//	result := canonicalizer.Canonicalize("30,5 */1 * * *")
//	// result.Canonical == "5,30 * * * *"
//	// result.Changes   == ["minute: list sorted to 5,30", "hour: */1 simplified to *"]
package canonicalizer
