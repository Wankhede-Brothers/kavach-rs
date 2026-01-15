package dsa

import (
	"context"
	"fmt"
	"testing"
	"time"
)

// ==================== Set Tests ====================

func TestSet_Basic(t *testing.T) {
	s := NewSet[string](10)

	// Add items
	if !s.Add("a") {
		t.Error("Add should return true for new item")
	}
	if s.Add("a") {
		t.Error("Add should return false for existing item")
	}

	// Contains
	if !s.Contains("a") {
		t.Error("Contains should return true for existing item")
	}
	if s.Contains("b") {
		t.Error("Contains should return false for non-existing item")
	}

	// Size
	if s.Size() != 1 {
		t.Errorf("Size should be 1, got %d", s.Size())
	}

	// Remove
	if !s.Remove("a") {
		t.Error("Remove should return true for existing item")
	}
	if s.Remove("a") {
		t.Error("Remove should return false for non-existing item")
	}
}

func TestSet_AddAll(t *testing.T) {
	s := NewSet[int](10)
	added := s.AddAll(1, 2, 3, 2, 1)
	if added != 3 {
		t.Errorf("AddAll should add 3 unique items, got %d", added)
	}
	if s.Size() != 3 {
		t.Errorf("Size should be 3, got %d", s.Size())
	}
}

func TestSet_Union(t *testing.T) {
	s1 := NewSet[int](10)
	s1.AddAll(1, 2, 3)

	s2 := NewSet[int](10)
	s2.AddAll(3, 4, 5)

	union := s1.Union(s2)
	if union.Size() != 5 {
		t.Errorf("Union size should be 5, got %d", union.Size())
	}
}

func TestSet_Intersection(t *testing.T) {
	s1 := NewSet[int](10)
	s1.AddAll(1, 2, 3)

	s2 := NewSet[int](10)
	s2.AddAll(2, 3, 4)

	intersection := s1.Intersection(s2)
	if intersection.Size() != 2 {
		t.Errorf("Intersection size should be 2, got %d", intersection.Size())
	}
}

// ==================== IndexedSlice Tests ====================

type TestItem struct {
	ID    string
	Value int
}

func TestIndexedSlice_Basic(t *testing.T) {
	s := NewIndexedSlice[string, TestItem](func(item TestItem) string {
		return item.ID
	}, 10)

	// Add
	if !s.Add(TestItem{ID: "a", Value: 1}) {
		t.Error("Add should return true for new item")
	}
	if s.Add(TestItem{ID: "a", Value: 2}) {
		t.Error("Add should return false for duplicate key")
	}

	// Get
	item, found := s.Get("a")
	if !found {
		t.Error("Get should find item")
	}
	if item.Value != 1 {
		t.Errorf("Value should be 1, got %d", item.Value)
	}

	// Update
	if !s.Update("a", TestItem{ID: "a", Value: 10}) {
		t.Error("Update should return true for existing key")
	}
	item, _ = s.Get("a")
	if item.Value != 10 {
		t.Errorf("Value should be 10 after update, got %d", item.Value)
	}

	// Len
	if s.Len() != 1 {
		t.Errorf("Len should be 1, got %d", s.Len())
	}
}

func TestIndexedSlice_AddOrUpdate(t *testing.T) {
	s := NewIndexedSlice[string, TestItem](func(item TestItem) string {
		return item.ID
	}, 10)

	// Add new
	if !s.AddOrUpdate(TestItem{ID: "a", Value: 1}) {
		t.Error("AddOrUpdate should return true for new item")
	}

	// Update existing
	if s.AddOrUpdate(TestItem{ID: "a", Value: 2}) {
		t.Error("AddOrUpdate should return false for update")
	}

	item, _ := s.Get("a")
	if item.Value != 2 {
		t.Errorf("Value should be 2 after update, got %d", item.Value)
	}
}

func TestIndexedSlice_Remove(t *testing.T) {
	s := NewIndexedSlice[string, TestItem](func(item TestItem) string {
		return item.ID
	}, 10)

	s.Add(TestItem{ID: "a", Value: 1})
	s.Add(TestItem{ID: "b", Value: 2})
	s.Add(TestItem{ID: "c", Value: 3})

	if !s.Remove("b") {
		t.Error("Remove should return true for existing key")
	}
	if s.Len() != 2 {
		t.Errorf("Len should be 2 after remove, got %d", s.Len())
	}

	// Verify order preserved
	items := s.Slice()
	if items[0].ID != "a" || items[1].ID != "c" {
		t.Error("Order not preserved after remove")
	}
}

// ==================== BloomFilter Tests ====================

func TestBloomFilter_Basic(t *testing.T) {
	bf := NewBloomFilter(100, 0.01)

	// Add items
	bf.Add("hello")
	bf.Add("world")

	// Check existence
	if !bf.MightContain("hello") {
		t.Error("MightContain should return true for added item")
	}
	if !bf.MightContain("world") {
		t.Error("MightContain should return true for added item")
	}

	// Check non-existence (may have false positives, but "definitely_not_here" is unlikely)
	// Note: This test may occasionally fail due to false positives
	falsePositives := 0
	for i := 0; i < 100; i++ {
		if bf.MightContain(fmt.Sprintf("definitely_not_here_%d", i)) {
			falsePositives++
		}
	}
	// With 1% false positive rate, expect < 5 false positives in 100 checks
	if falsePositives > 5 {
		t.Errorf("Too many false positives: %d", falsePositives)
	}
}

func TestBloomFilter_FillRatio(t *testing.T) {
	bf := NewBloomFilter(100, 0.01)

	if bf.EstimatedFillRatio() != 0 {
		t.Error("Empty filter should have 0 fill ratio")
	}

	for i := 0; i < 50; i++ {
		bf.Add(fmt.Sprintf("item_%d", i))
	}

	ratio := bf.EstimatedFillRatio()
	if ratio <= 0 || ratio >= 1 {
		t.Errorf("Fill ratio should be between 0 and 1, got %f", ratio)
	}
}

// ==================== Trie Tests ====================

func TestTrie_Basic(t *testing.T) {
	trie := NewTrie()

	// Insert
	trie.Insert("hello", "value1")
	trie.Insert("help", "value2")
	trie.Insert("world", "value3")

	// Search
	val, found := trie.Search("hello")
	if !found || val != "value1" {
		t.Error("Search should find 'hello'")
	}

	_, found = trie.Search("hel")
	if found {
		t.Error("Search should not find partial word 'hel'")
	}

	// HasPrefix
	if !trie.HasPrefix("hel") {
		t.Error("HasPrefix should return true for 'hel'")
	}
	if trie.HasPrefix("xyz") {
		t.Error("HasPrefix should return false for 'xyz'")
	}

	// Size
	if trie.Size() != 3 {
		t.Errorf("Size should be 3, got %d", trie.Size())
	}
}

func TestTrie_ContainsSubstring(t *testing.T) {
	trie := NewTrie()
	trie.Insert("node_modules", "blocked")
	trie.Insert(".lock", "lock file")

	// Should find pattern in path
	found, val := trie.ContainsSubstring("/home/user/project/node_modules/express/index.js")
	if !found || val != "blocked" {
		t.Error("ContainsSubstring should find 'node_modules'")
	}

	found, _ = trie.ContainsSubstring("/home/user/project/src/main.go")
	if found {
		t.Error("ContainsSubstring should not find anything in normal path")
	}
}

func TestTrie_Delete(t *testing.T) {
	trie := NewTrie()
	trie.Insert("hello", nil)
	trie.Insert("help", nil)

	if !trie.Delete("hello") {
		t.Error("Delete should return true for existing word")
	}
	if trie.Size() != 1 {
		t.Errorf("Size should be 1 after delete, got %d", trie.Size())
	}

	// "help" should still exist
	if _, found := trie.Search("help"); !found {
		t.Error("'help' should still exist after deleting 'hello'")
	}
}

// ==================== SuffixMap Tests ====================

func TestSuffixMap_Basic(t *testing.T) {
	sm := NewSuffixMap[string](10)

	// Add patterns
	sm.AddExact("package-lock.json", "npm lock file")
	sm.AddSuffix(".lock", "lock file")
	sm.AddContains("node_modules", "dependencies")

	// Test exact match
	val, found := sm.Get("package-lock.json")
	if !found || val != "npm lock file" {
		t.Error("Should match exact filename")
	}

	// Test suffix match
	val, found = sm.Get("/path/to/yarn.lock")
	if !found || val != "lock file" {
		t.Error("Should match suffix")
	}

	// Test contains match
	val, found = sm.Get("/home/user/project/node_modules/express/package.json")
	if !found || val != "dependencies" {
		t.Error("Should match contains pattern")
	}

	// Test no match
	_, found = sm.Get("/home/user/project/src/main.go")
	if found {
		t.Error("Should not match normal file")
	}
}

func TestSuffixMap_MatchExtension(t *testing.T) {
	sm := NewSuffixMap[string](10)
	sm.AddSuffix(".go", "Go source")
	sm.AddSuffix(".ts", "TypeScript source")

	val, found := sm.MatchExtension("/path/to/main.go")
	if !found || val != "Go source" {
		t.Error("Should match .go extension")
	}

	_, found = sm.MatchExtension("/path/to/readme.md")
	if found {
		t.Error("Should not match .md extension")
	}
}

// ==================== LRUCache Tests ====================

func TestLRUCache_Basic(t *testing.T) {
	cache := NewLRUCache[string, int](3, time.Hour)

	// Set and Get
	cache.Set("a", 1)
	cache.Set("b", 2)
	cache.Set("c", 3)

	val, found := cache.Get("a")
	if !found || val != 1 {
		t.Error("Should find 'a' with value 1")
	}

	// Eviction - add 4th item, oldest should be evicted
	cache.Set("d", 4)

	// 'b' should be evicted (oldest not accessed)
	// 'a' was accessed by Get, so 'b' is oldest
	_, found = cache.Get("b")
	if found {
		t.Error("'b' should be evicted")
	}

	if cache.Size() != 3 {
		t.Errorf("Size should be 3, got %d", cache.Size())
	}
}

func TestLRUCache_TTL(t *testing.T) {
	cache := NewLRUCache[string, int](10, 50*time.Millisecond)

	cache.Set("a", 1)

	// Should exist immediately
	if _, found := cache.Get("a"); !found {
		t.Error("Should find 'a' immediately")
	}

	// Wait for TTL to expire
	time.Sleep(60 * time.Millisecond)

	// Should not exist after TTL
	if _, found := cache.Get("a"); found {
		t.Error("'a' should be expired")
	}
}

func TestLRUCache_GetOrSet(t *testing.T) {
	cache := NewLRUCache[string, int](10, time.Hour)

	callCount := 0
	createFn := func() int {
		callCount++
		return 42
	}

	// First call should create
	val := cache.GetOrSet("key", createFn)
	if val != 42 || callCount != 1 {
		t.Error("First call should create value")
	}

	// Second call should not create
	val = cache.GetOrSet("key", createFn)
	if val != 42 || callCount != 1 {
		t.Error("Second call should use cached value")
	}
}

// ==================== PatternMatcher Tests ====================

func TestPatternMatcher_Basic(t *testing.T) {
	pm := NewPatternMatcher[string](100)

	// Add patterns
	pm.AddPattern("package-lock.json", "npm lock", "exact")
	pm.AddPattern(".lock", "lock file", "suffix")
	pm.AddPattern("node_modules", "deps", "contains")

	// Test matches
	val, found := pm.Match("package-lock.json")
	if !found || val != "npm lock" {
		t.Error("Should match exact pattern")
	}

	val, found = pm.Match("/path/to/yarn.lock")
	if !found || val != "lock file" {
		t.Error("Should match suffix pattern")
	}

	val, found = pm.Match("/project/node_modules/express/index.js")
	if !found || val != "deps" {
		t.Error("Should match contains pattern")
	}

	// Test non-match
	_, found = pm.Match("/project/src/main.go")
	if found {
		t.Error("Should not match normal file")
	}
}

// ==================== Benchmarks ====================

func BenchmarkSet_Contains(b *testing.B) {
	s := NewSet[string](1000)
	for i := 0; i < 1000; i++ {
		s.Add(fmt.Sprintf("item_%d", i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Contains("item_500")
	}
}

func BenchmarkIndexedSlice_Get(b *testing.B) {
	s := NewIndexedSlice[string, TestItem](func(item TestItem) string {
		return item.ID
	}, 1000)

	for i := 0; i < 1000; i++ {
		s.Add(TestItem{ID: fmt.Sprintf("item_%d", i), Value: i})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Get("item_500")
	}
}

func BenchmarkBloomFilter_MightContain(b *testing.B) {
	bf := NewBloomFilter(1000, 0.01)
	for i := 0; i < 1000; i++ {
		bf.Add(fmt.Sprintf("item_%d", i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bf.MightContain("item_500")
	}
}

func BenchmarkTrie_ContainsSubstring(b *testing.B) {
	trie := NewTrie()
	patterns := []string{"node_modules", ".lock", "package-lock.json", "yarn.lock", "dist", "build"}
	for _, p := range patterns {
		trie.Insert(p, p)
	}

	testPath := "/home/user/project/node_modules/express/index.js"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		trie.ContainsSubstring(testPath)
	}
}

func BenchmarkLRUCache_Get(b *testing.B) {
	cache := NewLRUCache[string, int](1000, time.Hour)
	for i := 0; i < 1000; i++ {
		cache.Set(fmt.Sprintf("key_%d", i), i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Get("key_500")
	}
}

// Benchmark comparison: slice scan vs indexed lookup
func BenchmarkComparison_SliceScan(b *testing.B) {
	items := make([]TestItem, 100)
	for i := 0; i < 100; i++ {
		items[i] = TestItem{ID: fmt.Sprintf("item_%d", i), Value: i}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, item := range items {
			if item.ID == "item_50" {
				break
			}
		}
	}
}

func BenchmarkComparison_IndexedLookup(b *testing.B) {
	s := NewIndexedSlice[string, TestItem](func(item TestItem) string {
		return item.ID
	}, 100)

	for i := 0; i < 100; i++ {
		s.Add(TestItem{ID: fmt.Sprintf("item_%d", i), Value: i})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Get("item_50")
	}
}

// ==================== LAZY LOADER TESTS ====================

func TestLazyLoader(t *testing.T) {
	callCount := 0
	loader := NewLazyLoader(func() (string, error) {
		callCount++
		return "loaded", nil
	})

	// Not loaded initially
	if loader.IsLoaded() {
		t.Error("should not be loaded initially")
	}

	// First call loads
	val, err := loader.Get()
	if err != nil || val != "loaded" {
		t.Errorf("expected 'loaded', got %v, %v", val, err)
	}
	if callCount != 1 {
		t.Errorf("expected 1 call, got %d", callCount)
	}

	// Second call uses cache
	val, _ = loader.Get()
	if callCount != 1 {
		t.Errorf("expected 1 call (cached), got %d", callCount)
	}

	// Is loaded now
	if !loader.IsLoaded() {
		t.Error("should be loaded after Get")
	}
}

func TestLazyMap(t *testing.T) {
	callCounts := make(map[string]int)
	factory := func(key string) func() (int, error) {
		return func() (int, error) {
			callCounts[key]++
			return len(key), nil
		}
	}

	lm := NewLazyMap[string, int](factory)

	// Load first key
	val, _ := lm.Get("hello")
	if val != 5 {
		t.Errorf("expected 5, got %d", val)
	}
	if callCounts["hello"] != 1 {
		t.Errorf("expected 1 call for 'hello', got %d", callCounts["hello"])
	}

	// Load same key again (cached)
	val, _ = lm.Get("hello")
	if callCounts["hello"] != 1 {
		t.Errorf("expected 1 call (cached), got %d", callCounts["hello"])
	}

	// Load different key
	val, _ = lm.Get("world")
	if val != 5 {
		t.Errorf("expected 5 for 'world', got %d", val)
	}
}

// ==================== FACT CACHE TESTS ====================

func TestFactCache(t *testing.T) {
	fc := NewFactCache()

	// Add facts
	fc.Add(&Fact{
		ID:        "fact1",
		Category:  "syntax",
		Value:     "test value",
		ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
	})

	// O(1) lookup
	fact := fc.Get("fact1")
	if fact == nil {
		t.Error("expected fact, got nil")
	}
	if fact.ID != "fact1" {
		t.Errorf("expected fact1, got %s", fact.ID)
	}

	// Get by category
	facts := fc.GetByCategory("syntax")
	if len(facts) != 1 {
		t.Errorf("expected 1 fact in syntax, got %d", len(facts))
	}
}

func TestFactCacheExpiry(t *testing.T) {
	fc := NewFactCache()

	// Add expired fact
	fc.Add(&Fact{
		ID:        "expired",
		Category:  "test",
		ExpiresAt: time.Now().Add(-1 * time.Hour).Unix(), // Already expired
	})

	// Should return nil for expired
	fact := fc.Get("expired")
	if fact != nil {
		t.Error("expected nil for expired fact")
	}

	// Clean should remove it
	removed := fc.CleanExpired()
	if removed != 1 {
		t.Errorf("expected 1 removed, got %d", removed)
	}
}

// ==================== PARALLEL LOADER TESTS ====================

func TestParallelLoader(t *testing.T) {
	loaders := []func() (int, error){
		func() (int, error) { time.Sleep(10 * time.Millisecond); return 1, nil },
		func() (int, error) { time.Sleep(10 * time.Millisecond); return 2, nil },
		func() (int, error) { time.Sleep(10 * time.Millisecond); return 3, nil },
	}

	start := time.Now()
	results := ParallelLoader(context.Background(), loaders)
	elapsed := time.Since(start)

	// Should complete in ~10ms, not 30ms (parallel)
	if elapsed > 50*time.Millisecond {
		t.Errorf("expected parallel execution, took %v", elapsed)
	}

	// Check results
	for i, r := range results {
		if r.Error != nil {
			t.Errorf("unexpected error at %d: %v", i, r.Error)
		}
		if r.Value != i+1 {
			t.Errorf("expected %d, got %d", i+1, r.Value)
		}
	}
}

// ==================== DYNAMIC CONTEXT TESTS ====================

func TestDynamicContext(t *testing.T) {
	dc := NewDynamicContext(1024 * 1024) // 1MB limit

	loadCount := 0
	dc.Register(&ContextLayer{
		Name:     "config",
		Priority: 10,
		TTL:      1 * time.Hour,
		Loader: func() (interface{}, error) {
			loadCount++
			return map[string]string{"key": "value"}, nil
		},
	})

	// Not loaded initially
	if dc.IsLoaded("config") {
		t.Error("should not be loaded initially")
	}

	// Get triggers load
	data, err := dc.Get("config")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if data == nil {
		t.Error("expected data, got nil")
	}

	// Should be loaded now
	if !dc.IsLoaded("config") {
		t.Error("should be loaded after Get")
	}

	// Second Get uses cache
	dc.Get("config")
	if loadCount != 1 {
		t.Errorf("expected 1 load (cached), got %d", loadCount)
	}
}

func TestDynamicContextDependencies(t *testing.T) {
	dc := NewDynamicContext(1024 * 1024)

	order := []string{}

	dc.Register(&ContextLayer{
		Name:     "base",
		Priority: 10,
		Loader: func() (interface{}, error) {
			order = append(order, "base")
			return "base", nil
		},
	})

	dc.Register(&ContextLayer{
		Name:     "derived",
		Priority: 5,
		Loader: func() (interface{}, error) {
			order = append(order, "derived")
			return "derived", nil
		},
	}, "base") // Depends on base

	// Loading derived should load base first
	dc.Get("derived")

	if len(order) != 2 || order[0] != "base" || order[1] != "derived" {
		t.Errorf("expected [base, derived], got %v", order)
	}
}
