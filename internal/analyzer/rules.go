package analyzer

import (
	"go/ast"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/tools/go/analysis"
)

const (
	ruleLowercase    = "lowercase"
	ruleEnglish      = "english"
	ruleSpecialChars = "specialchars"
	ruleSensitive    = "sensitive"
)

type ruleSpec struct {
	name     string
	message  string
	failed   func(ast.Expr, messageData) bool
	buildFix func(ast.Expr, messageData) (analysis.SuggestedFix, bool)
}

func (r *runner) checkMessage(pass *analysis.Pass, msgExpr ast.Expr) {
	data := collectMessageData(pass, msgExpr)
	textRules := []ruleSpec{
		{
			name:    ruleLowercase,
			message: "log message should start with a lowercase letter",
			failed: func(_ ast.Expr, d messageData) bool {
				return d.hasFullText && startsWithUpperASCII(d.fullText)
			},
			buildFix: func(expr ast.Expr, d messageData) (analysis.SuggestedFix, bool) {
				return buildLowercaseFix(expr, d.fullText)
			},
		},
		{
			name:    ruleEnglish,
			message: "log message should contain only English language",
			failed: func(_ ast.Expr, d messageData) bool {
				return d.hasFullText && containsNonEnglishLetters(d.fullText)
			},
			buildFix: func(expr ast.Expr, d messageData) (analysis.SuggestedFix, bool) {
				return buildEnglishOnlyFix(expr, d.fullText)
			},
		},
		{
			name:    ruleSpecialChars,
			message: "log message must not contain special symbols or emoji",
			failed: func(_ ast.Expr, d messageData) bool {
				return d.hasFullText && containsSpecialSymbolsOrEmoji(d.fullText)
			},
			buildFix: func(expr ast.Expr, d messageData) (analysis.SuggestedFix, bool) {
				return buildSpecialSymbolsFix(expr, d.fullText)
			},
		},
		{
			name:    ruleSensitive,
			message: "log message may contain sensitive data",
			failed: func(expr ast.Expr, d messageData) bool {
				return r.containsSensitiveData(expr, d)
			},
			buildFix: func(expr ast.Expr, _ messageData) (analysis.SuggestedFix, bool) {
				return buildSensitiveDataFix(expr)
			},
		},
	}

	for _, spec := range textRules {
		r.reportRuleViolation(pass, msgExpr, data, spec)
	}
}

func (r *runner) reportRuleViolation(pass *analysis.Pass, expr ast.Expr, data messageData, spec ruleSpec) {
	if !r.ruleEnabled(spec.name) || !spec.failed(expr, data) {
		return
	}

	diag := analysis.Diagnostic{
		Pos:     expr.Pos(),
		End:     expr.End(),
		Message: spec.message,
	}

	if !r.disableFixes && spec.buildFix != nil {
		if fix, ok := spec.buildFix(expr, data); ok {
			diag.SuggestedFixes = []analysis.SuggestedFix{fix}
		}
	}

	pass.Report(diag)
}

func startsWithUpperASCII(text string) bool {
	if text == "" {
		return false
	}

	r, _ := utf8.DecodeRuneInString(text)
	return r >= 'A' && r <= 'Z'
}

func buildLowercaseFix(expr ast.Expr, original string) (analysis.SuggestedFix, bool) {
	fixed, changed := lowercaseFirstASCII(original)
	if !changed {
		return analysis.SuggestedFix{}, false
	}

	return buildReplaceMessageExprFix(expr, fixed, "convert first letter to lowercase")
}

func lowercaseFirstASCII(text string) (string, bool) {
	if text == "" {
		return text, false
	}

	r, size := utf8.DecodeRuneInString(text)
	if r < 'A' || r > 'Z' {
		return text, false
	}

	return string(unicode.ToLower(r)) + text[size:], true
}

func containsNonEnglishLetters(text string) bool {
	for _, r := range text {
		if unicode.IsLetter(r) && !(r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z') {
			return true
		}
	}

	return false
}

func containsSpecialSymbolsOrEmoji(text string) bool {
	for _, r := range text {
		switch {
		case r == ' ':
			continue
		case unicode.IsLetter(r):
			continue
		case unicode.IsDigit(r):
			continue
		default:
			return true
		}
	}

	return false
}

func buildEnglishOnlyFix(expr ast.Expr, original string) (analysis.SuggestedFix, bool) {
	fixed := strings.TrimSpace(filterEnglishLettersOnly(original))
	if fixed == "" {
		fixed = "message"
	}

	if fixed == original {
		return analysis.SuggestedFix{}, false
	}

	return buildReplaceMessageExprFix(expr, fixed, "remove non-English letters")
}

func filterEnglishLettersOnly(text string) string {
	var b strings.Builder
	b.Grow(len(text))
	for _, r := range text {
		switch {
		case r >= 'a' && r <= 'z':
			b.WriteRune(r)
		case r >= 'A' && r <= 'Z':
			b.WriteRune(r)
		case unicode.IsLetter(r):
			// Отбрасываем неанглийские буквы.
			continue
		case unicode.IsDigit(r):
			b.WriteRune(r)
		case unicode.IsSpace(r):
			b.WriteRune(' ')
		default:
			// Пунктуацию и символы здесь сохраняем:
			// их обрабатывает отдельное правило и фиксер для спецсимволов.
			b.WriteRune(r)
		}
	}

	return strings.Join(strings.Fields(b.String()), " ")
}

func buildSpecialSymbolsFix(expr ast.Expr, original string) (analysis.SuggestedFix, bool) {
	fixed := strings.TrimSpace(filterToAlphaNumSpace(original))
	if fixed == "" {
		fixed = "message"
	}

	if fixed == original {
		return analysis.SuggestedFix{}, false
	}

	return buildReplaceMessageExprFix(expr, fixed, "remove special symbols and emoji")
}

func filterToAlphaNumSpace(text string) string {
	var b strings.Builder
	b.Grow(len(text))
	for _, r := range text {
		switch {
		case unicode.IsLetter(r):
			b.WriteRune(r)
		case unicode.IsDigit(r):
			b.WriteRune(r)
		case unicode.IsSpace(r):
			b.WriteRune(' ')
		}
	}

	return strings.Join(strings.Fields(b.String()), " ")
}

func buildSensitiveDataFix(expr ast.Expr) (analysis.SuggestedFix, bool) {
	return buildReplaceMessageExprFix(expr, "sensitive data redacted", "replace with neutral message")
}

func buildReplaceMessageExprFix(expr ast.Expr, fixed, message string) (analysis.SuggestedFix, bool) {
	if fixed == "" {
		return analysis.SuggestedFix{}, false
	}

	return analysis.SuggestedFix{
		Message: message,
		TextEdits: []analysis.TextEdit{
			{
				Pos:     expr.Pos(),
				End:     expr.End(),
				NewText: []byte(strconv.Quote(fixed)),
			},
		},
	}, true
}
