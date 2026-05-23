// Package annotator provides field-level annotation for crontab expressions.
//
// Each field of a five-part crontab expression is paired with a concise
// human-readable note (e.g. "every 5 minutes", "at hour 9") derived from
// the explainer package.
//
// Usage:
//
//	r := annotator.Annotate("0 9 * * 1")
//	for _, f := range r.Fields {
//	    fmt.Printf("%s (%s): %s\n", f.Name, f.Value, f.Note)
//	}
//
// The Inline helper returns the original expression with a compact comment
// appended:
//
//	fmt.Println(annotator.Inline("*/5 * * * *"))
//	// */5 * * * *  # every 5 min | any hour | any day | any month | any weekday
package annotator
