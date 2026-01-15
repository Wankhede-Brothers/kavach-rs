// Package dsa provides efficient data structures for Claude Code enforcement.
package dsa

import "sync"

// Set is a thread-safe generic set with O(1) operations.
type Set[T comparable] struct {
	data map[T]struct{}
	mu   sync.RWMutex
}

// NewSet creates a new Set with optional initial capacity.
func NewSet[T comparable](capacity int) *Set[T] {
	if capacity <= 0 {
		capacity = 16
	}
	return &Set[T]{
		data: make(map[T]struct{}, capacity),
	}
}

// Add inserts an item into the set. Returns true if item was new.
func (s *Set[T]) Add(item T) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.data[item]; exists {
		return false
	}
	s.data[item] = struct{}{}
	return true
}

// Contains checks if an item exists in the set. O(1) average.
func (s *Set[T]) Contains(item T) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, exists := s.data[item]
	return exists
}

// Remove deletes an item from the set. Returns true if item existed.
func (s *Set[T]) Remove(item T) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.data[item]; !exists {
		return false
	}
	delete(s.data, item)
	return true
}

// Size returns the number of items in the set.
func (s *Set[T]) Size() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.data)
}

// Clear removes all items from the set.
func (s *Set[T]) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data = make(map[T]struct{}, 16)
}

// ToSlice returns all items as a slice.
func (s *Set[T]) ToSlice() []T {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]T, 0, len(s.data))
	for item := range s.data {
		result = append(result, item)
	}
	return result
}

// AddAll adds multiple items to the set. Returns count of new items.
func (s *Set[T]) AddAll(items ...T) int {
	s.mu.Lock()
	defer s.mu.Unlock()

	added := 0
	for _, item := range items {
		if _, exists := s.data[item]; !exists {
			s.data[item] = struct{}{}
			added++
		}
	}
	return added
}

// Union returns a new set containing all items from both sets.
func (s *Set[T]) Union(other *Set[T]) *Set[T] {
	s.mu.RLock()
	other.mu.RLock()
	defer s.mu.RUnlock()
	defer other.mu.RUnlock()

	result := NewSet[T](len(s.data) + len(other.data))
	for item := range s.data {
		result.data[item] = struct{}{}
	}
	for item := range other.data {
		result.data[item] = struct{}{}
	}
	return result
}

// Intersection returns a new set containing items present in both sets.
func (s *Set[T]) Intersection(other *Set[T]) *Set[T] {
	s.mu.RLock()
	other.mu.RLock()
	defer s.mu.RUnlock()
	defer other.mu.RUnlock()

	// Iterate over smaller set for efficiency
	smaller, larger := s.data, other.data
	if len(s.data) > len(other.data) {
		smaller, larger = other.data, s.data
	}

	result := NewSet[T](len(smaller))
	for item := range smaller {
		if _, exists := larger[item]; exists {
			result.data[item] = struct{}{}
		}
	}
	return result
}
