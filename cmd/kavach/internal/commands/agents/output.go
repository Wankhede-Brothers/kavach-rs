// Package agents - output.go
// DACE: Single responsibility - TOON and Sutra Protocol output
package agents

import (
	"fmt"
	"strings"
	"time"
)

// OutputTOONList outputs all agents in TOON format
func OutputTOONList(agents []*Agent) {
	fmt.Println("[AGENTS]")
	fmt.Printf("count: %d\n", len(agents))
	fmt.Printf("date: %s\n", time.Now().Format("2006-01-02"))
	fmt.Println()

	byLevel := groupByLevel(agents)
	levels := []int{LevelNLU, LevelCEO, LevelEngineer, LevelReview, LevelAegis}

	for _, lvl := range levels {
		list := byLevel[lvl]
		if len(list) == 0 {
			continue
		}

		fmt.Printf("[%s]\n", LevelName(lvl))
		for _, a := range list {
			fmt.Printf("%s: %s (%s)\n", a.Name, a.Description, a.Model)
		}
		fmt.Println()
	}

	fmt.Println("[MODELS]")
	fmt.Println("opus: ceo,research-director,aegis-guardian")
	fmt.Println("sonnet: backend,frontend,devops,security,qa,code-reviewer")
	fmt.Println("haiku: nlu-intent-analyzer")
}

// OutputTOONSingle outputs single agent in TOON format
func OutputTOONSingle(a *Agent) {
	fmt.Printf("[AGENT:%s]\n", strings.ToUpper(a.Name))
	fmt.Printf("name: %s\n", a.Name)
	fmt.Printf("level: %d\n", a.Level)
	fmt.Printf("model: %s\n", a.Model)
	fmt.Printf("description: %s\n", a.Description)
	if a.Path != "" {
		fmt.Printf("source: %s\n", a.Path)
	}
	fmt.Println()

	if len(a.Triggers) > 0 {
		fmt.Println("[TRIGGERS]")
		for _, t := range a.Triggers {
			fmt.Printf("- %s\n", t)
		}
		fmt.Println()
	}

	if len(a.Tools) > 0 {
		fmt.Println("[TOOLS]")
		fmt.Println(strings.Join(a.Tools, ", "))
		fmt.Println()
	}

	if len(a.Research) > 0 {
		fmt.Println("[RESEARCH_CONTEXT]")
		for _, r := range a.Research {
			fmt.Printf("- %s\n", r)
		}
		fmt.Println()
	}

	if len(a.Patterns) > 0 {
		fmt.Println("[PATTERNS]")
		for _, p := range a.Patterns {
			fmt.Printf("- %s\n", p)
		}
	}
}

// OutputSutraList outputs all agents in Sutra Protocol
func OutputSutraList(agents []*Agent) {
	fmt.Println("[META]")
	fmt.Println("protocol: SP/3.0")
	fmt.Println("from: kavach/agents")
	fmt.Println("to: CEO")
	fmt.Printf("date: %s\n", time.Now().Format("2006-01-02"))
	fmt.Printf("count: %d\n", len(agents))
	fmt.Println()

	fmt.Println("[AGENT_HIERARCHY]")
	fmt.Println("L-1: nlu-intent-analyzer (haiku)")
	fmt.Println("L0: ceo, research-director (opus)")
	fmt.Println("L1: backend, frontend, devops, security, qa (sonnet)")
	fmt.Println("L1.5: code-reviewer (sonnet)")
	fmt.Println("L2: aegis-guardian (opus)")
	fmt.Println()

	fmt.Println("[DELEGATION_RULES]")
	fmt.Println("CEO_ONLY: orchestration,delegation,decisions")
	fmt.Println("ENGINEER: implementation,code_changes")
	fmt.Println("AEGIS: verification,quality_gate")
	fmt.Println()

	fmt.Println("[AVAILABLE]")
	for _, a := range agents {
		fmt.Printf("%s: %s\n", a.Name, a.Model)
	}
}

// OutputSutraSingle outputs single agent in Sutra Protocol
func OutputSutraSingle(a *Agent) {
	fmt.Println("[META]")
	fmt.Println("protocol: SP/3.0")
	fmt.Println("from: kavach/agents")
	fmt.Printf("to: %s\n", a.Name)
	fmt.Printf("date: %s\n", time.Now().Format("2006-01-02"))
	fmt.Println()

	fmt.Printf("[AGENT:%s]\n", strings.ToUpper(a.Name))
	fmt.Printf("name: %s\n", a.Name)
	fmt.Printf("level: %d\n", a.Level)
	fmt.Printf("model: %s\n", a.Model)
	fmt.Printf("desc: %s\n", a.Description)
	fmt.Println()

	fmt.Println("[CONTEXT]")
	fmt.Printf("tools: %s\n", strings.Join(a.Tools, ","))
	fmt.Printf("triggers: %s\n", strings.Join(a.Triggers, ","))
	fmt.Println()

	if len(a.Research) > 0 || len(a.Patterns) > 0 {
		fmt.Println("[INJECTED_CONTEXT]")
		fmt.Printf("research_entries: %d\n", len(a.Research))
		fmt.Printf("pattern_entries: %d\n", len(a.Patterns))
	}
}

func groupByLevel(agents []*Agent) map[int][]*Agent {
	byLevel := make(map[int][]*Agent)
	for _, a := range agents {
		byLevel[a.Level] = append(byLevel[a.Level], a)
	}
	return byLevel
}
