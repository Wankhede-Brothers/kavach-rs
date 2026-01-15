package session

import (
	"fmt"

	"github.com/claude/shared/pkg/enforce"
	"github.com/spf13/cobra"
)

var endCmd = &cobra.Command{
	Use:   "end",
	Short: "End session and save state",
	Long: `[END]
desc: End current session and persist state
hook: Stop
purpose: Save session state before termination

[SAVES]
- Session ID and date
- Research done state
- Memory queried state
- CEO/Aegis invocation states

[USAGE]
kavach session end

[OUTPUT]
[END]   date, session, project
[STATE] research, memory, ceo, aegis states`,
	Run: runEndCmd,
}

func runEndCmd(cmd *cobra.Command, args []string) {
	ctx := enforce.NewContext()
	session := enforce.GetOrCreateSession()

	fmt.Println("[END]")
	fmt.Printf("date: %s\nsession: %s\nproject: %s\n\n", ctx.Today, session.ID, session.Project)

	fmt.Println("[STATE]")
	fmt.Printf("research_done: %s\nmemory: %s\nceo: %s\naegis: %s\n",
		boolStr(session.ResearchDone), boolStr(session.MemoryQueried),
		boolStr(session.CEOInvoked), boolStr(session.AegisVerified))

	session.Save()
}
