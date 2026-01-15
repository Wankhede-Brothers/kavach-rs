// Package agentic provides skill-first routing.
// DACE: All routing patterns loaded from config (NO HARDCODING).
package agentic

import (
	"strings"
	"sync"

	"github.com/claude/shared/pkg/config"
)

// RoutingDecision represents the result of routing a request.
type RoutingDecision struct {
	UseSkill    bool   // Prefer skill over agent
	SkillName   string // Skill to invoke
	AgentName   string // Agent to invoke (if no skill)
	RequiresCEO bool   // Needs CEO orchestration
	Reason      string // Why this decision was made
}

// SkillFirstRouter routes requests preferring skills over agents.
// CORE PRINCIPLE: Skills are lightweight, agents are heavy.
type SkillFirstRouter struct {
	loader        *DynamicLoader
	skillTriggers map[string]string   // keyword -> skill
	agentSkills   map[string][]string // agent -> required skills
	mu            sync.RWMutex
}

// NewSkillFirstRouter creates a router that prefers skills.
func NewSkillFirstRouter(loader *DynamicLoader) *SkillFirstRouter {
	return &SkillFirstRouter{
		loader:        loader,
		skillTriggers: make(map[string]string),
		agentSkills:   make(map[string][]string),
	}
}

// RegisterSkillTrigger maps a keyword to a skill.
func (r *SkillFirstRouter) RegisterSkillTrigger(keyword, skillName string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.skillTriggers[strings.ToLower(keyword)] = skillName
}

// RegisterAgentSkills maps an agent to its required skills.
func (r *SkillFirstRouter) RegisterAgentSkills(agentName string, skills []string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.agentSkills[agentName] = skills
}

// Route determines where to send a request.
// PRIORITY: 1. Skill (if trigger matches)
//  2. Agent's skill (if agent requested)
//  3. Agent (if complex task)
func (r *SkillFirstRouter) Route(intent string, keywords []string) *RoutingDecision {
	r.mu.RLock()
	defer r.mu.RUnlock()

	decision := &RoutingDecision{}

	// 1. Check for skill trigger matches
	for _, kw := range keywords {
		if skillName, ok := r.skillTriggers[strings.ToLower(kw)]; ok {
			decision.UseSkill = true
			decision.SkillName = skillName
			decision.Reason = "Keyword trigger matched skill"
			return decision
		}
	}

	// 2. Check if intent maps to a skill
	// DACE: Load from config instead of hardcoding
	intentLower := strings.ToLower(intent)
	skillMappings := config.GetIntentSkillMappings()

	for keyword, skill := range skillMappings {
		if strings.Contains(intentLower, keyword) {
			decision.UseSkill = true
			decision.SkillName = skill
			decision.Reason = "Intent matched skill domain"
			return decision
		}
	}

	// 3. Complex tasks require CEO orchestration
	// DACE: Load from config instead of hardcoding
	complexIndicators := config.GetComplexIndicators()
	for _, indicator := range complexIndicators {
		if strings.Contains(intentLower, indicator) {
			decision.RequiresCEO = true
			decision.AgentName = "ceo"
			decision.Reason = "Complex task requires orchestration"
			return decision
		}
	}

	// 4. Default: use skill if possible, agent otherwise
	decision.UseSkill = false
	decision.AgentName = "backend-engineer" // Default agent
	decision.Reason = "No specific skill match, using default agent"
	return decision
}

// GetSkillForAgent returns the primary skill an agent should use.
func (r *SkillFirstRouter) GetSkillForAgent(agentName string) string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if skills, ok := r.agentSkills[agentName]; ok && len(skills) > 0 {
		return skills[0]
	}

	// DACE: Load defaults from config instead of hardcoding
	defaults := config.GetSkillAgentDefaults()
	return defaults[agentName]
}

// ShouldPreferSkill returns true if skills should be preferred for this task.
// DACE: Load skill preference keywords from config instead of hardcoding.
func (r *SkillFirstRouter) ShouldPreferSkill(taskType string) bool {
	// Skills are preferred for research, validation, single-domain tasks
	// DACE: Load from config instead of hardcoding
	skillPreferred := config.GetSkillPreferredKeywords()
	for _, pref := range skillPreferred {
		if strings.Contains(strings.ToLower(taskType), pref) {
			return true
		}
	}
	return false
}
