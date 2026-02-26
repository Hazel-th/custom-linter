package config

import (
	"encoding/json"
	"errors"
	"os"
)

// DefaultPath — путь к конфигу loglint по умолчанию.
const DefaultPath = ".loglint.json"

// Config содержит пользовательские настройки правил.
type Config struct {
	SensitivePatterns []string          `json:"sensitive_patterns"`
	CustomPatterns    map[string]string `json:"custom_patterns"`
	AutoFix           bool              `json:"auto_fix"`
	DisabledRules     []string          `json:"disabled_rules"`
}

// Defaults возвращает базовые настройки, если файл конфигурации не задан.
func Defaults() Config {
	return Config{
		SensitivePatterns: nil,
		CustomPatterns:    map[string]string{},
		AutoFix:           true,
		DisabledRules:     nil,
	}
}

// Load читает конфиг по path; если файла нет, возвращаются настройки по умолчанию.
func Load(path string) (Config, error) {
	if path == "" {
		path = DefaultPath
	}

	cfg := Defaults()
	raw, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return cfg, nil
		}
		return Config{}, err
	}

	if err := json.Unmarshal(raw, &cfg); err != nil {
		return Config{}, err
	}

	if cfg.CustomPatterns == nil {
		cfg.CustomPatterns = map[string]string{}
	}

	return cfg, nil
}
