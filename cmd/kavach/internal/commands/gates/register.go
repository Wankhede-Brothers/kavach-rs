// Package gates provides gate enforcement subcommands.
package gates

import "github.com/spf13/cobra"

// Register adds all gate commands to the parent gates command.
func Register(gatesCmd *cobra.Command) {
	// Umbrella gates (4 — called by hooks in settings.json)
	gatesCmd.AddCommand(preWriteCmd)  // PreToolUse:Write|Edit|NotebookEdit
	gatesCmd.AddCommand(postWriteCmd) // PostToolUse:Write|Edit|NotebookEdit
	gatesCmd.AddCommand(preToolCmd)   // PreToolUse:Bash|Read|Glob|Grep|Task|Skill|...
	gatesCmd.AddCommand(postToolCmd)  // PostToolUse:Bash|Read|Glob|Grep|Task|WebSearch|...

	// Intent gate (standalone — UserPromptSubmit)
	gatesCmd.AddCommand(intentCmd)

	// Legacy individual gates (kept for direct invocation / testing)
	gatesCmd.AddCommand(ceoCmd)
	gatesCmd.AddCommand(astCmd)
	gatesCmd.AddCommand(bashCmd)
	gatesCmd.AddCommand(readCmd)
	gatesCmd.AddCommand(skillCmd)
	gatesCmd.AddCommand(lintCmd)
	gatesCmd.AddCommand(researchCmd)
	gatesCmd.AddCommand(contentCmd)
	gatesCmd.AddCommand(qualityCmd)
	gatesCmd.AddCommand(enforcerCmd)
	gatesCmd.AddCommand(contextCmd)
	gatesCmd.AddCommand(dagCmd)
	gatesCmd.AddCommand(taskCmd)
	gatesCmd.AddCommand(codeGuardCmd)
	gatesCmd.AddCommand(chainCmd)
	gatesCmd.AddCommand(subagentCmd)
	gatesCmd.AddCommand(failureCmd)
	gatesCmd.AddCommand(mockdataCmd)
}
