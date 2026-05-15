// Package statistics analyses crontab expressions to produce frequency
// and interval statistics.
//
// It uses the scheduler package to generate a sample of future occurrences
// and derives metrics such as fires per hour/day/week and the minimum and
// maximum gap between consecutive firings.
//
// Basic usage:
//
//	s, err := statistics.Compute("*/15 * * * *", time.Now(), 100)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Print(statistics.FormatText(s))
package statistics
