package dsa

import (
	"container/list"
	"sync"
	"time"
)

// LRUCache is a generic thread-safe LRU cache with TTL support.
// K is the key type (must be comparable), V is the value type.
type LRUCache[K comparable, V any] struct {
	cache   map[K]*list.Element // O(1) lookup
	lru     *list.List          // O(1) eviction order
	maxSize int
	ttl     time.Duration
	mu      sync.RWMutex
}

// lruEntry wraps a value in the LRU list.
// P2 FIX #4: Added ttl field for per-entry TTL support.
type lruEntry[K comparable, V any] struct {
	key       K
	value     V
	timestamp int64
	ttl       int64 // P2 FIX #4: Per-entry TTL in milliseconds (0 = use cache default)
}

// NewLRUCache creates a new generic LRU cache.
//
// maxSize: maximum number of entries (evicts oldest when full)
// ttl: time-to-live for cache entries
func NewLRUCache[K comparable, V any](maxSize int, ttl time.Duration) *LRUCache[K, V] {
	if maxSize <= 0 {
		maxSize = 1000
	}
	return &LRUCache[K, V]{
		cache:   make(map[K]*list.Element, maxSize),
		lru:     list.New(),
		maxSize: maxSize,
		ttl:     ttl,
	}
}

// Get retrieves a cached value if it exists and hasn't expired.
// Returns the value and true if found, zero value and false otherwise.
func (c *LRUCache[K, V]) Get(key K) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	elem, exists := c.cache[key]
	if !exists {
		var zero V
		return zero, false
	}

	ent, ok := elem.Value.(*lruEntry[K, V])
	if !ok {
		c.removeElement(elem)
		var zero V
		return zero, false
	}

	// P2 FIX #4: Check per-entry TTL first, fallback to cache-wide TTL
	now := time.Now().UnixMilli()
	var ttlMs int64
	if ent.ttl > 0 {
		ttlMs = ent.ttl // Use per-entry TTL
	} else if c.ttl > 0 {
		ttlMs = int64(c.ttl.Milliseconds()) // Use cache-wide TTL
	}

	if ttlMs > 0 && now-ent.timestamp > ttlMs {
		c.removeElement(elem)
		var zero V
		return zero, false
	}

	// Move to front (most recently used)
	c.lru.MoveToFront(elem)
	return ent.value, true
}

// Set stores a value in the cache.
// If the cache is at capacity, the oldest entry is evicted.
func (c *LRUCache[K, V]) Set(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// If key already exists, update and move to front
	if elem, exists := c.cache[key]; exists {
		ent, ok := elem.Value.(*lruEntry[K, V])
		if ok {
			ent.value = value
			ent.timestamp = time.Now().UnixMilli()
			c.lru.MoveToFront(elem)
			return
		}
		c.removeElement(elem)
	}

	// Evict oldest if at capacity
	for c.lru.Len() >= c.maxSize {
		oldest := c.lru.Back()
		if oldest != nil {
			c.removeElement(oldest)
		} else {
			break
		}
	}

	// Add new entry
	ent := &lruEntry[K, V]{
		key:       key,
		value:     value,
		timestamp: time.Now().UnixMilli(),
	}
	elem := c.lru.PushFront(ent)
	c.cache[key] = elem
}

// SetWithTTL stores a value with a custom per-entry TTL.
// P2 FIX #4: Now properly stores per-entry TTL instead of ignoring it.
func (c *LRUCache[K, V]) SetWithTTL(key K, value V, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	ttlMs := int64(ttl.Milliseconds())

	// If key already exists, update with new TTL
	if elem, exists := c.cache[key]; exists {
		ent, ok := elem.Value.(*lruEntry[K, V])
		if ok {
			ent.value = value
			ent.timestamp = time.Now().UnixMilli()
			ent.ttl = ttlMs
			c.lru.MoveToFront(elem)
			return
		}
		c.removeElement(elem)
	}

	// Evict oldest if at capacity
	for c.lru.Len() >= c.maxSize {
		oldest := c.lru.Back()
		if oldest != nil {
			c.removeElement(oldest)
		} else {
			break
		}
	}

	// Add new entry with custom TTL
	ent := &lruEntry[K, V]{
		key:       key,
		value:     value,
		timestamp: time.Now().UnixMilli(),
		ttl:       ttlMs,
	}
	elem := c.lru.PushFront(ent)
	c.cache[key] = elem
}

// Delete removes a specific key from the cache.
// Returns true if the key was found and removed.
func (c *LRUCache[K, V]) Delete(key K) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, exists := c.cache[key]; exists {
		c.removeElement(elem)
		return true
	}
	return false
}

// Contains checks if a key exists in the cache (without updating LRU order).
// P2 FIX #4: Now checks per-entry TTL.
func (c *LRUCache[K, V]) Contains(key K) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	elem, exists := c.cache[key]
	if !exists {
		return false
	}

	// P2 FIX #4: Check per-entry TTL first, fallback to cache-wide TTL
	if ent, ok := elem.Value.(*lruEntry[K, V]); ok {
		now := time.Now().UnixMilli()
		var ttlMs int64
		if ent.ttl > 0 {
			ttlMs = ent.ttl
		} else if c.ttl > 0 {
			ttlMs = int64(c.ttl.Milliseconds())
		}
		if ttlMs > 0 && now-ent.timestamp > ttlMs {
			return false
		}
	}
	return true
}

// removeElement removes an element from both the map and the list.
// Must be called with the lock held.
func (c *LRUCache[K, V]) removeElement(elem *list.Element) {
	if ent, ok := elem.Value.(*lruEntry[K, V]); ok {
		delete(c.cache, ent.key)
	}
	c.lru.Remove(elem)
}

// Clear removes all entries from the cache.
func (c *LRUCache[K, V]) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache = make(map[K]*list.Element, c.maxSize)
	c.lru.Init()
}

// Size returns the current number of entries in the cache.
func (c *LRUCache[K, V]) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.lru.Len()
}

// Keys returns all keys in the cache (most recent first).
func (c *LRUCache[K, V]) Keys() []K {
	c.mu.RLock()
	defer c.mu.RUnlock()

	keys := make([]K, 0, c.lru.Len())
	for elem := c.lru.Front(); elem != nil; elem = elem.Next() {
		if ent, ok := elem.Value.(*lruEntry[K, V]); ok {
			keys = append(keys, ent.key)
		}
	}
	return keys
}

// ForEach iterates over all entries (most recent first).
func (c *LRUCache[K, V]) ForEach(fn func(K, V)) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	for elem := c.lru.Front(); elem != nil; elem = elem.Next() {
		if ent, ok := elem.Value.(*lruEntry[K, V]); ok {
			fn(ent.key, ent.value)
		}
	}
}

// Cleanup removes all expired entries.
// P2 FIX #4: Now checks per-entry TTL.
// Returns the number of entries removed.
func (c *LRUCache[K, V]) Cleanup() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now().UnixMilli()
	defaultTTLMs := int64(c.ttl.Milliseconds())
	count := 0

	// Iterate from back (oldest) to front
	for elem := c.lru.Back(); elem != nil; {
		prev := elem.Prev()
		if ent, ok := elem.Value.(*lruEntry[K, V]); ok {
			// P2 FIX #4: Use per-entry TTL if set, otherwise cache-wide TTL
			var ttlMs int64
			if ent.ttl > 0 {
				ttlMs = ent.ttl
			} else {
				ttlMs = defaultTTLMs
			}

			if ttlMs > 0 && now-ent.timestamp > ttlMs {
				c.removeElement(elem)
				count++
			}
		}
		elem = prev
	}

	return count
}

// GetOrSet retrieves a value or sets it if not present.
// The createFn is only called if the key doesn't exist.
func (c *LRUCache[K, V]) GetOrSet(key K, createFn func() V) V {
	// Try to get first
	if val, ok := c.Get(key); ok {
		return val
	}

	// Create and set
	c.mu.Lock()
	defer c.mu.Unlock()

	// Double-check after acquiring write lock
	if elem, exists := c.cache[key]; exists {
		if ent, ok := elem.Value.(*lruEntry[K, V]); ok {
			c.lru.MoveToFront(elem)
			return ent.value
		}
	}

	// Create new value
	value := createFn()

	// Evict if needed
	for c.lru.Len() >= c.maxSize {
		oldest := c.lru.Back()
		if oldest != nil {
			c.removeElement(oldest)
		} else {
			break
		}
	}

	// Add new entry
	ent := &lruEntry[K, V]{
		key:       key,
		value:     value,
		timestamp: time.Now().UnixMilli(),
	}
	elem := c.lru.PushFront(ent)
	c.cache[key] = elem

	return value
}
