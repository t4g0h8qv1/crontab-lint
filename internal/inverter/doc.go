// Package inverter computes the logical complement of a cron expression.
//
// Given a cron expression, Invert returns a new expression that fires at every
// minute that the original expression does NOT fire. Wildcard fields (*) are
// left as wildcards because their complement within a 5-field cross-product
// would be the empty set.
//
// Example:
//
//	res := inverter.Invert("0 * * * *")  // fires every hour on the hour
//	// res.Inverted => "1,2,...,59 * * * *"  (every minute except :00)
//
// Invert relies on the expander package to enumerate concrete values for each
// field before computing the set complement against the field's full range.
package inverter
