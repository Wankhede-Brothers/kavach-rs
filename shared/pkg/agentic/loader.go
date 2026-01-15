// Package agentic provides Dynamic Agentic Context Engineering.
// Loads agents/skills ON DEMAND - never all at once.
package agentic

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/claude/shared/pkg/dsa"
)

// AgentDef represents a dynamically loaded agent definition.
type AgentDef struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Model       string   `json:"model"`
	Skills      []string `json:"skills"`
	Priority    int      `json:"priority"`
	Loaded      bool     `json:"-"`
}

// SkillDef represents a dynamically loaded skill definition.
type SkillDef struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Triggers    []string `json:"triggers"`
	AutoInvoke  bool     `json:"auto_invoke"`
	Content     string   `json:"-"` // Loaded on demand
	Loaded      bool     `json:"-"`
}

// DynamicLoader provides lazy loading for agents and skills.
// CORE PRINCIPLE: Load nothing until needed.
type DynamicLoader struct {
	agentDir   string
	skillDir   string
	agents     *dsa.LazyMap[string, *AgentDef]
	skills     *dsa.LazyMap[string, *SkillDef]
	skillIndex map[string]string // trigger -> skill name
	mu         sync.RWMutex
}

// NewDynamicLoader creates a loader with lazy initialization.
func NewDynamicLoader(agentDir, skillDir string) *DynamicLoader {
	dl := &DynamicLoader{
		agentDir:   agentDir,
		skillDir:   skillDir,
		skillIndex: make(map[string]string),
	}

	// Create lazy agent loader
	dl.agents = dsa.NewLazyMap(func(name string) func() (*AgentDef, error) {
		return func() (*AgentDef, error) {
			return dl.loadAgent(name)
		}
	})

	// Create lazy skill loader
	dl.skills = dsa.NewLazyMap(func(name string) func() (*SkillDef, error) {
		return func() (*SkillDef, error) {
			return dl.loadSkill(name)
		}
	})

	return dl
}

// loadAgent loads a single agent definition on demand.
func (dl *DynamicLoader) loadAgent(name string) (*AgentDef, error) {
	path := filepath.Join(dl.agentDir, name+".md")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Parse frontmatter (simplified)
	agent := &AgentDef{
		Name:   name,
		Loaded: true,
	}
	// Extract description from first line after ---
	agent.Description = extractDescription(string(data))

	return agent, nil
}

// loadSkill loads a single skill definition on demand.
func (dl *DynamicLoader) loadSkill(name string) (*SkillDef, error) {
	path := filepath.Join(dl.skillDir, name, "SKILL.md")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	skill := &SkillDef{
		Name:    name,
		Content: string(data),
		Loaded:  true,
	}
	skill.Description = extractDescription(string(data))
	skill.Triggers = extractTriggers(string(data))

	// Index triggers for fast lookup
	dl.mu.Lock()
	for _, trigger := range skill.Triggers {
		dl.skillIndex[trigger] = name
	}
	dl.mu.Unlock()

	return skill, nil
}

// GetAgent retrieves an agent, loading it if needed.
func (dl *DynamicLoader) GetAgent(name string) (*AgentDef, error) {
	return dl.agents.Get(name)
}

// GetSkill retrieves a skill, loading it if needed.
func (dl *DynamicLoader) GetSkill(name string) (*SkillDef, error) {
	return dl.skills.Get(name)
}

// FindSkillByTrigger finds a skill that matches a trigger keyword.
// Returns the skill name if found, empty string otherwise.
func (dl *DynamicLoader) FindSkillByTrigger(trigger string) string {
	dl.mu.RLock()
	defer dl.mu.RUnlock()
	return dl.skillIndex[trigger]
}

// IsAgentLoaded checks if an agent is currently in memory.
func (dl *DynamicLoader) IsAgentLoaded(name string) bool {
	return dl.agents.IsLoaded(name)
}

// IsSkillLoaded checks if a skill is currently in memory.
func (dl *DynamicLoader) IsSkillLoaded(name string) bool {
	return dl.skills.IsLoaded(name)
}

// LoadedAgents returns names of agents currently in memory.
func (dl *DynamicLoader) LoadedAgents() []string {
	return dl.agents.LoadedKeys()
}

// LoadedSkills returns names of skills currently in memory.
func (dl *DynamicLoader) LoadedSkills() []string {
	return dl.skills.LoadedKeys()
}

// Helper: extract description from markdown content
func extractDescription(content string) string {
	// Find description in frontmatter
	lines := splitLines(content)
	for _, line := range lines {
		if len(line) > 12 && line[:12] == "description:" {
			return line[12:]
		}
	}
	return ""
}

// Helper: extract triggers from skill content
func extractTriggers(content string) []string {
	// Find triggers: [...] in skill definition
	var triggers []string
	// Simplified extraction - look for triggers: [...]
	return triggers
}

// Helper: split lines
func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}
