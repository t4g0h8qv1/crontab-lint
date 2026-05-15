// Package classifier categorizes crontab expressions into human-friendly
// schedule types such as "hourly", "daily", "weekly", etc.
package classifier

import "strings"

// Class represents a named schedule category.
type Class string

const (
	ClassEveryMinute Class = "every-minute"
	ClassHourly      Class = "hourly"
	ClassDaily       Class = "daily"
	ClassWeekly      Class = "weekly"
	ClassMonthly     Class = "monthly"
	ClassYearly      Class = "yearly"
	ClassCustom      Class = "custom"
	ClassUnknown     Class = "unknown"
)

// Result holds the classification outcome for a crontab expression.
type Result struct {
	Expression string `json:"expression"`
	Class      Class  `json:"class"`
	Label      string `json:"label"`
}

// Classify returns the schedule class for the given 5-field crontab expression.
// It expects a normalized expression (no aliases).
func Classify(expr string) Result {
	fields := strings.Fields(expr)
	if len(fields) != 5 {
		return Result{Expression: expr, Class: ClassUnknown, Label: "unknown (invalid field count)"}
	}

	min, hour, dom, month, dow := fields[0], fields[1], fields[2], fields[3], fields[4]

	switch {
	case isWild(min) && isWild(hour) && isWild(dom) && isWild(month) && isWild(dow):
		return Result{expr, ClassEveryMinute, "every minute"}

	case !isWild(min) && isWild(hour) && isWild(dom) && isWild(month) && isWild(dow):
		return Result{expr, ClassHourly, "hourly"}

	case !isWild(min) && !isWild(hour) && isWild(dom) && isWild(month) && isWild(dow):
		return Result{expr, ClassDaily, "daily"}

	case !isWild(min) && !isWild(hour) && isWild(dom) && isWild(month) && !isWild(dow):
		return Result{expr, ClassWeekly, "weekly"}

	case !isWild(min) && !isWild(hour) && !isWild(dom) && isWild(month) && isWild(dow):
		return Result{expr, ClassMonthly, "monthly"}

	case !isWild(min) && !isWild(hour) && !isWild(dom) && !isWild(month) && isWild(dow):
		return Result{expr, ClassYearly, "yearly"}

	default:
		return Result{expr, ClassCustom, "custom schedule"}
	}
}

// isWild reports whether a crontab field is a plain wildcard.
func isWild(field string) bool {
	return field == "*"
}
