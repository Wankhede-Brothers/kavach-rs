// Package skills - builtin.go
// DACE: Single responsibility - built-in skill definitions
package skills

// GetBuiltin returns the hardcoded skills
func GetBuiltin() []*Skill {
	return []*Skill{
		// Git skills
		{
			Name:        "commit",
			Category:    CatGit,
			Description: "Create git commit with conventional format",
			Triggers:    []string{"/commit", "commit changes"},
			Commands:    []string{"git add", "git commit"},
		},
		{
			Name:        "review-pr",
			Category:    CatGit,
			Description: "Review pull request changes",
			Triggers:    []string{"/review-pr", "review pr"},
			Commands:    []string{"gh pr view", "gh pr diff"},
		},
		{
			Name:        "create-pr",
			Category:    CatGit,
			Description: "Create pull request",
			Triggers:    []string{"/create-pr", "create pr"},
			Commands:    []string{"gh pr create"},
		},

		// Session skills
		{
			Name:        "init",
			Category:    CatSession,
			Description: "Initialize session with context",
			Triggers:    []string{"/init", "start session"},
			Commands:    []string{"kavach session init"},
		},
		{
			Name:        "status",
			Category:    CatSession,
			Description: "Show system status",
			Triggers:    []string{"/status", "show status"},
			Commands:    []string{"kavach status"},
		},

		// Memory skills
		{
			Name:        "memory",
			Category:    CatMemory,
			Description: "Query memory bank",
			Triggers:    []string{"/memory", "memory bank"},
			Commands:    []string{"kavach memory bank"},
		},
		{
			Name:        "kanban",
			Category:    CatMemory,
			Description: "Show kanban dashboard",
			Triggers:    []string{"/kanban", "show kanban"},
			Commands:    []string{"kavach memory kanban"},
		},

		// Research skills
		{
			Name:        "plan",
			Category:    CatResearch,
			Description: "Plan implementation approach",
			Triggers:    []string{"/plan", "create plan"},
			Commands:    []string{"EnterPlanMode"},
		},
		{
			Name:        "explore",
			Category:    CatResearch,
			Description: "Explore codebase",
			Triggers:    []string{"/explore", "explore code"},
			Commands:    []string{"Task(Explore)"},
		},

		// Build skills
		{
			Name:        "build",
			Category:    CatBuild,
			Description: "Build project",
			Triggers:    []string{"/build", "build project"},
			Commands:    []string{"make", "go build", "npm run build"},
		},

		// Test skills
		{
			Name:        "test",
			Category:    CatTest,
			Description: "Run tests",
			Triggers:    []string{"/test", "run tests"},
			Commands:    []string{"go test", "npm test", "pytest"},
		},
	}
}
