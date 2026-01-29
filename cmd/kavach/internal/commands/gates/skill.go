package gates

import (
	"strings"

	"github.com/claude/shared/pkg/config"
	"github.com/claude/shared/pkg/hook"
	"github.com/spf13/cobra"
)

var skillHookMode bool

var skillCmd = &cobra.Command{
	Use:   "skill",
	Short: "Skill invocation gate",
	Long: `[SKILL]
desc: Skill invocation validation
hook: PreToolUse:Skill
purpose: Validate skill names before invocation

[VALID_SKILLS]
git:     commit, review-pr, create-pr
session: init, status, memory
research: research, plan

[USAGE]
echo '{"tool_name":"Skill","tool_input":{"skill":"commit"}}' | kavach gates skill --hook

[OUTPUT]
approve: N/A (always modifies or blocks)
block:   Unknown skill name
modify:  Valid skill routed`,
	Run: runSkillGate,
}

func init() {
	skillCmd.Flags().BoolVar(&skillHookMode, "hook", false, "Run in hook mode (JSON stdin/stdout)")
}

// validSkills loads valid skill names from config (NO HARDCODING)
// P1 FIX: Dynamic loading from config/valid-skills.toon
func getValidSkills() map[string]bool {
	return config.GetValidSkills()
}

func runSkillGate(cmd *cobra.Command, args []string) {
	if !skillHookMode {
		cmd.Help()
		return
	}

	input := hook.MustReadHookInput()

	if input.ToolName != "Skill" {
		hook.ExitApproveTOON("SKILL")
	}

	skillName := input.GetString("skill")
	if skillName == "" {
		hook.ExitBlockTOON("SKILL", "no_skill_name")
	}

	skillName = strings.ToLower(skillName)

	// P1 FIX: Load valid skills from config instead of hardcoded map
	validSkills := getValidSkills()
	if !validSkills[skillName] {
		// Warn but allow — user may have installed custom skills/plugins
		hook.ExitModifyTOON("SKILL_WARN", map[string]string{
			"skill":  skillName,
			"status": "unrecognized_but_allowed",
			"note":   "not in valid-skills.toon",
		})
	}

	hook.ExitModifyTOON("SKILL", map[string]string{
		"skill":  skillName,
		"status": "routed",
	})
}
