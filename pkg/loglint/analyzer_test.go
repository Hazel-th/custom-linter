package loglint_test

import (
	"testing"

	"github.com/victornechaev/loglint/pkg/loglint"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	t.Parallel()

	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, loglint.Analyzer, "a", "ok")
}

func TestSuggestedFixes(t *testing.T) {
	t.Parallel()

	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, loglint.Analyzer, "fix")
}

func TestDisabledRule(t *testing.T) {
	t.Parallel()

	testdata := analysistest.TestData()
	analysistest.Run(
		t,
		testdata,
		loglint.NewAnalyzer(loglint.Options{DisabledRules: []string{"lowercase"}}),
		"configdisabled",
	)
}

func TestCustomSensitivePattern(t *testing.T) {
	t.Parallel()

	testdata := analysistest.TestData()
	analysistest.Run(
		t,
		testdata,
		loglint.NewAnalyzer(loglint.Options{
			CustomPatterns: map[string]string{
				"order-id": `\b\d{4}\b`,
			},
		}),
		"custompattern",
	)
}
