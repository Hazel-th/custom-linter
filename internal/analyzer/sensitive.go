package analyzer

import (
	"go/ast"
	"regexp"
	"strings"
)

func (r *runner) containsSensitiveData(expr ast.Expr, data messageData) bool {
	if data.hasFullText && containsCustomPattern(data.fullText, r.customPatterns) {
		return true
	}

	if data.hasFullText && containsSensitiveAssignment(data.fullText, r.sensitivePatterns) {
		return true
	}

	if !data.hasDynamic {
		return false
	}

	literalContext := strings.Join(data.literalParts, " ")
	if containsPattern(literalContext, r.sensitivePatterns) {
		return true
	}

	if containsCustomPattern(literalContext, r.customPatterns) {
		return true
	}

	return exprContainsSensitiveIdentifier(expr, r.sensitivePatterns)
}

func exprContainsSensitiveIdentifier(expr ast.Expr, patterns []string) bool {
	found := false
	ast.Inspect(expr, func(node ast.Node) bool {
		if found {
			return false
		}

		ident, ok := node.(*ast.Ident)
		if !ok {
			return true
		}

		if containsPattern(ident.Name, patterns) {
			found = true
			return false
		}

		return true
	})

	return found
}

func containsSensitiveAssignment(text string, patterns []string) bool {
	normalized := normalizeForSearch(text)
	for _, pattern := range patterns {
		idx := strings.Index(normalized, pattern)
		if idx == -1 {
			continue
		}

		after := strings.TrimSpace(normalized[idx+len(pattern):])
		if after == "" {
			continue
		}

		if strings.HasPrefix(after, ":") || strings.HasPrefix(after, "=") {
			return true
		}
	}

	return false
}

func containsPattern(text string, patterns []string) bool {
	normalized := normalizeForSearch(text)
	for _, pattern := range patterns {
		if strings.Contains(normalized, pattern) {
			return true
		}
	}

	return false
}

func normalizePatterns(patterns []string) []string {
	result := make([]string, 0, len(patterns))
	for _, p := range patterns {
		normalized := normalizeForSearch(p)
		if normalized == "" {
			continue
		}
		result = append(result, normalized)
	}

	return result
}

func containsCustomPattern(text string, patterns []*regexp.Regexp) bool {
	for _, pattern := range patterns {
		if pattern.MatchString(text) {
			return true
		}
	}

	return false
}

func normalizeForSearch(text string) string {
	lower := strings.ToLower(text)
	replacer := strings.NewReplacer("_", " ", "-", " ")
	return replacer.Replace(lower)
}

func defaultSensitivePatterns() []string {
	return []string{
		"password",
		"passwd",
		"pwd",
		"token",
		"secret",
		"api key",
		"apikey",
		"credential",
		"private key",
		"access key",
	}
}
