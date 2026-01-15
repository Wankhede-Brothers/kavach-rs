// Package skills - register.go
// DACE: Single responsibility - command registration and dispatch
package skills

import (
	"fmt"

	"github.com/claude/shared/events"
	"github.com/spf13/cobra"
)

var (
	getFlag    string
	sutraFlag  bool
	injectFlag bool
)

var cmd = &cobra.Command{
	Use:   "skills",
	Short: "Dynamic skill management (DACE + SP/3.0)",
	Long: `[SKILLS]
desc: Dynamic Agentic Context Engineering for skill management
protocol: SP/3.0 (Sutra Protocol)
sources: ~/.claude/skills/, .claude/skills/, Memory Bank
file: SKILL.md (case-sensitive for Claude Code)

[DACE_PRINCIPLES]
lazy_load:    Load skill context on-demand
skill_first:  Use kavach binary before spawning agents
on_demand:    Inject research/patterns when needed

[CATEGORIES]
git:      commit, review-pr, create-pr
session:  init, status
memory:   memory, kanban
research: plan, explore
build:    build
test:     test

[FLAGS]
--get <name>:  Get specific skill with full context
--sutra:       Output in Sutra Protocol format
--inject:      Inject research from Memory Bank

[USAGE]
kavach skills                        # TOON summary
kavach skills --get rust             # Full rust skill context
kavach skills --get rust --inject    # Rust with research`,
	Run: run,
}

func init() {
	cmd.Flags().StringVar(&getFlag, "get", "", "Get specific skill")
	cmd.Flags().BoolVar(&sutraFlag, "sutra", false, "Sutra Protocol output")
	cmd.Flags().BoolVar(&injectFlag, "inject", false, "Inject research context")
}

// Cmd returns the skills command
func Cmd() *cobra.Command {
	return cmd
}

func run(c *cobra.Command, args []string) {
	skills := Discover()

	if getFlag != "" {
		skill := Find(skills, getFlag)
		if skill == nil {
			fmt.Printf("[ERROR] Skill not found: %s\n", getFlag)
			return
		}

		// DACE: Publish EventSkillInvoke for telemetry/hooks
		eventBus := events.GetEventBus()
		eventBus.Publish(events.EventSkillInvoke, "kavach", map[string]interface{}{
			"skill":  skill.Name,
			"inject": injectFlag,
			"sutra":  sutraFlag,
		})

		if injectFlag {
			InjectContext(skill)
		}

		if sutraFlag {
			OutputSutraSingle(skill)
		} else {
			OutputTOONSingle(skill)
		}
		return
	}

	if sutraFlag {
		OutputSutraList(skills)
	} else {
		OutputTOONList(skills)
	}
}
