package loglint

import (
	internalanalyzer "github.com/victornechaev/loglint/internal/analyzer"
	"golang.org/x/tools/go/analysis"
)

// Options задает поведение правил анализатора loglint.
type Options = internalanalyzer.Options

// Analyzer — анализатор по умолчанию для тестов и простого подключения.
var Analyzer = internalanalyzer.New(Options{})

// NewAnalyzer создает новый анализатор с явно заданными параметрами.
func NewAnalyzer(options Options) *analysis.Analyzer {
	return internalanalyzer.New(options)
}
