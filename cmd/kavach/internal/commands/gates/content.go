package gates

import (
	"strings"

	"github.com/claude/shared/pkg/hook"
	"github.com/spf13/cobra"
)

var contentHookMode bool

var contentCmd = &cobra.Command{
	Use:   "content",
	Short: "Content validation gate",
	Long: `[CONTENT]
desc: Content safety and secrets detection
hook: PreToolUse:Write, PreToolUse:Edit
purpose: Prevent secrets, credentials, and sensitive data in code

[BLOCKS]
Secrets: password=, secret=, api_key=, token=, private_key
Keys:    BEGIN RSA PRIVATE, BEGIN OPENSSH PRIVATE

[WARNS]
Hardcoded: http://localhost, 127.0.0.1, 0.0.0.0 (suggest use_config)

[USAGE]
echo '{"tool_name":"Write","tool_input":{"file_path":"x.py","content":"password = secret"}}' | kavach gates content --hook

[OUTPUT]
approve: Content safe
block:   Sensitive pattern detected
modify:  Warning (hardcoded localhost)`,
	Run: runContentGate,
}

func init() {
	contentCmd.Flags().BoolVar(&contentHookMode, "hook", false, "Run in hook mode (JSON stdin/stdout)")
}

func runContentGate(cmd *cobra.Command, args []string) {
	if !contentHookMode {
		cmd.Help()
		return
	}

	input := hook.MustReadHookInput()

	if input.ToolName != "Write" && input.ToolName != "Edit" {
		hook.ExitApproveTOON("CONTENT")
	}

	content := input.GetString("content")
	if input.ToolName == "Edit" {
		content = input.GetString("new_string")
	}

	if content == "" {
		hook.ExitApproveTOON("CONTENT")
	}

	// Check for sensitive content patterns
	sensitivePatterns := []string{
		"password =",
		"secret =",
		"api_key =",
		"token =",
		"private_key",
		"BEGIN RSA PRIVATE",
		"BEGIN OPENSSH PRIVATE",
	}

	contentLower := strings.ToLower(content)
	for _, pattern := range sensitivePatterns {
		if strings.Contains(contentLower, strings.ToLower(pattern)) {
			hook.ExitBlockTOON("CONTENT", "sensitive:"+pattern)
		}
	}

	// Check for hardcoded IPs/URLs (warn only)
	if strings.Contains(content, "http://localhost") ||
		strings.Contains(content, "127.0.0.1") ||
		strings.Contains(content, "0.0.0.0") {
		hook.ExitModifyTOON("CONTENT", map[string]string{
			"warn": "hardcoded_localhost",
			"hint": "use_config",
		})
	}

	hook.ExitApproveTOON("CONTENT")
}
