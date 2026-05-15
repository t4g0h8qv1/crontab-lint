// Package differ provides field-level comparison of two crontab expressions.
//
// It normalizes both expressions (resolving aliases and named values) before
// comparing them, so semantically equivalent forms such as "@daily" and
// "0 0 * * *" are correctly identified as equal.
//
// Usage:
//
//	result, err := differ.Diff("0 * * * *", "30 * * * *")
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(differ.Summary(result))
package differ
