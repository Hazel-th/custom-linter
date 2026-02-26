package analyzer

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

func extractMessageExpr(pass *analysis.Pass, call *ast.CallExpr) (ast.Expr, bool) {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return nil, false
	}

	msgIndex, ok := resolveMessageIndex(pass, sel)
	if !ok || msgIndex < 0 || len(call.Args) <= msgIndex {
		return nil, false
	}

	return call.Args[msgIndex], true
}

func resolveMessageIndex(pass *analysis.Pass, sel *ast.SelectorExpr) (int, bool) {
	// Вызовы пакетного уровня (например, slog.Info / slog.InfoContext).
	if pkgPath, ok := packagePath(pass, sel.X); ok {
		switch pkgPath {
		case "log/slog":
			switch sel.Sel.Name {
			case "Debug", "Info", "Warn", "Error":
				return 0, true
			case "DebugContext", "InfoContext", "WarnContext", "ErrorContext":
				return 1, true
			}
		}
	}

	// Вызовы методов на инстансах логгеров (например, logger.Info / sugar.Infow).
	named := namedType(pass.TypesInfo.TypeOf(sel.X))
	if named == nil || named.Obj() == nil || named.Obj().Pkg() == nil {
		return -1, false
	}

	pkgPath := named.Obj().Pkg().Path()
	typeName := named.Obj().Name()
	methodName := sel.Sel.Name

	switch {
	case pkgPath == "log/slog" && typeName == "Logger":
		switch methodName {
		case "Debug", "Info", "Warn", "Error":
			return 0, true
		case "DebugContext", "InfoContext", "WarnContext", "ErrorContext":
			return 1, true
		}
	case pkgPath == "go.uber.org/zap" && typeName == "Logger":
		switch methodName {
		case "Debug", "Info", "Warn", "Error", "DPanic", "Panic", "Fatal":
			return 0, true
		}
	case pkgPath == "go.uber.org/zap" && typeName == "SugaredLogger":
		switch methodName {
		case "Debug", "Info", "Warn", "Error", "DPanic", "Panic", "Fatal",
			"Debugw", "Infow", "Warnw", "Errorw", "DPanicw", "Panicw", "Fatalw":
			return 0, true
		}
	}

	return -1, false
}

func packagePath(pass *analysis.Pass, expr ast.Expr) (string, bool) {
	// Для селекторов пакетного уровня путь импорта получаем через PkgName.
	ident, ok := expr.(*ast.Ident)
	if !ok {
		return "", false
	}

	usedObj, ok := pass.TypesInfo.Uses[ident]
	if !ok {
		return "", false
	}

	pkgName, ok := usedObj.(*types.PkgName)
	if !ok || pkgName.Imported() == nil {
		return "", false
	}

	return pkgName.Imported().Path(), true
}

func namedType(typ types.Type) *types.Named {
	if typ == nil {
		return nil
	}

	if ptr, ok := typ.(*types.Pointer); ok {
		typ = ptr.Elem()
	}

	named, _ := typ.(*types.Named)
	return named
}

func isStringExpr(pass *analysis.Pass, expr ast.Expr) bool {
	typ := pass.TypesInfo.TypeOf(expr)
	if typ == nil {
		return false
	}

	return types.AssignableTo(typ, types.Typ[types.String])
}
