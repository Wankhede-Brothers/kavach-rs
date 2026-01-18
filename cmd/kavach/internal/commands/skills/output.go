// Package skills - output.go
// DACE: Single responsibility - TOON and Sutra Protocol output
package skills

import (
	"fmt"
	"strings"
	"time"
)

// OutputTOONList outputs all skills in TOON format
func OutputTOONList(skills []*Skill) {
	fmt.Println("[SKILLS]")
	fmt.Printf("count: %d\n", len(skills))
	fmt.Printf("date: %s\n", time.Now().Format("2006-01-02"))
	fmt.Println()

	byCategory := ByCategory(skills)
	categories := []string{CatGit, CatSession, CatMemory, CatResearch, CatBuild, CatTest}

	for _, cat := range categories {
		list := byCategory[cat]
		if len(list) == 0 {
			continue
		}

		fmt.Printf("[%s]\n", strings.ToUpper(cat))
		for _, s := range list {
			fmt.Printf("%s: %s\n", s.Name, s.Description)
		}
		fmt.Println()
	}

	// Print uncategorized
	if general := byCategory["general"]; len(general) > 0 {
		fmt.Println("[GENERAL]")
		for _, s := range general {
			fmt.Printf("%s: %s\n", s.Name, s.Description)
		}
	}
}

// OutputTOONSingle outputs single skill in TOON format
func OutputTOONSingle(s *Skill) {
	fmt.Printf("[SKILL:%s]\n", strings.ToUpper(s.Name))
	fmt.Printf("name: %s\n", s.Name)
	fmt.Printf("category: %s\n", s.Category)
	fmt.Printf("description: %s\n", s.Description)
	if s.Path != "" {
		fmt.Printf("source: %s\n", s.Path)
	}
	fmt.Println()

	if len(s.Triggers) > 0 {
		fmt.Println("[TRIGGERS]")
		for _, t := range s.Triggers {
			fmt.Printf("- %s\n", t)
		}
		fmt.Println()
	}

	if len(s.Commands) > 0 {
		fmt.Println("[COMMANDS]")
		for _, c := range s.Commands {
			fmt.Printf("- %s\n", c)
		}
		fmt.Println()
	}

	if len(s.Research) > 0 {
		fmt.Println("[RESEARCH_CONTEXT]")
		for _, r := range s.Research {
			fmt.Printf("- %s\n", r)
		}
		fmt.Println()
	}

	if len(s.Patterns) > 0 {
		fmt.Println("[PATTERNS]")
		for _, p := range s.Patterns {
			fmt.Printf("- %s\n", p)
		}
	}
}

// OutputSutraList outputs all skills in Sutra Protocol
func OutputSutraList(skills []*Skill) {
	fmt.Println("[META]")
	fmt.Println("protocol: SP/1.0")
	fmt.Println("from: kavach/skills")
	fmt.Println("to: Claude")
	fmt.Printf("date: %s\n", time.Now().Format("2006-01-02"))
	fmt.Printf("count: %d\n", len(skills))
	fmt.Println()

	fmt.Println("[AVAILABLE_SKILLS]")
	for _, s := range skills {
		triggers := strings.Join(s.Triggers, ",")
		if triggers == "" {
			triggers = "/" + s.Name
		}
		fmt.Printf("%s: %s\n", s.Name, triggers)
	}
	fmt.Println()

	fmt.Println("[CATEGORIES]")
	byCategory := ByCategory(skills)
	for cat, list := range byCategory {
		names := make([]string, len(list))
		for i, s := range list {
			names[i] = s.Name
		}
		fmt.Printf("%s: %s\n", cat, strings.Join(names, ","))
	}
}

// OutputSutraSingle outputs single skill in Sutra Protocol
func OutputSutraSingle(s *Skill) {
	fmt.Println("[META]")
	fmt.Println("protocol: SP/1.0")
	fmt.Println("from: kavach/skills")
	fmt.Printf("skill: %s\n", s.Name)
	fmt.Printf("date: %s\n", time.Now().Format("2006-01-02"))
	fmt.Println()

	fmt.Printf("[SKILL:%s]\n", strings.ToUpper(s.Name))
	fmt.Printf("name: %s\n", s.Name)
	fmt.Printf("category: %s\n", s.Category)
	fmt.Printf("desc: %s\n", s.Description)
	fmt.Println()

	fmt.Println("[CONTEXT]")
	fmt.Printf("triggers: %s\n", strings.Join(s.Triggers, ","))
	fmt.Printf("commands: %s\n", strings.Join(s.Commands, ","))
	fmt.Println()

	if len(s.Research) > 0 || len(s.Patterns) > 0 {
		fmt.Println("[INJECTED_CONTEXT]")
		fmt.Printf("research_entries: %d\n", len(s.Research))
		fmt.Printf("pattern_entries: %d\n", len(s.Patterns))
	}
}
