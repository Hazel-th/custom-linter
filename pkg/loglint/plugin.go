package loglint

import (
	"github.com/golangci/plugin-module-register/register"
	"github.com/victornechaev/loglint/internal/config"
	"golang.org/x/tools/go/analysis"
)

func init() {
	register.Plugin("loglint", New)
}

// Settings описывает конфигурацию loglint из YAML-настроек golangci-lint.
type Settings struct {
	SensitivePatterns []string          `json:"sensitive-patterns"`
	CustomPatterns    map[string]string `json:"custom-patterns"`
	DisabledRules     []string          `json:"disabled-rules"`
	AutoFix           *bool             `json:"auto-fix"`
	ConfigPath        string            `json:"config-path"`
}

// Plugin — адаптер module-plugin, который ожидает golangci-lint.
type Plugin struct {
	settings Settings
}

// New декодирует настройки плагина и создает его экземпляр.
func New(settings any) (register.LinterPlugin, error) {
	cfg, err := register.DecodeSettings[Settings](settings)
	if err != nil {
		return nil, err
	}

	return &Plugin{settings: cfg}, nil
}

func (p *Plugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	fileCfg, err := config.Load(p.settings.ConfigPath)
	if err != nil {
		return nil, err
	}

	options := mergeConfigWithSettings(fileCfg, p.settings)
	analyzer := NewAnalyzer(options)
	return []*analysis.Analyzer{analyzer}, nil
}

func (p *Plugin) GetLoadMode() string {
	return register.LoadModeTypesInfo
}

func mergeConfigWithSettings(cfg config.Config, settings Settings) Options {
	// Явно заданная YAML-опция должна переопределять значение из файла конфигурации.
	autoFix := cfg.AutoFix
	if settings.AutoFix != nil {
		autoFix = *settings.AutoFix
	}

	return Options{
		SensitivePatterns: mergeStringSlices(cfg.SensitivePatterns, settings.SensitivePatterns),
		CustomPatterns:    mergeStringMaps(cfg.CustomPatterns, settings.CustomPatterns),
		DisabledRules:     mergeStringSlices(cfg.DisabledRules, settings.DisabledRules),
		DisableFixes:      !autoFix,
	}
}

func mergeStringMaps(base, override map[string]string) map[string]string {
	merged := make(map[string]string, len(base)+len(override))
	for key, value := range base {
		merged[key] = value
	}
	for key, value := range override {
		merged[key] = value
	}

	return merged
}

func mergeStringSlices(base, override []string) []string {
	merged := make([]string, 0, len(base)+len(override))
	merged = append(merged, base...)
	merged = append(merged, override...)
	return merged
}
