// Package statistics analyses crontab expressions to produce frequency
// and interval statistics.
//
// It uses the scheduler package to generate a sample of future occurrences
// and derives metrics such as fires per hour/day/week and the minimum and
// maximum gap between consecutive firings.
//
// # Overview
//
// The primary entry point is [Compute], which accepts a cron expression, a
// reference start time, and the number of sample occurrences to evaluate.
// A larger sample size yields more accurate interval statistics at the cost
// of additional computation.
//
// # Output formats
//
// Results can be rendered in two ways:
//
//   - [FormatText] returns a human-readable, plain-text summary.
//   - [FormatJSON] returns a JSON-encoded representation suitable for
//     programmatic consumption.
//
// # Basic usage
//
//	s, err := statistics.Compute("*/15 * * * *", time.Now(), 100)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Print(statistics.FormatText(s))
package statistics
