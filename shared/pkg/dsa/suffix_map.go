package dsa

import (
	"path/filepath"
	"strings"
	"sync"
)

// SuffixMap provides O(1) lookup for file extensions and patterns.
type SuffixMap[V any] struct {
	exact    map[string]V // Exact filename matches (e.g., "package-lock.json")
	suffix   map[string]V // Suffix matches (e.g., ".lock")
	contains map[string]V // Contains matches (e.g., "node_modules")
	mu       sync.RWMutex
}

// NewSuffixMap creates a new SuffixMap with optional initial capacity.
func NewSuffixMap[V any](capacity int) *SuffixMap[V] {
	if capacity <= 0 {
		capacity = 16
	}
	return &SuffixMap[V]{
		exact:    make(map[string]V, capacity),
		suffix:   make(map[string]V, capacity),
		contains: make(map[string]V, capacity),
	}
}

// AddExact adds an exact filename match pattern.
func (m *SuffixMap[V]) AddExact(pattern string, value V) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.exact[pattern] = value
}

// AddSuffix adds a suffix match pattern (e.g., ".lock", "-lock.json").
func (m *SuffixMap[V]) AddSuffix(pattern string, value V) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.suffix[pattern] = value
}

// AddContains adds a contains match pattern (e.g., "node_modules").
func (m *SuffixMap[V]) AddContains(pattern string, value V) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.contains[pattern] = value
}

// Get checks if a path matches any pattern and returns the associated value.
// Checks in order: exact match, suffix match, contains match.
func (m *SuffixMap[V]) Get(path string) (V, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Normalize path
	base := filepath.Base(path)

	// 1. Exact match on base filename
	if v, ok := m.exact[base]; ok {
		return v, true
	}

	// 2. Suffix match (extension or pattern)
	for suffix, v := range m.suffix {
		if strings.HasSuffix(base, suffix) || strings.HasSuffix(path, suffix) {
			return v, true
		}
	}

	// 3. Contains match (directory patterns)
	for pattern, v := range m.contains {
		if strings.Contains(path, pattern) {
			return v, true
		}
	}

	var zero V
	return zero, false
}

// GetExact checks only exact filename matches.
func (m *SuffixMap[V]) GetExact(filename string) (V, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	v, ok := m.exact[filename]
	return v, ok
}

// GetBySuffix checks only suffix matches.
func (m *SuffixMap[V]) GetBySuffix(path string) (V, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for suffix, v := range m.suffix {
		if strings.HasSuffix(path, suffix) {
			return v, true
		}
	}
	var zero V
	return zero, false
}

// MatchExtension checks if a path's extension matches any suffix pattern.
func (m *SuffixMap[V]) MatchExtension(path string) (V, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	ext := filepath.Ext(path)
	if ext == "" {
		var zero V
		return zero, false
	}

	if v, ok := m.suffix[ext]; ok {
		return v, true
	}
	var zero V
	return zero, false
}

// Size returns total number of patterns.
func (m *SuffixMap[V]) Size() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.exact) + len(m.suffix) + len(m.contains)
}

// Clear removes all patterns.
func (m *SuffixMap[V]) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.exact = make(map[string]V, 16)
	m.suffix = make(map[string]V, 16)
	m.contains = make(map[string]V, 16)
}

// PatternMatcher provides a combined matcher using multiple strategies.
type PatternMatcher[V any] struct {
	bloom       *BloomFilter
	suffixMap   *SuffixMap[V]
	trie        *Trie
	hasContains bool // Track if any contains patterns exist
	mu          sync.RWMutex
}

// NewPatternMatcher creates a new PatternMatcher optimized for the given pattern count.
func NewPatternMatcher[V any](expectedPatterns int) *PatternMatcher[V] {
	return &PatternMatcher[V]{
		bloom:     NewBloomFilter(expectedPatterns, 0.01),
		suffixMap: NewSuffixMap[V](expectedPatterns),
		trie:      NewTrie(),
	}
}

// AddPattern adds a pattern with its value and match strategy.
func (pm *PatternMatcher[V]) AddPattern(pattern string, value V, matchType string) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	// Add to Bloom filter for fast negative check
	pm.bloom.Add(pattern)
	pm.bloom.Add(filepath.Base(pattern))
	if ext := filepath.Ext(pattern); ext != "" {
		pm.bloom.Add(ext)
	}

	// Add to appropriate structure based on match type
	switch matchType {
	case "exact":
		pm.suffixMap.AddExact(pattern, value)
	case "suffix":
		pm.suffixMap.AddSuffix(pattern, value)
	case "contains":
		pm.suffixMap.AddContains(pattern, value)
		pm.trie.Insert(pattern, value)
		pm.hasContains = true
	default:
		// Auto-detect based on pattern
		if strings.HasPrefix(pattern, ".") {
			pm.suffixMap.AddSuffix(pattern, value)
		} else if strings.Contains(pattern, "/") || strings.Contains(pattern, "\\") {
			pm.suffixMap.AddContains(pattern, value)
			pm.trie.Insert(pattern, value)
			pm.hasContains = true
		} else {
			pm.suffixMap.AddExact(pattern, value)
		}
	}
}

// Match checks if a path matches any pattern.
// Uses Bloom filter for fast negative check, then detailed match.
func (pm *PatternMatcher[V]) Match(path string) (V, bool) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	base := filepath.Base(path)
	ext := filepath.Ext(path)

	// Fast negative check with Bloom filter
	mightMatch := pm.bloom.MightContain(base) || pm.bloom.MightContain(ext) || pm.bloom.MightContain(path)

	// For contains patterns, also check path segments
	if !mightMatch && pm.hasContains {
		segments := strings.Split(path, string(filepath.Separator))
		for _, seg := range segments {
			if seg != "" && pm.bloom.MightContain(seg) {
				mightMatch = true
				break
			}
		}
	}

	if !mightMatch {
		var zero V
		return zero, false
	}

	// Check suffix map first (most common case)
	if v, ok := pm.suffixMap.Get(path); ok {
		return v, true
	}

	// Fall back to Trie for complex patterns
	if found, v := pm.trie.ContainsSubstring(path); found {
		return v.(V), true
	}

	var zero V
	return zero, false
}
