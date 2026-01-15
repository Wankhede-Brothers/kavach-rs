// Package patterns provides dynamic pattern loading from TOON config.
// loader.go: TOON config loader entry point.
// DACE: Dynamic patterns, never hardcode.
package patterns

import (
	"path/filepath"

	"github.com/claude/shared/pkg/util"
)

var cached *Config

// Load returns patterns from TOON config (cached after first load).
func Load() *Config {
	if cached != nil {
		return cached
	}
	cached = loadFromTOON()
	return cached
}

// Reload forces reload from TOON config.
func Reload() *Config {
	cached = nil
	return Load()
}

func loadFromTOON() *Config {
	cfg := &Config{
		ValidAgents: make(map[string][]string),
		IntentWords: make(map[string][]string),
		LoadedFrom:  "default",
	}

	configPath := filepath.Join(util.ClaudeDir(), "config", "patterns.toon")
	if !util.FileExists(configPath) {
		cfg.LoadedFrom = "defaults"
		loadDefaults(cfg)
		return cfg
	}

	content := util.ReadFileString(configPath)
	if content == "" {
		loadDefaults(cfg)
		return cfg
	}

	cfg.LoadedFrom = configPath
	parseTOON(content, cfg)
	return cfg
}
