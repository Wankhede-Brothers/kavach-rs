// Package patterns provides dynamic pattern loading from TOON config.
// defaults.go: Minimal fallback defaults.
// DACE: Defaults only used when TOON config missing.
package patterns

func loadDefaults(cfg *Config) {
	// Minimal defaults - full config should be in TOON file
	cfg.Sensitive = []string{".env", "credentials", "secret", "password", "private_key"}
	cfg.Blocked = getBlockedDefaults()
	cfg.CodeExts = []string{".go", ".rs", ".ts", ".py", ".js"}
	cfg.LargeExts = []string{".log", ".csv", ".sql"}

	cfg.ValidAgents = map[string][]string{
		"L-1": {"nlu-intent-analyzer"},
		"L0":  {"ceo", "research-director"},
		"L1":  {"backend-engineer", "frontend-engineer"},
		"L2":  {"aegis-guardian"},
	}

	// BUG-001 FIX: Consolidated with gates/intent_nlu.go classifications
	cfg.IntentWords = map[string][]string{
		"debug":     {"fix", "bug", "error", "broken", "crash", "failing", "not working", "doesnt work"},
		"optimize":  {"optimize", "faster", "slow", "performance", "speed up", "efficient"},
		"implement": {"implement", "create", "build", "add", "develop", "write", "new feature"},
		"refactor":  {"refactor", "restructure", "clean up", "improve code", "technical debt"},
		"research":  {"research", "find", "search", "explore", "explain", "how does", "what is"},
		"docs":      {"document", "documentation", "readme", "api docs", "jsdoc", "rustdoc"},
		"audit":     {"audit", "review", "security scan", "vulnerability", "compliance"},
		"status":    {"status", "progress", "state"},
	}
}

// getBlockedDefaults returns comprehensive blocked command patterns.
// Security: Blocks destructive, RCE, DoS, and privilege escalation commands.
func getBlockedDefaults() []string {
	return []string{
		// Destructive file operations
		"rm -rf /",
		"rm -rf /*",
		"rm -rf ~",
		"> /etc/passwd",
		"> /etc/shadow",
		// Disk destruction
		"dd if=/dev/zero",
		"dd if=/dev/random",
		"mkfs.",
		"fdisk",
		"parted",
		// System control
		"shutdown",
		"reboot",
		"init 0",
		"init 6",
		"poweroff",
		"halt",
		// Remote code execution (pipe to shell)
		"| bash",
		"| sh",
		"|bash",
		"|sh",
		"| /bin/bash",
		"| /bin/sh",
		// Fork bomb patterns
		":()",
		":(){",
		// Dangerous permission changes
		"chmod -r 777 /",
		"chmod 777 /",
		"chown -r",
		// Network attacks
		"nc -e",
		"ncat -e",
		// History manipulation
		"history -c",
		"export histsize=0",
		// Kernel manipulation
		"insmod",
		"rmmod",
		"modprobe -r",
	}
}
