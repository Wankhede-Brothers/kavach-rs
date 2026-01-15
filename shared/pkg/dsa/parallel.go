// Package dsa provides parallel loading utilities.
package dsa

import (
	"context"
	"sync"
)

// ParallelResult holds the result of a parallel operation.
type ParallelResult[V any] struct {
	Value V
	Error error
	Index int
}

// ParallelLoader executes multiple loaders concurrently.
// Returns results in the same order as input loaders.
func ParallelLoader[V any](ctx context.Context, loaders []func() (V, error)) []ParallelResult[V] {
	results := make([]ParallelResult[V], len(loaders))
	var wg sync.WaitGroup

	for i, loader := range loaders {
		wg.Add(1)
		go func(idx int, load func() (V, error)) {
			defer wg.Done()

			select {
			case <-ctx.Done():
				results[idx] = ParallelResult[V]{Index: idx, Error: ctx.Err()}
				return
			default:
			}

			val, err := load()
			results[idx] = ParallelResult[V]{Value: val, Index: idx, Error: err}
		}(i, loader)
	}

	wg.Wait()
	return results
}

// ParallelMap applies a function to items in parallel with bounded concurrency.
func ParallelMap[T, R any](ctx context.Context, items []T, maxWorkers int, fn func(T) (R, error)) []ParallelResult[R] {
	if maxWorkers <= 0 {
		maxWorkers = len(items)
	}

	results := make([]ParallelResult[R], len(items))
	jobs := make(chan int, len(items))
	var wg sync.WaitGroup

	// Start workers
	for w := 0; w < maxWorkers && w < len(items); w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for idx := range jobs {
				select {
				case <-ctx.Done():
					results[idx] = ParallelResult[R]{Index: idx, Error: ctx.Err()}
					continue
				default:
				}

				val, err := fn(items[idx])
				results[idx] = ParallelResult[R]{Value: val, Index: idx, Error: err}
			}
		}()
	}

	// Send jobs
	for i := range items {
		jobs <- i
	}
	close(jobs)

	wg.Wait()
	return results
}

// ParallelFiles loads multiple files concurrently.
type FileLoader struct {
	Path   string
	Loader func(path string) ([]byte, error)
}

// LoadFilesParallel loads files in parallel and returns results.
func LoadFilesParallel(ctx context.Context, files []FileLoader) []ParallelResult[[]byte] {
	loaders := make([]func() ([]byte, error), len(files))
	for i, f := range files {
		file := f // Capture for closure
		loaders[i] = func() ([]byte, error) {
			return file.Loader(file.Path)
		}
	}
	return ParallelLoader(ctx, loaders)
}

// OnDemandPool manages a pool of lazily loaded resources.
type OnDemandPool[K comparable, V any] struct {
	pool      map[K]*LazyLoader[V]
	factory   func(K) (V, error)
	mu        sync.RWMutex
	maxActive int
	active    int
}

// NewOnDemandPool creates a pool with maximum active resources.
func NewOnDemandPool[K comparable, V any](maxActive int, factory func(K) (V, error)) *OnDemandPool[K, V] {
	return &OnDemandPool[K, V]{
		pool:      make(map[K]*LazyLoader[V]),
		factory:   factory,
		maxActive: maxActive,
	}
}

// Get retrieves or creates a resource for the given key.
func (p *OnDemandPool[K, V]) Get(key K) (V, error) {
	p.mu.RLock()
	loader, exists := p.pool[key]
	p.mu.RUnlock()

	if exists {
		return loader.Get()
	}

	p.mu.Lock()
	if loader, exists = p.pool[key]; exists {
		p.mu.Unlock()
		return loader.Get()
	}

	// Check pool limits
	if p.maxActive > 0 && p.active >= p.maxActive {
		// Evict oldest unused
		for k, l := range p.pool {
			if !l.IsLoaded() {
				delete(p.pool, k)
				break
			}
		}
	}

	factory := p.factory
	loader = NewLazyLoader(func() (V, error) {
		return factory(key)
	})
	p.pool[key] = loader
	p.active++
	p.mu.Unlock()

	return loader.Get()
}
