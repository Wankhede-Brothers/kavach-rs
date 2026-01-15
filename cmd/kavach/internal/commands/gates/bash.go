// Package gates provides hook gates for Claude Code.
// bash.go: Bash command sanitizer with Rust CLI enforcement.
// NO HARDCODING - Legacy/Rust mappings from config/rust-cli.toon
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

	if patterns.IsBlocked(command) {
		hook.ExitBlockTOON("BASH", "blocked_command")
	}

	if legacy, rust, reason := detectLegacyCommand(command); legacy != "" {
		msg := "LEGACY_BLOCKED:" + legacy + ":USE:" + rust + ":" + reason
		hook.ExitBlockTOON("RUST_CLI", msg)
	}

	if strings.HasPrefix(strings.TrimSpace(command), "sudo") {
		hook.ExitModifyTOON("BASH", map[string]string{
			"warn": "sudo_detected",
		})
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
