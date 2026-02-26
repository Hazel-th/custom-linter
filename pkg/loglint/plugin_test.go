package loglint

import (
	"testing"

	"github.com/victornechaev/loglint/internal/config"
)

func TestMergeConfigWithSettings(t *testing.T) {
	t.Parallel()

	autoFix := false
	options := mergeConfigWithSettings(
		config.Config{
			SensitivePatterns: []string{"token"},
			CustomPatterns: map[string]string{
				"email": `.+@.+`,
			},
			AutoFix:       true,
			DisabledRules: []string{"english"},
		},
		Settings{
			SensitivePatterns: []string{"refresh token"},
			CustomPatterns: map[string]string{
				"order-id": `\b\d{4}\b`,
			},
			DisabledRules: []string{"lowercase"},
			AutoFix:       &autoFix,
		},
	)

	if !options.DisableFixes {
		t.Fatalf("expected DisableFixes = true")
	}

	if len(options.SensitivePatterns) != 2 {
		t.Fatalf("unexpected sensitive patterns: %#v", options.SensitivePatterns)
	}

	if got := options.CustomPatterns["email"]; got != `.+@.+` {
		t.Fatalf("unexpected email pattern: %q", got)
	}

	if got := options.CustomPatterns["order-id"]; got != `\b\d{4}\b` {
		t.Fatalf("unexpected order-id pattern: %q", got)
	}

	if len(options.DisabledRules) != 2 {
		t.Fatalf("unexpected disabled rules: %#v", options.DisabledRules)
	}
}
