// Package expander resolves each field of a crontab expression into the
// complete, sorted set of integer values it matches.
//
// Given an expression such as "*/15 8-10 * * 1,3", Expand returns a
// Result containing one ExpandedField per cron field.  Each ExpandedField
// lists every concrete integer the field will match during scheduling,
// making it straightforward to audit or visualise what a cron entry
// actually covers without having to mentally evaluate ranges and steps.
//
// Example:
//
//	r := expander.Expand("*/15 9 * * *")
//	for _, f := range r.Fields {
//		fmt.Printf("%s: %s\n", f.Name, expander.JoinValues(f.Values))
//	}
//	// minute: 0,15,30,45
//	// hour:   9
//	// ...
package expander
