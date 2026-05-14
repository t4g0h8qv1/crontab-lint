package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yourorg/crontab-lint/internal/formatter"
	"github.com/yourorg/crontab-lint/internal/linter"
	"github.com/yourorg/crontab-lint/internal/validator"
	"github.com/yourorg/crontab-lint/internal/explainer"
)

func main() {
	jsonOutput := flag.Bool("json", false, "Output results as JSON")
	verbose := flag.Bool("verbose", false, "Include human-readable explanation")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: crontab-lint [options] <cron expression>\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExample:\n  crontab-lint \"*/5 * * * *\"\n")
	}
	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	expr := flag.Arg(0)

	validationErrors := validator.Validate(expr)
	lintIssues := linter.Lint(expr)

	var explanation string
	if *verbose {
		var err error
		explanation, err = explainer.Explain(expr)
		if err != nil {
			explanation = fmt.Sprintf("(explanation unavailable: %v)", err)
		}
	}

	result := buildResult(expr, validationErrors, lintIssues, explanation)

	var output string
	if *jsonOutput {
		output = formatter.FormatJSON(result)
	} else {
		output = formatter.FormatText(result)
	}

	fmt.Print(output)

	if len(validationErrors) > 0 || len(lintIssues) > 0 {
		os.Exit(2)
	}
}

func buildResult(expr string, validationErrors []string, lintIssues []linter.Issue, explanation string) formatter.Result {
	return formatter.Result{
		Expression:  expr,
		Valid:       len(validationErrors) == 0,
		Errors:      validationErrors,
		Issues:      lintIssues,
		Explanation: explanation,
	}
}
