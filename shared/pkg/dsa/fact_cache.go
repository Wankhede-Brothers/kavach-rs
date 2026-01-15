// Package dsa provides fact caching with O(1) lookups and TTL expiry.
package dsa

import (
	"sync"
	"time"
)

// Fact represents a cached fact with expiry tracking.
type Fact struct {
	ID        string
	Category  string
	Value     interface{}
	ExpiresAt int64 // Unix timestamp for O(1) comparison
	CreatedAt int64
	Source    string
}

// FactCache provides O(1) fact lookups with automatic expiry.
type FactCache struct {
	facts      map[string]*Fact    // O(1) by ID
	byCategory map[string][]string // Category -> []ID for O(1) category lookup
	mu         sync.RWMutex
}

// NewFactCache creates a new fact cache.
func NewFactCache() *FactCache {
	return &FactCache{
		facts:      make(map[string]*Fact, 256),
		byCategory: make(map[string][]string, 16),
	}
}

// Add inserts or updates a fact. O(1) operation.
func (c *FactCache) Add(fact *Fact) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Remove from old category if exists
	if old, exists := c.facts[fact.ID]; exists {
		c.removeCategoryIndex(old.Category, fact.ID)
	}

	c.facts[fact.ID] = fact
	c.addCategoryIndex(fact.Category, fact.ID)
}

// Get retrieves a fact by ID. Returns nil if not found or expired.
// O(1) operation.
func (c *FactCache) Get(id string) *Fact {
	c.mu.RLock()
	defer c.mu.RUnlock()

	fact, exists := c.facts[id]
	if !exists {
		return nil
	}

	// O(1) expiry check using int64 comparison
	if fact.ExpiresAt > 0 && time.Now().Unix() > fact.ExpiresAt {
		return nil
	}

	return fact
}

// GetByCategory retrieves all non-expired facts in a category.
// O(k) where k is facts in category.
func (c *FactCache) GetByCategory(category string) []*Fact {
	c.mu.RLock()
	defer c.mu.RUnlock()

	ids, exists := c.byCategory[category]
	if !exists {
		return nil
	}

	now := time.Now().Unix()
	result := make([]*Fact, 0, len(ids))

	for _, id := range ids {
		if fact, ok := c.facts[id]; ok {
			if fact.ExpiresAt == 0 || now <= fact.ExpiresAt {
				result = append(result, fact)
			}
		}
	}

	return result
}

// GetActive retrieves all non-expired facts. O(n) but with fast int64 check.
func (c *FactCache) GetActive() []*Fact {
	c.mu.RLock()
	defer c.mu.RUnlock()

	now := time.Now().Unix()
	result := make([]*Fact, 0, len(c.facts))

	for _, fact := range c.facts {
		// O(1) timestamp comparison - no time.Parse needed
		if fact.ExpiresAt == 0 || now <= fact.ExpiresAt {
			result = append(result, fact)
		}
	}

	return result
}

// Remove deletes a fact by ID. O(1) operation.
func (c *FactCache) Remove(id string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	fact, exists := c.facts[id]
	if !exists {
		return false
	}

	c.removeCategoryIndex(fact.Category, id)
	delete(c.facts, id)
	return true
}

// CleanExpired removes all expired facts. Returns count removed.
func (c *FactCache) CleanExpired() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now().Unix()
	removed := 0

	for id, fact := range c.facts {
		if fact.ExpiresAt > 0 && now > fact.ExpiresAt {
			c.removeCategoryIndex(fact.Category, id)
			delete(c.facts, id)
			removed++
		}
	}

	return removed
}

// Size returns the total number of facts (including expired).
func (c *FactCache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.facts)
}

// Categories returns all category names.
func (c *FactCache) Categories() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	cats := make([]string, 0, len(c.byCategory))
	for cat := range c.byCategory {
		cats = append(cats, cat)
	}
	return cats
}

// Helper: add ID to category index
func (c *FactCache) addCategoryIndex(category, id string) {
	c.byCategory[category] = append(c.byCategory[category], id)
}

// Helper: remove ID from category index
func (c *FactCache) removeCategoryIndex(category, id string) {
	ids := c.byCategory[category]
	for i, fid := range ids {
		if fid == id {
			c.byCategory[category] = append(ids[:i], ids[i+1:]...)
			break
		}
	}
}

// ParseTimestamp converts RFC3339 string to Unix timestamp.
// Call once at load time, then use int64 for O(1) comparisons.
func ParseTimestamp(s string) int64 {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return 0
	}
	return t.Unix()
}

// FormatTimestamp converts Unix timestamp to RFC3339 string.
func FormatTimestamp(ts int64) string {
	if ts == 0 {
		return ""
	}
	return time.Unix(ts, 0).Format(time.RFC3339)
}
