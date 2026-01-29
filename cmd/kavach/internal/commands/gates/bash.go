// Package gates provides hook gates for Claude Code.
// bash.go: Bash command sanitizer with Rust CLI enforcement.
// DACE: Uses shared/pkg/config for JSON-based security rules.
package gates

import (
	"strings"

	"github.com/claude/shared/pkg/config"
	"github.com/claude/shared/pkg/hook"
	"github.com/claude/shared/pkg/patterns"
	"github.com/spf13/cobra"
)

var bashHookMode bool

var bashCmd = &cobra.Command{
	Use:   "bash",
	Short: "Bash command sanitizer with Rust CLI enforcement",
	Run:   runBashGate,
}

func init() {
	bashCmd.Flags().BoolVar(&bashHookMode, "hook", false, "Hook mode")
}

func runBashGate(cmd *cobra.Command, args []string) {
	if !bashHookMode {
		cmd.Help()
		return
	}

	input := hook.MustReadHookInput()

	if input.ToolName != "Bash" {
		hook.ExitSilent()
	}

	command := input.GetString("command")
	if command == "" {
		hook.ExitBlockTOON("BASH", "empty_command")
	}

	// Check blocked commands from gates/config.json (priority)
	if config.IsBlockedCommand(command) {
		hook.ExitBlockTOON("BASH", "blocked_command")
	}

	// Legacy: Check patterns from patterns.toon
	if patterns.IsBlocked(command) {
		hook.ExitBlockTOON("BASH", "blocked_command")
	}

	// Check for legacy CLI commands that should use Rust alternatives
	if legacy, rust, reason := detectLegacyCommand(command); legacy != "" {
		msg := "LEGACY_BLOCKED:" + legacy + ":USE:" + rust + ":" + reason
		hook.ExitBlockTOON("RUST_CLI", msg)
	}

	// Warn on sudo commands
	if strings.HasPrefix(strings.TrimSpace(command), "sudo") {
		hook.ExitModifyTOON("BASH", map[string]string{
			"warn": "sudo_detected",
		})
	}

	// Warn on other risky patterns from config
	cfg := config.LoadGatesConfig()
	cmdLower := strings.ToLower(command)
	for _, warn := range cfg.Bash.WarnCommands {
		if strings.Contains(cmdLower, strings.ToLower(warn)) {
			hook.ExitModifyTOON("BASH", map[string]string{
				"warn": warn + "_detected",
			})
		}
	}

	hook.ExitSilent()
}

func detectLegacyCommand(command string) (string, string, string) {
	cfg := config.LoadPatterns("rust-cli.toon")
	blocked := cfg["LEGACY:BLOCKED"]
	allowed := cfg["ALLOWED:LEGACY"]

	parts := strings.Fields(strings.TrimSpace(command))
	if len(parts) == 0 {
		return "", "", ""
	}

	fullCmd := parts[0]
	cmdName := fullCmd

	if strings.Contains(fullCmd, "/") {
		cmdParts := strings.Split(fullCmd, "/")
		cmdName = cmdParts[len(cmdParts)-1]
	}

	for _, a := range allowed {
		if fullCmd == a || cmdName == a {
			return "", "", ""
		}
	}

	for _, mapping := range blocked {
		mparts := strings.SplitN(mapping, ":", 3)
		if len(mparts) >= 2 && cmdName == mparts[0] {
			reason := ""
			if len(mparts) >= 3 {
				reason = mparts[2]
			}
			return mparts[0], mparts[1], reason
		}
	}

	return "", "", ""
}
