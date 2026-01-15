package dsa

import "sync"

// IndexedSlice maintains both ordered slice and O(1) lookup map.
// K is the key type, V is the value type.
type IndexedSlice[K comparable, V any] struct {
	items []V
	index map[K]int // key -> slice index
	keyFn func(V) K // extract key from value
	mu    sync.RWMutex
}

// NewIndexedSlice creates a new IndexedSlice with a key extraction function.
func NewIndexedSlice[K comparable, V any](keyFn func(V) K, capacity int) *IndexedSlice[K, V] {
	if capacity <= 0 {
		capacity = 16
	}
	return &IndexedSlice[K, V]{
		items: make([]V, 0, capacity),
		index: make(map[K]int, capacity),
		keyFn: keyFn,
	}
}

// Add appends an item if not already present. Returns true if added.
func (s *IndexedSlice[K, V]) Add(item V) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := s.keyFn(item)
	if _, exists := s.index[key]; exists {
		return false
	}

	s.index[key] = len(s.items)
	s.items = append(s.items, item)
	return true
}

// Get retrieves an item by key. O(1) average.
func (s *IndexedSlice[K, V]) Get(key K) (V, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	idx, exists := s.index[key]
	if !exists {
		var zero V
		return zero, false
	}
	return s.items[idx], true
}

// GetPtr retrieves a pointer to an item by key for in-place updates.
func (s *IndexedSlice[K, V]) GetPtr(key K) *V {
	s.mu.RLock()
	defer s.mu.RUnlock()

	idx, exists := s.index[key]
	if !exists {
		return nil
	}
	return &s.items[idx]
}

// Update replaces an existing item. Returns false if key not found.
func (s *IndexedSlice[K, V]) Update(key K, item V) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	idx, exists := s.index[key]
	if !exists {
		return false
	}
	s.items[idx] = item
	return true
}

// AddOrUpdate adds a new item or updates existing one.
func (s *IndexedSlice[K, V]) AddOrUpdate(item V) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := s.keyFn(item)
	if idx, exists := s.index[key]; exists {
		s.items[idx] = item
		return false // updated
	}

	s.index[key] = len(s.items)
	s.items = append(s.items, item)
	return true // added
}

// Remove deletes an item by key. O(n) due to slice compaction.
func (s *IndexedSlice[K, V]) Remove(key K) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	idx, exists := s.index[key]
	if !exists {
		return false
	}

	// Remove from slice
	s.items = append(s.items[:idx], s.items[idx+1:]...)

	// Remove from index
	delete(s.index, key)

	// Update indices for items after removed one
	for k, i := range s.index {
		if i > idx {
			s.index[k] = i - 1
		}
	}

	return true
}

// Contains checks if a key exists. O(1) average.
func (s *IndexedSlice[K, V]) Contains(key K) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, exists := s.index[key]
	return exists
}

// Slice returns a copy of the underlying slice.
func (s *IndexedSlice[K, V]) Slice() []V {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]V, len(s.items))
	copy(result, s.items)
	return result
}

// Len returns the number of items.
func (s *IndexedSlice[K, V]) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.items)
}

// Clear removes all items.
func (s *IndexedSlice[K, V]) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.items = s.items[:0]
	s.index = make(map[K]int, 16)
}

// ForEach iterates over all items in order.
func (s *IndexedSlice[K, V]) ForEach(fn func(V)) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, item := range s.items {
		fn(item)
	}
}

// Keys returns all keys.
func (s *IndexedSlice[K, V]) Keys() []K {
	s.mu.RLock()
	defer s.mu.RUnlock()

	keys := make([]K, 0, len(s.index))
	for k := range s.index {
		keys = append(keys, k)
	}
	return keys
}
