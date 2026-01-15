// Package config provides dynamic configuration loading from TOON files.
// NO HARDCODING - All patterns loaded from config/*.toon at runtime.
// P1 FIX #5: TTL-based cache invalidation for hot-reloading support.
package config

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// CacheTTL is the time-to-live for cached config (5 minutes default).
// P1 FIX #5: Allows hot-reloading of config files without restart.
const CacheTTL = 5 * time.Minute

// MaxCacheSize limits the number of cached config files.
// P1 FIX: Prevents unbounded memory growth from config cache.
const MaxCacheSize = 50

// cacheEntry holds cached data with timestamp.
type cacheEntry struct {
	data       map[string][]string
	timestamp  time.Time
	lastAccess time.Time // P1 FIX: Track last access for LRU eviction
}

var (
	configDir   string
	configCache = make(map[string]*cacheEntry)
	cacheMu     sync.RWMutex
)

func init() {
	if dir := os.Getenv("KAVACH_CONFIG_DIR"); dir != "" {
		configDir = dir
	} else {
		home, _ := os.UserHomeDir()
		candidates := []string{
			filepath.Join(home, ".config", "kavach"),
			"/etc/kavach",
			"./config",
		}
		for _, c := range candidates {
			if _, err := os.Stat(c); err == nil {
				configDir = c
				break
			}
		}
	}
}

// LoadPatterns loads patterns from a TOON config file dynamically.
// P1 FIX #5: TTL-based cache - reloads after CacheTTL expires.
// P1 FIX: Updates lastAccess on hit, evicts LRU when cache full.
func LoadPatterns(filename string) map[string][]string {
	cacheMu.RLock()
	if cached, ok := configCache[filename]; ok {
		// P1 FIX #5: Check if cache is still valid (within TTL)
		if time.Since(cached.timestamp) < CacheTTL {
			cacheMu.RUnlock()
			// P1 FIX: Update last access time (requires write lock)
			cacheMu.Lock()
			cached.lastAccess = time.Now()
			cacheMu.Unlock()
			return cached.data
		}
	}
	cacheMu.RUnlock()

	result := make(map[string][]string)
	paths := []string{
		filepath.Join(configDir, filename),
		filepath.Join(".", "config", filename),
	}

	var file *os.File
	var err error
	for _, p := range paths {
		file, err = os.Open(p)
		if err == nil {
			break
		}
	}
	if err != nil {
		return result
	}
	defer file.Close()

	var currentSection string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			currentSection = strings.Trim(line, "[]")
			continue
		}
		if currentSection != "" {
			// Handle "keywords: a,b,c" format
			if strings.HasPrefix(line, "keywords:") {
				for _, kw := range strings.Split(strings.TrimPrefix(line, "keywords:"), ",") {
					if kw = strings.TrimSpace(kw); kw != "" {
						result[currentSection] = append(result[currentSection], kw)
					}
				}
			} else {
				// Add all other lines (including those with colons like "cat:bat:reason")
				result[currentSection] = append(result[currentSection], line)
			}
		}
	}

	// P1 FIX #5: Store with timestamp for TTL-based invalidation
	now := time.Now()
	cacheMu.Lock()
	// P1 FIX: Evict LRU entry if cache is full
	if len(configCache) >= MaxCacheSize {
		evictLRU()
	}
	configCache[filename] = &cacheEntry{
		data:       result,
		timestamp:  now,
		lastAccess: now,
	}
	cacheMu.Unlock()
	return result
}

// evictLRU removes the least recently used cache entry.
// P1 FIX: Prevents unbounded cache growth.
// Must be called with cacheMu held.
func evictLRU() {
	var oldestKey string
	var oldestTime time.Time
	first := true

	for key, entry := range configCache {
		if first || entry.lastAccess.Before(oldestTime) {
			oldestKey = key
			oldestTime = entry.lastAccess
			first = false
		}
	}

	if oldestKey != "" {
		delete(configCache, oldestKey)
	}
}

// GetNLUPatterns returns NLU patterns for intent classification
func GetNLUPatterns() map[string][]string { return LoadPatterns("nlu-patterns.toon") }

// GetSkillPatterns returns skill detection patterns
func GetSkillPatterns() map[string][]string { return LoadPatterns("skill-patterns.toon") }

// GetAgentMappings returns agent configuration
func GetAgentMappings() map[string][]string { return LoadPatterns("agent-mappings.toon") }

// GetValidAgents returns list of valid agent names from config
func GetValidAgents() []string { return GetAgentMappings()["VALID:AGENTS"] }

// GetEngineers returns list of engineer agents from config
func GetEngineers() []string { return GetAgentMappings()["ENGINEERS:LIST"] }

// IsValidAgent checks if agent name is valid (from config)
func IsValidAgent(agent string) bool {
	for _, a := range GetValidAgents() {
		if a == agent {
			return true
		}
	}
	return false
}

// IsEngineer checks if agent is an engineer (from config)
func IsEngineer(agent string) bool {
	for _, e := range GetEngineers() {
		if e == agent {
			return true
		}
	}
	return false
}

// ClearCache clears config cache for immediate reloading.
// P1 FIX #5: Forces reload on next access.
func ClearCache() {
	cacheMu.Lock()
	configCache = make(map[string]*cacheEntry)
	cacheMu.Unlock()
}

// InvalidateFile invalidates cache for a specific file.
// P1 FIX #5: Allows targeted cache invalidation.
func InvalidateFile(filename string) {
	cacheMu.Lock()
	delete(configCache, filename)
	cacheMu.Unlock()
}

// GetRouterMappings returns router configuration from router-mappings.toon
func GetRouterMappings() map[string][]string {
	return LoadPatterns("router-mappings.toon")
}

// GetIntentSkillMappings returns intent to skill mappings (DACE: dynamic loading)
func GetIntentSkillMappings() map[string]string {
	data := GetRouterMappings()["SKILL:INTENT_MAPPINGS"]
	result := make(map[string]string)
	for _, line := range data {
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			result[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}
	return result
}

// GetComplexIndicators returns keywords indicating complex tasks
func GetComplexIndicators() []string {
	return GetRouterMappings()["COMPLEX:INDICATORS"]
}

// GetSkillAgentDefaults returns default skill for each agent type
func GetSkillAgentDefaults() map[string]string {
	data := GetRouterMappings()["SKILL:AGENT_DEFAULTS"]
	result := make(map[string]string)
	for _, line := range data {
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			result[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}
	return result
}

// GetSkillPreferredKeywords returns task types that prefer skills
func GetSkillPreferredKeywords() []string {
	return GetRouterMappings()["SKILL:PREFERRED_KEYWORDS"]
}

// GetValidSkills returns the list of valid skill names from config
func GetValidSkills() map[string]bool {
	data := LoadPatterns("valid-skills.toon")
	skills := data["VALID_SKILLS"]
	result := make(map[string]bool, len(skills))
	for _, skill := range skills {
		result[skill] = true
	}
	return result
}

// GetFrameworkPatterns returns framework patterns for research detection
func GetFrameworkPatterns() map[string][]string {
	return LoadPatterns("frameworks.toon")
}
