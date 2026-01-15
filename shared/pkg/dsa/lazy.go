// Package dsa provides lazy loading patterns for on-demand data loading.
package dsa

import (
	"sync"
	"sync/atomic"
)

// LazyLoader provides thread-safe on-demand loading with sync.Once semantics.
// V is the value type to be loaded lazily.
type LazyLoader[V any] struct {
	value  atomic.Pointer[V]
	loader func() (V, error)
	once   sync.Once
	err    error
	loaded atomic.Bool
}

// NewLazyLoader creates a lazy loader with the given load function.
// The loader function is only called once, on first access.
func NewLazyLoader[V any](loader func() (V, error)) *LazyLoader[V] {
	return &LazyLoader[V]{
		loader: loader,
	}
}

// Get returns the lazily loaded value, loading it on first access.
// Subsequent calls return the cached value without calling loader.
// Thread-safe: multiple goroutines can safely call Get concurrently.
func (l *LazyLoader[V]) Get() (V, error) {
	l.once.Do(func() {
		val, err := l.loader()
		l.err = err
		if err == nil {
			l.value.Store(&val)
			l.loaded.Store(true)
		}
	})

	if ptr := l.value.Load(); ptr != nil {
		return *ptr, l.err
	}
	var zero V
	return zero, l.err
}

// IsLoaded returns true if the value has been loaded.
func (l *LazyLoader[V]) IsLoaded() bool {
	return l.loaded.Load()
}

// Reset clears the cached value, allowing reload on next Get.
// Use sparingly - breaks the "load once" guarantee.
func (l *LazyLoader[V]) Reset() {
	l.once = sync.Once{}
	l.value.Store(nil)
	l.loaded.Store(false)
	l.err = nil
}

// LazyMap provides lazy loading for map values.
// Keys are loaded on-demand, not all at once.
type LazyMap[K comparable, V any] struct {
	cache   map[K]*LazyLoader[V]
	factory func(K) func() (V, error)
	mu      sync.RWMutex
}

// NewLazyMap creates a lazy map with a factory function.
// The factory creates loaders for each key on demand.
func NewLazyMap[K comparable, V any](factory func(K) func() (V, error)) *LazyMap[K, V] {
	return &LazyMap[K, V]{
		cache:   make(map[K]*LazyLoader[V]),
		factory: factory,
	}
}

// Get retrieves or creates a lazy loader for the given key.
func (m *LazyMap[K, V]) Get(key K) (V, error) {
	m.mu.RLock()
	loader, exists := m.cache[key]
	m.mu.RUnlock()

	if exists {
		return loader.Get()
	}

	m.mu.Lock()
	// Double-check after acquiring write lock
	if loader, exists = m.cache[key]; exists {
		m.mu.Unlock()
		return loader.Get()
	}

	loader = NewLazyLoader(m.factory(key))
	m.cache[key] = loader
	m.mu.Unlock()

	return loader.Get()
}

// IsLoaded checks if a key's value has been loaded.
func (m *LazyMap[K, V]) IsLoaded(key K) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if loader, exists := m.cache[key]; exists {
		return loader.IsLoaded()
	}
	return false
}

// Keys returns all keys that have loaders (loaded or not).
func (m *LazyMap[K, V]) Keys() []K {
	m.mu.RLock()
	defer m.mu.RUnlock()

	keys := make([]K, 0, len(m.cache))
	for k := range m.cache {
		keys = append(keys, k)
	}
	return keys
}

// LoadedKeys returns only keys whose values have been loaded.
func (m *LazyMap[K, V]) LoadedKeys() []K {
	m.mu.RLock()
	defer m.mu.RUnlock()

	keys := make([]K, 0)
	for k, loader := range m.cache {
		if loader.IsLoaded() {
			keys = append(keys, k)
		}
	}
	return keys
}
