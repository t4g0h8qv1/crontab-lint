package normalizer

// AliasInfo holds metadata about a crontab shorthand alias.
type AliasInfo struct {
	// Alias is the shorthand string (e.g. "@daily").
	Alias string
	// Expansion is the canonical five-field expression.
	Expansion string
	// Description is a human-readable explanation of when the job runs.
	Description string
}

// KnownAliases is the list of all recognised crontab shorthand aliases
// along with their canonical expansions and descriptions.
var KnownAliases = []AliasInfo{
	{"@yearly", "0 0 1 1 *", "Run once a year at midnight on January 1st"},
	{"@annually", "0 0 1 1 *", "Run once a year at midnight on January 1st (same as @yearly)"},
	{"@monthly", "0 0 1 * *", "Run once a month at midnight on the first day"},
	{"@weekly", "0 0 * * 0", "Run once a week at midnight on Sunday"},
	{"@daily", "0 0 * * *", "Run once a day at midnight"},
	{"@midnight", "0 0 * * *", "Run once a day at midnight (same as @daily)"},
	{"@hourly", "0 * * * *", "Run once an hour at the beginning of the hour"},
}

// IsAlias reports whether the given expression is a recognised shorthand alias.
func IsAlias(expr string) bool {
	for _, a := range KnownAliases {
		if a.Alias == expr {
			return true
		}
	}
	return false
}

// DescribeAlias returns the human-readable description for a known alias,
// and a boolean indicating whether the alias was found.
func DescribeAlias(expr string) (string, bool) {
	for _, a := range KnownAliases {
		if a.Alias == expr {
			return a.Description, true
		}
	}
	return "", false
}
