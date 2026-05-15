package comparator

import (
	"encoding/json"
	"fmt"
)

// FormatText returns a human-readable comparison summary.
func FormatText(r Result) string {
	if r.Relation == "equal" {
		return fmt.Sprintf(
			"Both expressions fire at the same frequency: %d times/day.\n"+
				"  A: %s\n"+
				"  B: %s\n",
			r.FrequencyA, r.ExpressionA, r.ExpressionB,
		)
	}

	faster, slower := r.ExpressionB, r.ExpressionA
	fasterCount, slowerCount := r.FrequencyB, r.FrequencyA
	if r.Relation == "slower" {
		faster, slower = r.ExpressionA, r.ExpressionB
		fasterCount, slowerCount = r.FrequencyA, r.FrequencyB
	}

	return fmt.Sprintf(
		"Expression B is %s than A by %.1f%%.\n"+
			"  A: %s — %d times/day\n"+
			"  B: %s — %d times/day\n"+
			"  Delta: %+d fires/day\n"+
			"  Faster: %s (%d/day)  Slower: %s (%d/day)\n",
		r.Relation, r.RatioPercent,
		r.ExpressionA, r.FrequencyA,
		r.ExpressionB, r.FrequencyB,
		r.Delta,
		faster, fasterCount, slower, slowerCount,
	)
}

// jsonResult is the JSON-serialisable form of Result.
type jsonResult struct {
	ExpressionA  string  `json:"expression_a"`
	ExpressionB  string  `json:"expression_b"`
	FrequencyA   int     `json:"frequency_a_per_day"`
	FrequencyB   int     `json:"frequency_b_per_day"`
	Delta        int     `json:"delta"`
	Relation     string  `json:"relation"`
	RatioPercent float64 `json:"ratio_percent"`
}

// FormatJSON returns a JSON-encoded comparison result.
func FormatJSON(r Result) (string, error) {
	data := jsonResult{
		ExpressionA:  r.ExpressionA,
		ExpressionB:  r.ExpressionB,
		FrequencyA:   r.FrequencyA,
		FrequencyB:   r.FrequencyB,
		Delta:        r.Delta,
		Relation:     r.Relation,
		RatioPercent: r.RatioPercent,
	}
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}
