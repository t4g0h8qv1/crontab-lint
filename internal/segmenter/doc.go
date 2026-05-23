// Package segmenter provides positional decomposition of cron expressions.
//
// It splits a five-field cron expression into typed Segment values that
// record each field's name, raw text, and byte offsets within the original
// string.  This is particularly useful for language-server-style tooling
// where a cursor position must be mapped back to the corresponding cron
// field for hover documentation or inline diagnostics.
//
// Example:
//
//	segs, err := segmenter.Segment("*/15 0 * * 1-5")
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, s := range segs {
//		fmt.Printf("%s: %s\n", s.Name, s.Value)
//	}
//
// AtOffset can be used to resolve a cursor byte position to the relevant
// field:
//
//	field := segmenter.AtOffset(segs, cursorPos)
//	if field != nil {
//		fmt.Println("cursor is in field:", field.Name)
//	}
package segmenter
