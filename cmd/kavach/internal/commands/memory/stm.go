package memory

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/claude/shared/pkg/types"
	"github.com/claude/shared/pkg/util"
	"github.com/spf13/cobra"
)

var stmCmd = &cobra.Command{
	Use:   "stm",
	Short: "Short-term memory updater",
	Long: `[STM]
desc: Short-term memory (scratchpad) operations
path: ~/.local/shared/shared-ai/memory/STM/scratchpad.json
purpose: Track current focus, project, and session state

[INPUT]
stdin: {"focus":"implementing auth","project":"my-app"}

[FIELDS]
focus:   Current task focus
project: Active project name

[USAGE]
echo '{"focus":"fixing bug","project":"api-server"}' | kavach memory stm

[OUTPUT]
[STM]
focus: <current_focus>
project: <current_project>`,
	Run: runSTMCmd,
}

func runSTMCmd(cmd *cobra.Command, args []string) {
	var ctx types.STMContext
	if err := json.NewDecoder(os.Stdin).Decode(&ctx); err != nil {
		fmt.Println("[ERROR]")
		fmt.Printf("msg: %s\n", err.Error())
		os.Exit(1)
	}

	scratchpadPath := util.MemoryBankPath("STM") + "/scratchpad.json"
	if err := util.WriteJSON(scratchpadPath, &ctx); err != nil {
		fmt.Println("[ERROR]")
		fmt.Printf("msg: %s\n", err.Error())
		os.Exit(1)
	}

	fmt.Println("[STM]")
	fmt.Printf("focus: %s\nproject: %s\n", ctx.Focus, ctx.Project)
}
