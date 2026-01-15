package dsa

import (
	"hash/fnv"
	"sync"
)

// BloomFilter is a space-efficient probabilistic data structure.
// False positives are possible, but false negatives are not.
type BloomFilter struct {
	bits    []uint64
	size    uint64 // number of bits
	numHash int    // number of hash functions
	mu      sync.RWMutex
}

// NewBloomFilter creates a new Bloom filter.
// expectedItems: expected number of items to add
// falsePositiveRate: desired false positive rate (e.g., 0.01 for 1%)
func NewBloomFilter(expectedItems int, falsePositiveRate float64) *BloomFilter {
	if expectedItems <= 0 {
		expectedItems = 100
	}
	if falsePositiveRate <= 0 || falsePositiveRate >= 1 {
		falsePositiveRate = 0.01
	}

	// Optimal size: m = -n*ln(p) / (ln(2)^2)
	// Optimal hash count: k = (m/n) * ln(2)
	n := float64(expectedItems)
	p := falsePositiveRate

	// Calculate optimal size in bits
	m := int64(-n * 2.302585 * 2.302585 / (0.480453 * p)) // ln(p) / ln(2)^2
	if m < 64 {
		m = 64
	}

	// Calculate optimal number of hash functions
	k := int(float64(m) / n * 0.693147) // ln(2)
	if k < 1 {
		k = 1
	}
	if k > 10 {
		k = 10
	}

	// Round up to nearest 64-bit word
	words := (m + 63) / 64

	return &BloomFilter{
		bits:    make([]uint64, words),
		size:    uint64(words * 64),
		numHash: k,
	}
}

// Add inserts an item into the Bloom filter.
func (bf *BloomFilter) Add(item string) {
	bf.mu.Lock()
	defer bf.mu.Unlock()

	h1, h2 := bf.hash(item)
	for i := 0; i < bf.numHash; i++ {
		// Double hashing: hash_i = h1 + i*h2
		pos := (h1 + uint64(i)*h2) % bf.size
		bf.bits[pos/64] |= 1 << (pos % 64)
	}
}

// MightContain checks if an item might be in the set.
// Returns false = definitely not in set
// Returns true = might be in set (could be false positive)
func (bf *BloomFilter) MightContain(item string) bool {
	bf.mu.RLock()
	defer bf.mu.RUnlock()

	h1, h2 := bf.hash(item)
	for i := 0; i < bf.numHash; i++ {
		pos := (h1 + uint64(i)*h2) % bf.size
		if bf.bits[pos/64]&(1<<(pos%64)) == 0 {
			return false // Definitely not present
		}
	}
	return true // Might be present
}

// hash returns two hash values for double hashing.
func (bf *BloomFilter) hash(item string) (uint64, uint64) {
	h := fnv.New64a()
	h.Write([]byte(item))
	h1 := h.Sum64()

	h.Reset()
	h.Write([]byte(item))
	h.Write([]byte{0xff}) // Salt for second hash
	h2 := h.Sum64()

	return h1, h2
}

// Clear removes all items from the filter.
func (bf *BloomFilter) Clear() {
	bf.mu.Lock()
	defer bf.mu.Unlock()

	for i := range bf.bits {
		bf.bits[i] = 0
	}
}

// EstimatedFillRatio returns the approximate fill ratio.
func (bf *BloomFilter) EstimatedFillRatio() float64 {
	bf.mu.RLock()
	defer bf.mu.RUnlock()

	setBits := 0
	for _, word := range bf.bits {
		setBits += popcount(word)
	}
	return float64(setBits) / float64(bf.size)
}

// popcount counts the number of set bits in a uint64.
func popcount(x uint64) int {
	// Brian Kernighan's algorithm
	count := 0
	for x != 0 {
		x &= x - 1
		count++
	}
	return count
}

// AddAll adds multiple items to the filter.
func (bf *BloomFilter) AddAll(items ...string) {
	bf.mu.Lock()
	defer bf.mu.Unlock()

	for _, item := range items {
		h1, h2 := bf.hash(item)
		for i := 0; i < bf.numHash; i++ {
			pos := (h1 + uint64(i)*h2) % bf.size
			bf.bits[pos/64] |= 1 << (pos % 64)
		}
	}
}
