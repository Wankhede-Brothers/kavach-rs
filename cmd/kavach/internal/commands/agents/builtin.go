// Package agents - builtin.go
// DACE: Single responsibility - built-in agent definitions
package agents

// GetBuiltin returns the hardcoded agent hierarchy
func GetBuiltin() []*Agent {
	return []*Agent{
		{
			Name:        "nlu-intent-analyzer",
			Level:       LevelNLU,
			Model:       "haiku",
			Description: "Fast intent parsing - routes to CEO",
			Triggers:    []string{"user_prompt"},
			Tools:       []string{"Read", "Grep", "Glob"},
		},
		{
			Name:        "ceo",
			Level:       LevelCEO,
			Model:       "opus",
			Description: "Orchestrator - delegates, never writes code",
			Triggers:    []string{"task_delegation", "complex_request"},
			Tools:       []string{"Task", "Read", "Grep", "Glob", "WebSearch"},
		},
		{
			Name:        "research-director",
			Level:       LevelCEO,
			Model:       "opus",
			Description: "Evidence-based research findings",
			Triggers:    []string{"research_needed", "verify_facts"},
			Tools:       []string{"WebSearch", "WebFetch", "Read", "Grep"},
		},
		{
			Name:        "backend-engineer",
			Level:       LevelEngineer,
			Model:       "sonnet",
			Description: "Rust backend - Axum, Tonic, Zig",
			Triggers:    []string{"backend_task", "api_implementation"},
			Tools:       []string{"Read", "Edit", "Write", "Bash", "Grep"},
		},
		{
			Name:        "frontend-engineer",
			Level:       LevelEngineer,
			Model:       "sonnet",
			Description: "TypeScript + React frontend",
			Triggers:    []string{"frontend_task", "ui_implementation"},
			Tools:       []string{"Read", "Edit", "Write", "Bash", "Grep"},
		},
		{
			Name:        "devops-engineer",
			Level:       LevelEngineer,
			Model:       "sonnet",
			Description: "Docker, K8s, CI/CD pipelines",
			Triggers:    []string{"devops_task", "deployment"},
			Tools:       []string{"Read", "Edit", "Write", "Bash"},
		},
		{
			Name:        "security-engineer",
			Level:       LevelEngineer,
			Model:       "sonnet",
			Description: "Security analysis, OWASP compliance",
			Triggers:    []string{"security_review", "vulnerability_check"},
			Tools:       []string{"Read", "Grep", "WebSearch"},
		},
		{
			Name:        "qa-lead",
			Level:       LevelEngineer,
			Model:       "sonnet",
			Description: "Test strategy and coverage",
			Triggers:    []string{"testing_task", "coverage_analysis"},
			Tools:       []string{"Read", "Edit", "Write", "Bash"},
		},
		{
			Name:        "code-reviewer",
			Level:       LevelReview,
			Model:       "sonnet",
			Description: "Post-implementation code review",
			Triggers:    []string{"code_review", "pr_review"},
			Tools:       []string{"Read", "Grep", "Bash"},
		},
		{
			Name:        "aegis-guardian",
			Level:       LevelAegis,
			Model:       "opus",
			Description: "Verification Guardian - Quality, Security, Testing",
			Triggers:    []string{"final_verification", "quality_gate"},
			Tools:       []string{"Read", "Grep", "Bash"},
		},
	}
}
