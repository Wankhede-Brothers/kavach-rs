// Package agentic provides integration for Go binaries.
package agentic

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// AgenticSystem coordinates dynamic agent/skill loading with research gates.
// CORE PRINCIPLE: Everything is lazy-loaded and research-verified.
type AgenticSystem struct {
	loader  *DynamicLoader
	router  *SkillFirstRouter
	gate    *ResearchGate
	baseDir string
	mu      sync.RWMutex
}

// SystemConfig configures the agentic system.
type SystemConfig struct {
	AgentDir string // Directory containing agent .md files
	SkillDir string // Directory containing skill directories
}

// NewAgenticSystem creates a coordinated agentic system.
func NewAgenticSystem(cfg SystemConfig) *AgenticSystem {
	loader := NewDynamicLoader(cfg.AgentDir, cfg.SkillDir)
	router := NewSkillFirstRouter(loader)
	gate := NewResearchGate()

	return &AgenticSystem{
		loader:  loader,
		router:  router,
		gate:    gate,
		baseDir: filepath.Dir(cfg.AgentDir),
	}
}

// DefaultSystemConfig returns config for opencode structure.
func DefaultSystemConfig() SystemConfig {
	home, _ := os.UserHomeDir()
	return SystemConfig{
		AgentDir: filepath.Join(home, ".config/opencode/agent"),
		SkillDir: filepath.Join(home, ".config/opencode/skill"),
	}
}

// RouteRequest routes a user request to the appropriate handler.
// Returns routing decision with skill-first preference.
func (sys *AgenticSystem) RouteRequest(intent string, keywords []string) *RoutingDecision {
	return sys.router.Route(intent, keywords)
}

// GetSkill retrieves a skill (lazy loaded).
func (sys *AgenticSystem) GetSkill(name string) (*SkillDef, error) {
	return sys.loader.GetSkill(name)
}

// GetAgent retrieves an agent (lazy loaded).
func (sys *AgenticSystem) GetAgent(name string) (*AgentDef, error) {
	return sys.loader.GetAgent(name)
}

// RequireResearch checks if task needs research before implementation.
func (sys *AgenticSystem) RequireResearch(task string) *ResearchRequirement {
	return sys.gate.RequireResearch(task)
}

// ValidateResponse checks response for research compliance.
func (sys *AgenticSystem) ValidateResponse(response string) *ValidationResult {
	result := &ValidationResult{
		Valid:      true,
		Violations: []string{},
	}

	// Check for forbidden phrases
	forbidden := sys.gate.CheckForForbiddenPhrases(response)
	if len(forbidden) > 0 {
		result.Valid = false
		result.Violations = append(result.Violations, forbidden...)
		result.Reason = "Response contains forbidden phrases indicating skipped research"
	}

	// Check for research evidence
	if !sys.gate.ValidateResearchDone(response) {
		result.Valid = false
		result.Violations = append(result.Violations, "no-research-evidence")
		result.Reason = "Response lacks evidence of research (WebSearch/WebFetch)"
	}

	return result
}

// ValidationResult holds response validation outcome.
type ValidationResult struct {
	Valid      bool
	Violations []string
	Reason     string
}

// BuildSearchQuery creates a research query with date injection.
func (sys *AgenticSystem) BuildSearchQuery(topic string) string {
	return sys.gate.BuildSearchQuery(topic)
}

// Today returns the injected current date.
func (sys *AgenticSystem) Today() string {
	return sys.gate.Today()
}

// Year returns the injected current year.
func (sys *AgenticSystem) Year() string {
	return sys.gate.Year()
}

// Stats returns system statistics.
func (sys *AgenticSystem) Stats() map[string]interface{} {
	return map[string]interface{}{
		"loaded_agents": sys.loader.LoadedAgents(),
		"loaded_skills": sys.loader.LoadedSkills(),
		"today":         sys.gate.Today(),
		"year":          sys.gate.Year(),
	}
}

// ProcessRequest handles a complete request through the agentic pipeline.
// This is the main entry point for integrating with Go binaries.
func (sys *AgenticSystem) ProcessRequest(ctx context.Context, req *Request) (*Response, error) {
	resp := &Response{
		RequestID: req.ID,
		Date:      sys.gate.Today(),
	}

	// 1. Check if research is required
	if resReq := sys.gate.RequireResearch(req.Task); resReq != nil && resReq.Mandatory {
		resp.RequiresResearch = true
		resp.ResearchQuery = sys.gate.BuildSearchQuery(resReq.Topic)
		resp.ResearchReason = resReq.Reason
	}

	// 2. Route to skill or agent
	decision := sys.router.Route(req.Intent, req.Keywords)
	resp.UseSkill = decision.UseSkill
	resp.SkillName = decision.SkillName
	resp.AgentName = decision.AgentName
	resp.RequiresCEO = decision.RequiresCEO
	resp.RoutingReason = decision.Reason

	// 3. Load skill/agent content if needed
	if decision.UseSkill && decision.SkillName != "" {
		skill, err := sys.loader.GetSkill(decision.SkillName)
		if err != nil {
			return nil, fmt.Errorf("failed to load skill %s: %w", decision.SkillName, err)
		}
		resp.SkillContent = skill.Content
	}

	if decision.AgentName != "" && !decision.UseSkill {
		agent, err := sys.loader.GetAgent(decision.AgentName)
		if err != nil {
			return nil, fmt.Errorf("failed to load agent %s: %w", decision.AgentName, err)
		}
		resp.AgentDescription = agent.Description
	}

	return resp, nil
}

// Request represents an incoming request to the agentic system.
type Request struct {
	ID       string   // Unique request ID
	Task     string   // User's task description
	Intent   string   // Parsed intent
	Keywords []string // Extracted keywords
}

// Response represents the agentic system's routing decision.
type Response struct {
	RequestID        string // Original request ID
	Date             string // Today's date (injected)
	RequiresResearch bool   // Whether research is mandatory
	ResearchQuery    string // Suggested search query
	ResearchReason   string // Why research is needed
	UseSkill         bool   // True if skill should be used
	SkillName        string // Skill to invoke
	SkillContent     string // Loaded skill content
	AgentName        string // Agent to invoke (if no skill)
	AgentDescription string // Loaded agent description
	RequiresCEO      bool   // Needs CEO orchestration
	RoutingReason    string // Why this routing decision
}

// RegisterDefaultTriggers sets up standard skill triggers.
func (sys *AgenticSystem) RegisterDefaultTriggers() {
	triggers := map[string]string{
		"rust":       "backend",
		"axum":       "backend",
		"tonic":      "backend",
		"grpc":       "backend",
		"react":      "frontend",
		"typescript": "frontend",
		"vue":        "frontend",
		"sql":        "sql",
		"postgres":   "sql",
		"database":   "sql",
		"security":   "security",
		"auth":       "security",
		"test":       "testing",
		"pytest":     "testing",
		"terraform":  "cloud-infrastructure-mastery",
		"docker":     "cloud-infrastructure-mastery",
		"k8s":        "cloud-infrastructure-mastery",
		"kubernetes": "cloud-infrastructure-mastery",
		"api":        "api-design",
		"rest":       "api-design",
		"debug":      "debug-like-expert",
		"arch":       "arch",
		"dsa":        "dsa",
	}

	for keyword, skill := range triggers {
		sys.router.RegisterSkillTrigger(keyword, skill)
	}
}
