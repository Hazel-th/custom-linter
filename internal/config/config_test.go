package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDefaultsWhenFileMissing(t *testing.T) {
	t.Parallel()

	missingPath := filepath.Join(t.TempDir(), "missing.json")
	cfg, err := Load(missingPath)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if !cfg.AutoFix {
		t.Fatalf("expected AutoFix to be true by default")
	}

	if cfg.CustomPatterns == nil {
		t.Fatalf("expected CustomPatterns map to be initialized")
	}
}

func TestLoadFromFile(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	cfgPath := filepath.Join(dir, ".loglint.json")
	content := `{
		"sensitive_patterns": ["refresh token"],
		"custom_patterns": {"order": "\\\\b\\\\d{4}\\\\b"},
		"auto_fix": false,
		"disabled_rules": ["lowercase"]
	}`
	if err := os.WriteFile(cfgPath, []byte(content), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	cfg, err := Load(cfgPath)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.AutoFix {
		t.Fatalf("expected AutoFix to be false")
	}

	if len(cfg.SensitivePatterns) != 1 || cfg.SensitivePatterns[0] != "refresh token" {
		t.Fatalf("unexpected SensitivePatterns: %#v", cfg.SensitivePatterns)
	}

	if got := cfg.CustomPatterns["order"]; got != `\\b\\d{4}\\b` {
		t.Fatalf("unexpected custom pattern: %q", got)
	}

	if len(cfg.DisabledRules) != 1 || cfg.DisabledRules[0] != "lowercase" {
		t.Fatalf("unexpected DisabledRules: %#v", cfg.DisabledRules)
	}
}
