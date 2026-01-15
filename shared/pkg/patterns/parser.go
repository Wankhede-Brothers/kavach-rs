// Package patterns provides dynamic pattern loading from TOON config.
// parser.go: TOON format parser for patterns.
// DACE: Single responsibility - parsing only.
package patterns

import "strings"

func parseTOON(content string, cfg *Config) {
	var currentBlock string
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			currentBlock = strings.Trim(line, "[]")
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		values := parseList(val)

		applyToConfig(cfg, currentBlock, key, values)
	}
}

func applyToConfig(cfg *Config, block, key string, values []string) {
	switch block {
	case "SENSITIVE":
		cfg.Sensitive = append(cfg.Sensitive, values...)
	case "BLOCKED":
		cfg.Blocked = append(cfg.Blocked, values...)
	case "CODE_EXTS":
		cfg.CodeExts = append(cfg.CodeExts, values...)
	case "LARGE_EXTS":
		cfg.LargeExts = append(cfg.LargeExts, values...)
	case "AGENTS":
		cfg.ValidAgents[key] = values
	case "INTENT":
		cfg.IntentWords[key] = values
	}
}

func parseList(val string) []string {
	val = strings.Trim(val, "[]")
	parts := strings.Split(val, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}
