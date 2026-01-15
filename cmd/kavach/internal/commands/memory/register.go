// Package memory provides memory bank subcommands.
package memory

import "github.com/spf13/cobra"

// Register adds all memory commands to the parent memory command.
func Register(memoryCmd *cobra.Command) {
	memoryCmd.AddCommand(bankCmd)
	memoryCmd.AddCommand(writeCmd)
	memoryCmd.AddCommand(rpcCmd)
	memoryCmd.AddCommand(stmCmd)
	memoryCmd.AddCommand(injectCmd)
	memoryCmd.AddCommand(specCmd)
	memoryCmd.AddCommand(kanbanCmd)
	memoryCmd.AddCommand(viewCmd)
	memoryCmd.AddCommand(syncCmd)
}
