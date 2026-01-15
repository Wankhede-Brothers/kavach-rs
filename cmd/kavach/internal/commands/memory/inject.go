package memory

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var injectCmd = &cobra.Command{
	Use:   "inject",
	Short: "RPC context injector",
	Long: `[INJECT]
desc: Inject context data into response
purpose: Add memory/context to hook additionalContext

[INPUT]
stdin: {"context":"session","data":{"user":"dev","project":"api"}}

[FIELDS]
context: Context type identifier
data:    Key-value pairs to inject

[USAGE]
echo '{"context":"memory","data":{"category":"decisions"}}' | kavach memory inject

[OUTPUT]
[INJECT]
context: <type>
<key>: <value>
...`,
	Run: runInjectCmd,
}

type InjectRequest struct {
	Context string                 `json:"context"`
	Data    map[string]interface{} `json:"data,omitempty"`
}

func runInjectCmd(cmd *cobra.Command, args []string) {
	var req InjectRequest
	if err := json.NewDecoder(os.Stdin).Decode(&req); err != nil {
		fmt.Println("[ERROR]")
		fmt.Printf("msg: %s\n", err.Error())
		os.Exit(1)
	}

	fmt.Println("[INJECT]")
	fmt.Printf("context: %s\n", req.Context)
	for k, v := range req.Data {
		fmt.Printf("%s: %v\n", k, v)
	}
}
