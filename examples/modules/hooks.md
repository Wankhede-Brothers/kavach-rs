# Hooks Reference - SP/3.0
# Reference only, not loaded into context

HOOKS{event,command,purpose}
  SessionStart,kavach session init,Inject date + load memory
  UserPromptSubmit,kavach intent --hook,Classify intent
  PreToolUse,kavach gates enforcer --hook,Pipeline enforcement
  PreToolUse:Task,kavach gates ceo --hook,Validate delegation
  PreToolUse:Bash,kavach gates bash --hook,Block dangerous commands
  PreToolUse:Read,kavach gates read --hook,Block sensitive files
  PreToolUse:Write|Edit,kavach gates enforcer --hook,AST + research gate
  PostToolUse,kavach tracking context,Track context usage
  PreCompact,kavach session compact,Save state to Memory Bank
  Stop,kavach session end,Save session state

HOOK:OUTPUT_SCHEMA
  decision: "approve" | "block"
  reason: string
  additionalContext: string (injected to model)

GATES:AVAILABLE
  enforcer: Full pipeline (intent→ceo→quality→aegis)
  ceo: Task orchestration, subagent validation
  ast: AST validation for code changes
  bash: Command sanitization
  read: File access control
  intent: Intent classification
  skill: Skill invocation validation
