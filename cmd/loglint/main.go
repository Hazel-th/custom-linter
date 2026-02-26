package main

import (
	"log"

	"github.com/victornechaev/loglint/internal/config"
	"github.com/victornechaev/loglint/pkg/loglint"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	cfg, err := config.Load(config.DefaultPath)
	if err != nil {
		log.Fatal(err)
	}

	options := loglint.Options{
		SensitivePatterns: cfg.SensitivePatterns,
		CustomPatterns:    cfg.CustomPatterns,
		DisabledRules:     cfg.DisabledRules,
		DisableFixes:      !cfg.AutoFix,
	}

	singlechecker.Main(loglint.NewAnalyzer(options))
}
