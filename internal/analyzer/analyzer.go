package analyzer

import (
	"go/ast"
	"regexp"
	"strings"

	"golang.org/x/tools/go/analysis"
)

type Options struct {
	SensitivePatterns []string
	CustomPatterns    map[string]string
	DisabledRules     []string
	DisableFixes      bool
}

type runner struct {
	sensitivePatterns []string
	customPatterns    []*regexp.Regexp
	disabledRules     map[string]struct{}
	disableFixes      bool
}

func New(options Options) *analysis.Analyzer {
	patterns := defaultSensitivePatterns()
	for _, p := range options.SensitivePatterns {
		trimmed := strings.TrimSpace(p)
		if trimmed == "" {
			continue
		}
		patterns = append(patterns, trimmed)
	}

	r := &runner{
		sensitivePatterns: normalizePatterns(patterns),
		customPatterns:    compilePatterns(options.CustomPatterns),
		disabledRules:     normalizeDisabledRules(options.DisabledRules),
		disableFixes:      options.DisableFixes,
	}

	return &analysis.Analyzer{
		Name: "loglint",
		Doc:  "checks slog and zap log messages for style and security issues",
		Run:  r.run,
	}
}

func compilePatterns(patterns map[string]string) []*regexp.Regexp {
	if len(patterns) == 0 {
		return nil
	}

	result := make([]*regexp.Regexp, 0, len(patterns))
	for _, pattern := range patterns {
		re, err := regexp.Compile(pattern)
		if err != nil {
			continue
		}
		result = append(result, re)
	}

	return result
}

func normalizeDisabledRules(rules []string) map[string]struct{} {
	if len(rules) == 0 {
		return nil
	}

	result := make(map[string]struct{}, len(rules))
	for _, rule := range rules {
		key := strings.TrimSpace(strings.ToLower(rule))
		if key == "" {
			continue
		}
		result[key] = struct{}{}
	}

	return result
}

func (r *runner) ruleEnabled(rule string) bool {
	if len(r.disabledRules) == 0 {
		return true
	}

	_, disabled := r.disabledRules[rule]
	return !disabled
}

func (r *runner) run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			call, ok := node.(*ast.CallExpr)
			if !ok {
				return true
			}

			msgExpr, ok := extractMessageExpr(pass, call)
			if !ok {
				return true
			}

			if !isStringExpr(pass, msgExpr) {
				return true
			}

			r.checkMessage(pass, msgExpr)
			return true
		})
	}

	return nil, nil
}
