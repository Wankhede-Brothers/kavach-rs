// Package agents - register.go
// DACE: Single responsibility - command registration and dispatch
package agents

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
	Use:   "agents",
	Short: "Dynamic agent management (DACE + SP/3.0)",
	Long: `[AGENTS]
desc: Dynamic Agentic Context Engineering for agent management
protocol: SP/3.0 (Sutra Protocol)
sources: ~/.claude/agents/, .claude/agents/, Memory Bank

[DACE_PRINCIPLES]
lazy_load:    Load agent context on-demand
skill_first:  Use kavach binary before spawning
on_demand:    Inject research/patterns when needed

[HIERARCHY]
L-1:   nlu-intent-analyzer (haiku)
L0:    ceo, research-director (opus)
L1:    backend, frontend, devops, security, qa (sonnet)
L1.5:  code-reviewer (sonnet)
L2:    aegis-guardian (opus)

[FLAGS]
--get <name>:  Get specific agent with full context
--sutra:       Output in Sutra Protocol format
--inject:      Inject research from Memory Bank

[USAGE]
kavach agents                        # TOON summary
kavach agents --get ceo              # Full CEO context
kavach agents --get ceo --inject     # CEO with research`,
	Run: run,
}

func init() {
	cmd.Flags().StringVar(&getFlag, "get", "", "Get specific agent")
	cmd.Flags().BoolVar(&sutraFlag, "sutra", false, "Sutra Protocol output")
	cmd.Flags().BoolVar(&injectFlag, "inject", false, "Inject research context")
}

// Cmd returns the agents command
func Cmd() *cobra.Command {
	return cmd
}

func run(c *cobra.Command, args []string) {
	agents := Discover()

	if getFlag != "" {
		agent := Find(agents, getFlag)
		if agent == nil {
			fmt.Printf("[ERROR] Agent not found: %s\n", getFlag)
			return
		}

		// DACE: Publish EventAgentInvoke for telemetry/hooks
		eventBus := events.GetEventBus()
		eventBus.Publish(events.EventAgentInvoke, "kavach", map[string]interface{}{
			"agent":  agent.Name,
			"level":  agent.Level,
			"inject": injectFlag,
			"sutra":  sutraFlag,
		})

		if injectFlag {
			InjectContext(agent)
		}

		if sutraFlag {
			OutputSutraSingle(agent)
		} else {
			OutputTOONSingle(agent)
		}
		return
	}

	if sutraFlag {
		OutputSutraList(agents)
	} else {
		OutputTOONList(agents)
	}
}
