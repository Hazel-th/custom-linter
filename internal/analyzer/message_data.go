package analyzer

import (
	"go/ast"
	"go/constant"
	"go/token"
	"strconv"
	"strings"

	"golang.org/x/tools/go/analysis"
)

type messageData struct {
	fullText     string
	hasFullText  bool
	literalParts []string
	hasDynamic   bool
}

func collectMessageData(pass *analysis.Pass, expr ast.Expr) messageData {
	// Сначала пробуем вычислить значение как константу времени компиляции.
	// Это покрывает обычные литералы и полностью вычислимые выражения.
	tv, ok := pass.TypesInfo.Types[expr]
	if ok && tv.Value != nil && tv.Value.Kind() == constant.String {
		value := constant.StringVal(tv.Value)
		return messageData{
			fullText:     value,
			hasFullText:  true,
			literalParts: []string{value},
			hasDynamic:   false,
		}
	}

	// Если полностью вычислить выражение нельзя, сохраняем литеральные части.
	// Они используются как контекст для sensitive-проверок динамических выражений.
	parts, dynamic := collectLiteralParts(expr)
	if !dynamic && len(parts) > 0 {
		joined := strings.Join(parts, "")
		return messageData{
			fullText:     joined,
			hasFullText:  true,
			literalParts: parts,
			hasDynamic:   false,
		}
	}

	return messageData{
		literalParts: parts,
		hasDynamic:   dynamic,
	}
}

func collectLiteralParts(expr ast.Expr) ([]string, bool) {
	switch e := expr.(type) {
	case *ast.BasicLit:
		if e.Kind != token.STRING {
			return nil, true
		}

		value, err := strconv.Unquote(e.Value)
		if err != nil {
			return nil, true
		}

		return []string{value}, false
	case *ast.BinaryExpr:
		if e.Op != token.ADD {
			return nil, true
		}

		leftParts, leftDynamic := collectLiteralParts(e.X)
		rightParts, rightDynamic := collectLiteralParts(e.Y)
		return append(leftParts, rightParts...), leftDynamic || rightDynamic
	case *ast.ParenExpr:
		return collectLiteralParts(e.X)
	default:
		return nil, true
	}
}
