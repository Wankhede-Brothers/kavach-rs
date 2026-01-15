package dsa

import "sync"

// trieNode represents a node in the Trie.
type trieNode struct {
	children map[rune]*trieNode
	isEnd    bool
	value    interface{} // Associated value at end of word
}

// Trie is a prefix tree for efficient string operations.
type Trie struct {
	root *trieNode
	size int
	mu   sync.RWMutex
}

// NewTrie creates a new empty Trie.
func NewTrie() *Trie {
	return &Trie{
		root: &trieNode{children: make(map[rune]*trieNode)},
	}
}

// Insert adds a string to the Trie with an optional associated value.
func (t *Trie) Insert(key string, value interface{}) {
	t.mu.Lock()
	defer t.mu.Unlock()

	node := t.root
	for _, ch := range key {
		if node.children[ch] == nil {
			node.children[ch] = &trieNode{children: make(map[rune]*trieNode)}
		}
		node = node.children[ch]
	}

	if !node.isEnd {
		t.size++
	}
	node.isEnd = true
	node.value = value
}

// Search checks if a string exists in the Trie.
func (t *Trie) Search(key string) (interface{}, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	node := t.root
	for _, ch := range key {
		if node.children[ch] == nil {
			return nil, false
		}
		node = node.children[ch]
	}
	return node.value, node.isEnd
}

// HasPrefix checks if any string in the Trie starts with the given prefix.
func (t *Trie) HasPrefix(prefix string) bool {
	t.mu.RLock()
	defer t.mu.RUnlock()

	node := t.root
	for _, ch := range prefix {
		if node.children[ch] == nil {
			return false
		}
		node = node.children[ch]
	}
	return true
}

// ContainsSubstring checks if any stored pattern is a substring of the input.
// Returns (found, value) where value is the matched pattern's associated value.
func (t *Trie) ContainsSubstring(input string) (bool, interface{}) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	runes := []rune(input)
	for i := range runes {
		node := t.root
		for j := i; j < len(runes); j++ {
			if node.children[runes[j]] == nil {
				break
			}
			node = node.children[runes[j]]
			if node.isEnd {
				return true, node.value
			}
		}
	}
	return false, nil
}

// ContainsSuffix checks if any stored pattern matches the suffix of input.
func (t *Trie) ContainsSuffix(input string) (bool, interface{}) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	runes := []rune(input)
	for i := len(runes) - 1; i >= 0; i-- {
		node := t.root
		match := true
		for j := i; j < len(runes); j++ {
			if node.children[runes[j]] == nil {
				match = false
				break
			}
			node = node.children[runes[j]]
		}
		if match && node.isEnd {
			return true, node.value
		}
	}
	return false, nil
}

// GetAllWithPrefix returns all strings with the given prefix.
func (t *Trie) GetAllWithPrefix(prefix string) []string {
	t.mu.RLock()
	defer t.mu.RUnlock()

	node := t.root
	for _, ch := range prefix {
		if node.children[ch] == nil {
			return nil
		}
		node = node.children[ch]
	}

	var results []string
	t.collectAll(node, prefix, &results)
	return results
}

// collectAll recursively collects all strings from a node.
func (t *Trie) collectAll(node *trieNode, prefix string, results *[]string) {
	if node.isEnd {
		*results = append(*results, prefix)
	}
	for ch, child := range node.children {
		t.collectAll(child, prefix+string(ch), results)
	}
}

// Size returns the number of strings in the Trie.
func (t *Trie) Size() int {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.size
}

// Clear removes all strings from the Trie.
func (t *Trie) Clear() {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.root = &trieNode{children: make(map[rune]*trieNode)}
	t.size = 0
}

// Delete removes a string from the Trie.
// Returns true if the string was found and removed.
func (t *Trie) Delete(key string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	deleted := false
	t.deleteHelper(t.root, []rune(key), 0, &deleted)
	return deleted
}

// deleteHelper recursively deletes a key and prunes empty nodes.
// Returns true if the current node should be pruned from parent.
func (t *Trie) deleteHelper(node *trieNode, runes []rune, depth int, deleted *bool) bool {
	if depth == len(runes) {
		if !node.isEnd {
			return false
		}
		node.isEnd = false
		node.value = nil
		t.size--
		*deleted = true
		return len(node.children) == 0
	}

	ch := runes[depth]
	child := node.children[ch]
	if child == nil {
		return false
	}

	shouldPrune := t.deleteHelper(child, runes, depth+1, deleted)
	if shouldPrune {
		delete(node.children, ch)
		return !node.isEnd && len(node.children) == 0
	}
	return false
}
