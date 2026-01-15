// Package dsa provides dynamic context engineering for agentic systems.
package dsa

import (
	"context"
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

// ContextLayer represents a loadable context layer.
type ContextLayer struct {
	Name     string
	Priority int // Higher = more important
	Loader   func() (interface{}, error)
	TTL      time.Duration
}

// LoadedLayer holds a loaded context layer with metadata.
type LoadedLayer struct {
	Name      string
	Data      interface{}
	LoadedAt  int64
	ExpiresAt int64
	Size      int // Approximate memory size
}

// DynamicContext provides on-demand context loading.
// Only loads what's needed, when it's needed.
type DynamicContext struct {
	layers       map[string]*ContextLayer
	loaded       map[string]*LoadedLayer
	loaders      map[string]*LazyLoader[*LoadedLayer]
	dependencies map[string][]string // Layer -> required layers
	mu           sync.RWMutex
	totalSize    atomic.Int64
	maxSize      int64
}

// NewDynamicContext creates a context manager with max memory limit.
func NewDynamicContext(maxSizeBytes int64) *DynamicContext {
	return &DynamicContext{
		layers:       make(map[string]*ContextLayer),
		loaded:       make(map[string]*LoadedLayer),
		loaders:      make(map[string]*LazyLoader[*LoadedLayer]),
		dependencies: make(map[string][]string),
		maxSize:      maxSizeBytes,
	}
}

// Register adds a context layer definition without loading it.
func (dc *DynamicContext) Register(layer *ContextLayer, dependsOn ...string) {
	dc.mu.Lock()
	defer dc.mu.Unlock()

	dc.layers[layer.Name] = layer
	dc.dependencies[layer.Name] = dependsOn

	// Create lazy loader
	dc.loaders[layer.Name] = NewLazyLoader(func() (*LoadedLayer, error) {
		data, err := layer.Loader()
		if err != nil {
			return nil, err
		}

		now := time.Now().Unix()
		var expiresAt int64
		if layer.TTL > 0 {
			expiresAt = now + int64(layer.TTL.Seconds())
		}

		return &LoadedLayer{
			Name:      layer.Name,
			Data:      data,
			LoadedAt:  now,
			ExpiresAt: expiresAt,
			Size:      estimateSize(data),
		}, nil
	})
}

// Get retrieves a context layer, loading it and dependencies if needed.
// This is the core of "load only when needed" pattern.
func (dc *DynamicContext) Get(name string) (interface{}, error) {
	dc.mu.RLock()
	loader, exists := dc.loaders[name]
	deps := dc.dependencies[name]
	dc.mu.RUnlock()

	if !exists {
		return nil, nil
	}

	// Load dependencies first (recursive)
	for _, dep := range deps {
		if _, err := dc.Get(dep); err != nil {
			return nil, err
		}
	}

	// Load this layer
	loaded, err := loader.Get()
	if err != nil {
		return nil, err
	}

	// Check if expired and needs reload
	if loaded.ExpiresAt > 0 && time.Now().Unix() > loaded.ExpiresAt {
		dc.Invalidate(name)
		loaded, err = loader.Get()
		if err != nil {
			return nil, err
		}
	}

	// Track in loaded map
	// P0 FIX: Prevent memory leak by subtracting old size before adding new
	dc.mu.Lock()
	if existing, exists := dc.loaded[name]; exists {
		// Subtract old size to prevent double-counting
		dc.totalSize.Add(-int64(existing.Size))
	}
	dc.loaded[name] = loaded
	dc.totalSize.Add(int64(loaded.Size))
	dc.mu.Unlock()

	// Evict if over limit
	dc.evictIfNeeded()

	return loaded.Data, nil
}

// IsLoaded checks if a layer is currently loaded.
func (dc *DynamicContext) IsLoaded(name string) bool {
	dc.mu.RLock()
	defer dc.mu.RUnlock()

	loader, exists := dc.loaders[name]
	if !exists {
		return false
	}
	return loader.IsLoaded()
}

// LoadedLayers returns names of currently loaded layers (sorted for deterministic output).
// P1 FIX: Sort output to prevent random order from map iteration.
func (dc *DynamicContext) LoadedLayers() []string {
	dc.mu.RLock()
	defer dc.mu.RUnlock()

	names := make([]string, 0, len(dc.loaded))
	for name := range dc.loaded {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// Invalidate marks a layer for reload on next access.
func (dc *DynamicContext) Invalidate(name string) {
	dc.mu.Lock()
	defer dc.mu.Unlock()

	if loader, exists := dc.loaders[name]; exists {
		loader.Reset()
	}
	if loaded, exists := dc.loaded[name]; exists {
		dc.totalSize.Add(-int64(loaded.Size))
		delete(dc.loaded, name)
	}

	// Also invalidate dependents
	for layerName, deps := range dc.dependencies {
		for _, dep := range deps {
			if dep == name {
				if loader, exists := dc.loaders[layerName]; exists {
					loader.Reset()
				}
				delete(dc.loaded, layerName)
			}
		}
	}
}

// evictIfNeeded removes least important layers if over memory limit.
func (dc *DynamicContext) evictIfNeeded() {
	if dc.maxSize <= 0 {
		return
	}

	dc.mu.Lock()
	defer dc.mu.Unlock()

	// Find lowest priority loaded layers
	for dc.totalSize.Load() > dc.maxSize && len(dc.loaded) > 0 {
		var evictName string
		evictPriority := int(^uint(0) >> 1) // Max int

		for name, loaded := range dc.loaded {
			layer := dc.layers[name]
			if layer != nil && layer.Priority < evictPriority {
				evictPriority = layer.Priority
				evictName = name
				_ = loaded // Use loaded
			}
		}

		if evictName != "" {
			if loaded := dc.loaded[evictName]; loaded != nil {
				dc.totalSize.Add(-int64(loaded.Size))
			}
			delete(dc.loaded, evictName)
			if loader := dc.loaders[evictName]; loader != nil {
				loader.Reset()
			}
		} else {
			break
		}
	}
}

// Preload loads specific layers in parallel (for known-needed contexts).
func (dc *DynamicContext) Preload(ctx context.Context, names ...string) error {
	loaders := make([]func() (interface{}, error), len(names))
	for i, name := range names {
		n := name // Capture
		loaders[i] = func() (interface{}, error) {
			return dc.Get(n)
		}
	}

	results := ParallelLoader(ctx, loaders)
	for _, r := range results {
		if r.Error != nil {
			return r.Error
		}
	}
	return nil
}

// Stats returns current memory usage statistics.
func (dc *DynamicContext) Stats() map[string]interface{} {
	dc.mu.RLock()
	defer dc.mu.RUnlock()

	return map[string]interface{}{
		"total_layers":    len(dc.layers),
		"loaded_layers":   len(dc.loaded),
		"total_size":      dc.totalSize.Load(),
		"max_size":        dc.maxSize,
		"utilization_pct": float64(dc.totalSize.Load()) / float64(dc.maxSize) * 100,
	}
}

// estimateSize provides rough memory size estimation.
func estimateSize(v interface{}) int {
	switch val := v.(type) {
	case string:
		return len(val)
	case []byte:
		return len(val)
	case []string:
		total := 0
		for _, s := range val {
			total += len(s)
		}
		return total
	case map[string]interface{}:
		return len(val) * 64 // Rough estimate
	default:
		return 64 // Default estimate
	}
}
